package cache

// 根据传入的 key 选择相应节点 PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 用于从对应 group 查找缓存值
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

