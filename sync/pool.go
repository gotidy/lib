package sync

import "sync"

// Pool is typed implementation of the sync.Pool.
// See: https://pkg.go.dev/sync#Pool
type Pool[T any] struct {
	pool sync.Pool
}

// NewPool creates a new Poll.
func NewPool[T any](new func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{New: func() any { return new }},
	}
}

// Get selects an arbitrary item from the Pool, removes it from the Pool, and returns it to the caller.
// Get may choose to ignore the pool and treat it as empty. Callers should not assume any relation
// between values passed to Put and the values returned by Get.
//
// If Get would otherwise return nil and new func is non-nil, Get returns the result of calling new.
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put adds x to the pool.
func (p *Pool[T]) Put(x T) {
	p.pool.Put(x)
}
