package caching

import (
	"sync"
)

const defaultAuth = "Authorization"

var ch = New()

type Cache struct {
	mux  sync.RWMutex
	dict map[interface{}]interface{}
}

func New() *Cache {
	return &Cache{dict: make(map[interface{}]interface{})}
}

func (c *Cache) SafeClean() {
	c.mux.Lock()
	for k := range c.dict {
		delete(c.dict, k)
	}
	c.mux.Unlock()
}

func (c *Cache) SafePut(k, v interface{}) {
	c.mux.Lock()
	c.dict[k] = v
	c.mux.Unlock()
}

func (c *Cache) SafeGet(k interface{}) (interface{}, bool) {
	c.mux.RLock()
	if v, ok := c.dict[k]; ok {
		c.mux.RUnlock()
		return v, ok
	}
	c.mux.RUnlock()
	return nil, false
}

func (c *Cache) SafePop(k interface{}) interface{} {
	c.mux.Lock()
	if v, ok := c.dict[k]; ok {
		delete(c.dict, k)
		c.mux.Unlock()
		return v
	}
	c.mux.Unlock()
	return nil
}

func SafeGet[T any](c *Cache, k interface{}) (T, bool) {
	defer recover()

	if src, ok := c.SafeGet(k); ok && src != nil {
		if dst, ok := src.(T); ok {
			return dst, ok
		}
	}
	var dst T
	return dst, false
}

func (c *Cache) Caching(k interface{}, fn func() interface{}) interface{} {
	v, ok := c.SafeGet(k)
	if ok {
		return v
	}
	rst := fn()

	if err, ok := rst.(error); ok {
		return err
	}

	c.SafePut(k, rst)
	return rst
}
