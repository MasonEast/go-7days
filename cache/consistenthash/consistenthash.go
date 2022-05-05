package consistenthash

/*
一致性hash算法

*/
import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash Hash
	replicas int	// 虚拟节点倍数
	keys []int	// 哈希环keys
	hashMap map[int]string	// 虚拟节点和真实节点的映射表，键是虚拟节点的hash值，值是真实节点的名称
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash: fn,
		hashMap: make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// 允许传入0或多个真实节点的名称
func (m *Map) Add(keys ...string) {
	// 对每个真实节点对应创建m.replicas个虚拟接滴啊安
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}

	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	// 计算key的hash值
	hash := int(m.hash([]byte(key)))

	// 顺时针找到第一个匹配的虚拟节点下标idx
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 通过hashMap映射得到真实节点
	return m.hashMap[m.keys[idx%len(m.keys)]]
}