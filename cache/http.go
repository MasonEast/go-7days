package cache

/*
HTTPPool
	- 提供 HTTP 服务的能力;
	- 根据具体的 key，创建 HTTP 客户端从远程节点获取缓存值的能力;
*/

import (
	"cache/consistenthash"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	self string		// 记录自己的地址
	basePath string		// 作为节点间通讯地址的前缀

	mu          sync.Mutex // guards peers and httpGetters
	peers       *consistenthash.Map	// 用来根据具体的key选择节点
	httpGetters map[string]*httpGetter // 映射远程节点与对应的httpGetter 
}

type httpGetter struct {
	baseURL string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self : self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// 服务端功能
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 判断访问路径前缀是否正确
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path:" + r.URL.Path)
	}

	p.Log("%s %s", r.Method, r.URL.Path)

	// 约定访问路径格式为 /<basepath>/<groupname>/<key>
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)

	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if(group == nil) {
		http.Error(w, "no sucn group:"+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

// 客户端功能
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)

	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

var _ PeerGetter = (*httpGetter)(nil)


func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// 实例化了一致性哈希算法，并且添加了传入的节点
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)

	// 为每一个节点创建了一个HTTP客户端 httpGetter
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// 根据具体的key选择节点，并返回节点对应的客户端
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)