package cache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// 函数类型实现某一个接口又称为接口型函数，方便使用者在调用时既能传入函数作为参数，也能传入实现了该接口的结构体作为参数。
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 一个Group可以理解为一个缓存的命名空间，name唯一
type Group struct {
	name string
	getter Getter
	mainCache cache
}

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name: name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
	}

	// 将group存储到全局变量groups中
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// 从mainCache中查找缓存，如果存在则返回缓存值， 缓存不存在则调用load方法获取源数据，并将源数据添加到缓存mainCache中
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{
		b : cloneBytes(bytes),
	}

	g.populateCache(key, value)

	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}