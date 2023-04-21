/**
 * @Author : NewtSun
 * @Date : 2023/4/17 16:12
 * @Description :
 **/

package controller

import (
	"GoCache/cachepolicy"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *cachepolicy.LruCache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Lazy Initialization
	if c.lru == nil {
		c.lru = cachepolicy.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
