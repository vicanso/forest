package global

// 全局使用的缓存配置，可用于缓存一些全局信息

import (
	"sync"

	lru "github.com/hashicorp/golang-lru"
)

const (
	defaultCacheSize   = 1024 * 10
	connectingCountKey = "connecting"
	routeInfosKey      = "routeInfo"
)

var (
	m        = &sync.Map{}
	lruCache *lru.Cache
)

func init() {
	l, err := NewLRU(defaultCacheSize)
	if err != nil {
		panic(err)
	}
	lruCache = l
}

// SaveConnectingCount save the current connecting count
func SaveConnectingCount(v uint32) {
	m.Store(connectingCountKey, v)
}

// GetConnectingCount get the current connecting count
func GetConnectingCount() (connectingCount uint32) {
	v, ok := m.Load(connectingCountKey)
	if !ok || v == nil {
		return 0
	}
	return v.(uint32)
}

// Load get data from cache
func Load(key interface{}) (interface{}, bool) {
	return m.Load(key)
}

// Store store data to cache
func Store(key, value interface{}) {
	m.Store(key, value)
}

// LoadOrStore load the data from cache, if not exists, store it
func LoadOrStore(key, value interface{}) (interface{}, bool) {
	return m.LoadOrStore(key, value)
}

// NewLRU new a lru cache
func NewLRU(size int) (*lru.Cache, error) {
	return lru.New(size)
}

// Add add value to lru cache（default cache）
func Add(key, value interface{}) (evicted bool) {
	return lruCache.Add(key, value)
}

// Get get the value from lru cache
func Get(key interface{}) (value interface{}, found bool) {
	return lruCache.Get(key)
}

// Remove remove the key from lru cache
func Remove(key interface{}) {
	lruCache.Remove(key)
}
