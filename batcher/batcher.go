package batcher

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Batcher accumulate messages in an internal buffer and flushes it if it is full
// or flush interval is expired.
type Batcher[T any] struct {
	queue      chan T
	close      chan struct{}
	recuperate bool
	wg         sync.WaitGroup
	closed     atomic.Bool

	batchSize     atomic.Uint64
	flushInterval atomic.Int64 // time.Duration

	handler func(ctx context.Context, batch []T)
}

// New creates a new batcher with defined handler and options.
// If handler is nil, then Batcher panic.
// Batches are reused after the handler is called, to prevent this use WithoutBatchReuse.
func New[T any](ctx context.Context, handler func(ctx context.Context, batch []T), opts ...Option[T]) *Batcher[T] {
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

	b := &Batcher[T]{
		handler:    handler,
		close:      make(chan struct{}),
		queue:      make(chan T, o.queueSize),
		recuperate: o.recuperate,
	}

	b.batchSize.Store(o.batchSize)
	b.flushInterval.Store(int64(o.flushInterval))

	b.start(ctx)

	return b
}

// Put the value to the batch.
func (b *Batcher[T]) Put(v T) {
	if b.closed.Load() {
		return
	}

	b.queue <- v
}

// TryPut tries to put the value to the batch.
func (b *Batcher[T]) TryPut(v T) {
	if b.closed.Load() {
		return
	}

	select {
	case b.queue <- v:
	default:
	}
}

func (b *Batcher[T]) start(ctx context.Context) {
	b.wg.Add(1)

	go func() {
		var batch []T
		ctx, cancel := context.WithCancelCause(ctx)

		flush := func() {
			b.handler(ctx, batch)
			if b.recuperate {
				batch = batch[:0]
			} else {
				batch = nil
			}
		}

		defer b.wg.Done()
		for {
			t := time.NewTimer(time.Duration(b.flushInterval.Load()))
			select {
			case <-b.close:
				cancel(errors.New("batch processing stopped"))
				return
			case <-ctx.Done():
			case <-t.C:
				flush()
			case v, ok := <-b.queue:
				t.Stop()
				if !ok {
					flush()
					return
				}

				batch = append(batch, v)
				if len(batch) >= int(b.batchSize.Load()) {
					flush()
				}
			}
		}
	}()
}

func (b *Batcher[T]) Close(flush bool) {
	b.closed.Store(true)
	if flush {
		close(b.queue)
	} else {
		close(b.close)
	}
	b.wg.Wait()
}

// SetBatchSize sets the batch size.
func (b *Batcher[T]) SetBatchSize(size int) {
	b.batchSize.Store(uint64(size))
}

// SetFlushInterval sets the flush interval.
func (b *Batcher[T]) SetFlushInterval(interval time.Duration) {
	b.flushInterval.Store(int64(interval))
}
