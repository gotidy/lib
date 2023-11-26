package batcher

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// BatcherFast accumulate messages in an internal buffer and flushes it if it is full
// or flush interval is expired.
type BatcherFast[T any] struct {
	batch      []T
	mu         sync.Mutex
	queue      chan []T
	close      chan struct{}
	closeTimer chan struct{}
	pool       sync.Pool
	recuperate bool
	wg         sync.WaitGroup
	closed     atomic.Bool

	batchSize     atomic.Uint64
	flushInterval atomic.Int64 // time.Duration

	handler func(ctx context.Context, batch []T)

	recuperated int64
	allocated   int64
}

// New creates a new batcher with defined handler and options.
// If handler is nil, then BatcherFast panic.
// Batches are reused after the handler is called, to prevent this use WithoutBatchReuse.
// It is not support TryPut().
func NewFast[T any](ctx context.Context, handler func(ctx context.Context, batch []T), opts ...Option[T]) *BatcherFast[T] {
	if handler == nil {
		panic(errors.New("handler must be defined"))
	}

	o := &options{
		batchSize:     DefaultBatchSize,
		flushInterval: DefaultFlushInterval,
		recuperate:    true,
		queueSize:     1,
	}
	for _, opt := range opts {
		opt(o)
	}

	b := &BatcherFast[T]{
		handler:    handler,
		close:      make(chan struct{}),
		closeTimer: make(chan struct{}),
		queue:      make(chan []T, o.queueSize),
		recuperate: o.recuperate,
	}

	b.batchSize.Store(o.batchSize)
	b.flushInterval.Store(int64(o.flushInterval))

	b.pool = sync.Pool{New: func() any {
		atomic.AddInt64(&b.allocated, 1)
		return make([]T, 0, b.batchSize.Load())
	}}

	b.start(ctx)

	return b
}

// Put the value to the batch.
func (b *BatcherFast[T]) Put(v T) {
	if b.closed.Load() {
		return
	}

	b.mu.Lock()
	b.batch = append(b.batch, v)
	if len(b.batch) >= int(b.batchSize.Load()) {
		b.enqueue()
	}
	b.mu.Unlock()
}

func (b *BatcherFast[T]) start(ctx context.Context) {
	b.wg.Add(2)

	go func() {
		defer b.wg.Done()
		for {
			select {
			case <-b.close:
				return
			case batch, ok := <-b.queue:
				b.handler(ctx, batch)
				if !ok {
					return
				}
				if b.recuperate {
					atomic.AddInt64(&b.recuperated, 1)
					b.pool.Put(batch)
				}
			}
		}
	}()

	go func() {
		defer b.wg.Done()
		for {
			select {
			case <-time.After(time.Duration(b.flushInterval.Load())):
				b.mu.Lock()
				b.enqueue()
				b.mu.Unlock()
			case <-b.closeTimer:
				return
			}
		}
	}()
}

func (b *BatcherFast[T]) enqueue() {
	b.queue <- b.batch
	b.batch = b.pool.Get().([]T)[:0]
}

func (b *BatcherFast[T]) Close(flush bool) {
	b.closed.Store(true)

	close(b.closeTimer)

	if flush {
		b.queue <- b.batch
		close(b.queue)
	} else {
		close(b.close)
	}

	b.wg.Wait()
}

// SetBatchSize sets the batch size.
func (b *BatcherFast[T]) SetBatchSize(size int) {
	b.batchSize.Store(uint64(size))
}

// SetFlushInterval sets the flush interval.
func (b *BatcherFast[T]) SetFlushInterval(interval time.Duration) {
	b.flushInterval.Store(int64(interval))
}
