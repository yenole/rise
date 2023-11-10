package router

import (
	"reflect"
)

type IMap interface {
	Get(string) (any, bool)
	Set(string, any)
}

func Wrap[T any](c IMap, dist T) {
	key := reflect.TypeOf(dist).String()
	c.Set(key, dist)
}

func Unwrap[T any](c IMap) (t T) {
	key := reflect.TypeOf(t).String()
	value, ok := c.Get(key)
	if !ok {
		return t
	}
	if dist, ok := value.(T); ok {
		return dist
	}
	return
}
