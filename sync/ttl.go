package sync

import (
	"sync"
	"sync/atomic"
	"time"
)

// TTLPointer is a pointer to a value that will be updated periodically.
type TTLPointer[T any] struct {
	u       atomic.Pointer[T]
	ttl     atomic.Int64 // time.Duration
	updated atomic.Int64 // Unix time
	mu      sync.Mutex
}

// NewTTLPointer creates a new TTLPointer.
func NewTTLPointer[T any](ttl time.Duration) *TTLPointer[T] {
	p := &TTLPointer[T]{}
	p.ttl.Store(int64(ttl))
	return p
}

// Get returns the value stored in the pointer. If the pointer is expired, it returns nil.
func (p *TTLPointer[T]) Get() *T {
	if p.ttl.Load() == 0 {
		return p.u.Load()
	}
	if time.Now().Sub(time.Unix(p.updated.Load(), 0)) > time.Duration(p.ttl.Load()) {
		return nil
	}
	return p.u.Load()
}

// Set sets the value stored in the pointer.
func (p *TTLPointer[T]) Set(t *T) {
	p.u.Store(t)
	p.updated.Store(time.Now().Unix())
}

// SetTTL sets the TTL of the pointer.
func (p *TTLPointer[T]) SetTTL(ttl time.Duration) {
	p.ttl.Store(int64(ttl))
}

// GetSet returns the value stored in the pointer. If the value is nil, it will return the value of the getter function.
func (p *TTLPointer[T]) GetSet(getter func() (*T, error)) (*T, error) {
	v := p.Get()
	if v != nil {
		return v, nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	// Recheck if the value is still nil.
	v = p.Get()
	if v != nil {
		return v, nil
	}

	var err error
	if v, err = getter(); err != nil {
		return nil, err
	}
	p.Set(v)
	return v, nil
}
