package batcher

import "time"

const (
	// Default value of the batch size
	DefaultBatchSize = 1000
	// Default value of the flush interval.
	DefaultFlushInterval = time.Second
)

type options struct {
	batchSize     uint64
	flushInterval time.Duration
	recuperate    bool
	queueSize     int
}

// Option sets the batcher option.
type Option[T any] func(*options)

// WithBatchSize sets batch size.
func WithBatchSize(size uint64) func(*options) {
	return func(b *options) {
		b.batchSize = size
	}
}

// WithFlushInterval sets flush interval.
func WithFlushInterval(interval time.Duration) func(*options) {
	return func(b *options) {
		b.flushInterval = interval
	}
}

// WithoutBatchReuse prevents the batch from being reused after the handler has been called.
func WithoutBatchReuse() func(*options) {
	return func(b *options) {
		b.recuperate = false
	}
}

// WithQueueSize sets batch queue size. Default value is 1
//
// For preventing parallel batch processing set value to 0.
//
// For cases where the value is greater than 0, the formed batches are processed
// in parallel in turn, while new ones are formed at this time.
func WithQueueSize(size int) func(*options) {
	return func(b *options) {
		b.queueSize = size
	}
}
