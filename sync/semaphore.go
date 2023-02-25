// Package semaphore provides a weighted semaphore implementation.
package sync

import "context"

// Semaphore provides a way to bound concurrent access to a resource.
type Semaphore struct {
	tokens chan struct{}
}

// NewSemaphore creates a new semaphore with the specified number of tokens
func NewSemaphore(n int) *Semaphore {
	return &Semaphore{tokens: make(chan struct{}, n)}
}

// Acquire waits until a token is available, then acquires it
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.tokens <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryAcquire acquires the semaphore without blocking. On success, returns true.
// On failure, returns false and leaves the semaphore unchanged.
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.tokens <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release releases a token that was previously acquired.
func (s *Semaphore) Release() {
	select {
	case <-s.tokens:
	default:
		panic("nothing to release, 0 tokens acquired")
	}
}
