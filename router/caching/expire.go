package caching

import (
	"sync/atomic"
	"time"
)

type Expired[T any] struct {
	value     T
	expired   int64
	periold   int64
	onExpired func(v interface{})
}

func Wrap[T any](v T, periold int64, fn func(v interface{})) *Expired[T] {
	unix := time.Now().Unix()
	return &Expired[T]{value: v, expired: unix + periold, periold: periold, onExpired: fn}

}

func (e *Expired[T]) Unwrap() (t T) {
	unix := time.Now().Unix()
	if v := atomic.LoadInt64(&e.expired); v >= unix {
		atomic.StoreInt64(&e.expired, unix+e.periold)
		return e.value
	}
	if e.onExpired != nil {
		e.onExpired(e.value)
	}
	return
}
