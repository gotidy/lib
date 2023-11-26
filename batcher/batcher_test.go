package batcher

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"
)

const batchSize = 100

func Benchmark_No_Recuperate_Put(b *testing.B) {
	batch := New(
		context.Background(),
		func(ctx context.Context, batch []string) {
			// b.Logf("Batch: %+v", batch)
		},
		WithQueueSize(1000),
		WithBatchSize(batchSize),
		WithFlushInterval(time.Minute),
		WithoutBatchReuse(),
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		batch.Put("some data")
	}

	b.StopTimer()

	batch.Close(true)
}

func Benchmark_Put(b *testing.B) {
	batch := New(
		context.Background(),
		func(ctx context.Context, batch []string) {
			// b.Logf("Batch: %+v", batch)
		},
		WithQueueSize(1000),
		WithBatchSize(batchSize),
		WithFlushInterval(time.Minute),
	)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		batch.Put("some data")
	}

	b.StopTimer()

	batch.Close(true)
}

func Benchmark_No_Recuperate_Put_Fast(b *testing.B) {
	batch := NewFast(
		context.Background(),
		func(ctx context.Context, batch []string) {
			// b.Logf("Batch: %+v", batch)
		},
		WithBatchSize(batchSize),
		WithFlushInterval(time.Minute),
		WithoutBatchReuse(),
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		batch.Put("some data")
	}

	b.StopTimer()

	batch.Close(true)
}

func Benchmark_Put_Fast(b *testing.B) {
	batch := NewFast(
		context.Background(),
		func(ctx context.Context, batch []string) {
			// b.Logf("Batch: %+v", batch)
		},
		WithBatchSize(batchSize),
		WithFlushInterval(time.Minute),
	)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		batch.Put("some data")
	}

	b.StopTimer()

	batch.Close(true)
}

func TestBatcher(t *testing.T) {
	var data []int

	flusher := func(ctx context.Context, batch []int) {
		for _, v := range batch {
			data = append(data, v)
		}
	}

	batch := New(
		context.Background(),
		flusher,
		WithBatchSize(10),
		WithFlushInterval(time.Second),
	)

	var sample []int
	for i := 0; i < 30; i++ {
		batch.Put(i)
		sample = append(sample, i)
	}

	batch.Close(true)

	sort.Ints(data)

	if !reflect.DeepEqual(data, sample) {
		t.Errorf("expected %+v, got %+v", sample, data)
	}
}

func TestBatcherPeriodicalFlush(t *testing.T) {
	var data []int

	flusher := func(ctx context.Context, batch []int) {
		for _, v := range batch {
			data = append(data, v)
		}
	}

	batch := New(
		context.Background(),
		flusher,
		WithBatchSize(10),
		WithFlushInterval(time.Second),
	)

	var sample []int

	for i := 0; i < 5; i++ {
		batch.Put(i)
		sample = append(sample, i)
	}

	time.Sleep(time.Second * 5)

	for i := 5; i < 10; i++ {
		batch.Put(i)
		sample = append(sample, i)
	}

	batch.Close(true)

	sort.Ints(data)

	if !reflect.DeepEqual(data, sample) {
		t.Errorf("expected %+v, got %+v", sample, data)
	}
}

func TestBatcherNewPanic(t *testing.T) {
	defer func() { recover() }()
	_ = New[int](context.Background(), nil, WithBatchSize(10), WithFlushInterval(time.Second))
	t.Error("must panic when flusher is nil")
}

func TestBatcherPutToClosed(t *testing.T) {
	var data []int

	flusher := func(ctx context.Context, batch []int) {
		for _, v := range batch {
			data = append(data, v)
		}
	}

	batch := New(context.Background(), flusher, WithBatchSize(1), WithFlushInterval(time.Second))

	batch.Put(1)

	batch.Close(true)

	if !reflect.DeepEqual(data, []int{1}) {
		t.Errorf("expected %+v, got %+v", []int{1}, data)
	}

	data = nil

	batch.Put(2)

	if data != nil {
		t.Errorf("closed batcher must not accumulate data")
	}
}
