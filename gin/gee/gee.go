/*
 * @Description:
 * @Date: 2022-04-30 10:20:55
 * @Author: mason
 */
package gee

/*
实现思路：
1. 从http.ListenAndServer的第二参数入手，这个参数从net/http源码可以了解，它是一个实现了ServeHTTP方法的接口，所以这里
只要我们自己传入一个实现了ServeHTTP的接口实例，就可以拿到请求自行处理；
2. 创建一个Engine的空结构体，并实现ServeHTTP方法，让它作为http.ListenAndServer的第二参数；
3. ServeHTTP方法接收两个参数：w http.ResponseWriter, req *http.Request，我们可以抽出一个Context，用来接收这两个参数
并做一些封装处理；
4. Context可以根据这两个参数实现，获取请求参数，对外暴露设置接口返回头，状态，返回体的方法（String，JSON， Data， HTML等）
5. 封装Engine，想让它可以处理各种请求可以抽出Router；
6. router主要是实现添加各种请求，并维护一个每种请求对应handler的map；
------------

7. router使用map存储有个问题就是不支持动态路由，这时可以考虑使用前缀树结构来解决；
	- 首先添加路由的时候，要在routers的roots里记录每种请求方法的请求路径的根路径，通过递归将子节点就作为根路径的children存储，遇到*或者:
		就将该子节点标记为非精确匹配；
	- 当收到请求后，根据根路径去匹配，通过递归一层层匹配，如果请求路径中有:或者*，就可以匹配到之前标记为非精确匹配的子节点；

8. 添加路由分组：通过路由分组我们可以共用路由前缀，还可以对分组使用特有的中间件；
	- 增加RouterGroup；
	- RouterGroup需要访问路由，我们可以在Group中保存一个Engine的指针；
	- Engine作为顶层分组，保存RouterGroup的指针，拥有RouterGroup所有能力；
	- 这样我们就可以将所有的路由相关函数交给RouterGroup处理；

9. 使用中间件：中间件允许用户在handler中做一些自定义操作，并且可以改变Context；
	- 通过next方法支持使用多个中间件链式调用；
	- 中间件应该保存在Context中，因为中间件要支持在用户handler执行后依然可以调用；
	- 当我们接收到一个具体请求时，要判断该请求适用于哪些中间件，在这里简单通过 URL 的前缀来判断
*/
import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 添加路由分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
