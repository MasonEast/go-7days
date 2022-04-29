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


*/
import (
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	e.router.addRoute(method, pattern, handler)
}

func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}