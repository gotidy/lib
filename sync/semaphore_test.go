package sync

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func Hammer(sem *Semaphore, loops int) {
	for i := 0; i < loops; i++ {
		sem.Acquire(context.Background())
		time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond / 1000)
		sem.Release()
	}
}

func TestSemaphore(t *testing.T) {
	t.Parallel()

	n := runtime.GOMAXPROCS(0)
	loops := 10000 / n
	sem := NewSemaphore(n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			Hammer(sem, loops)
		}()
	}
	wg.Wait()
}

func TestSemaphorePanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("release of an unacquired weighted semaphore did not panic")
		}
	}()
	w := NewSemaphore(1)
	w.Release()
}

func TestSemaphoreTryAcquire(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	sem := NewSemaphore(2)
	tries := []bool{}
	sem.Acquire(ctx)
	tries = append(tries, sem.TryAcquire())
	tries = append(tries, sem.TryAcquire())

	sem.Release()
	sem.Release()

	tries = append(tries, sem.TryAcquire())
	sem.Acquire(ctx)
	tries = append(tries, sem.TryAcquire())

	want := []bool{true, false, true, false}
	for i := range tries {
		if tries[i] != want[i] {
			t.Errorf("tries[%d]: got %t, want %t", i, tries[i], want[i])
		}
	}
}

func TestSemaphoreAcquire(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	sem := NewSemaphore(2)
	tryAcquire := func() bool {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		defer cancel()
		return sem.Acquire(ctx) == nil
	}

	tries := []bool{}
	sem.Acquire(ctx)
	tries = append(tries, tryAcquire())
	tries = append(tries, tryAcquire())

	sem.Release()
	sem.Release()

	tries = append(tries, tryAcquire())
	sem.Acquire(ctx)
	tries = append(tries, tryAcquire())

	want := []bool{true, false, true, false}
	for i := range tries {
		if tries[i] != want[i] {
			t.Errorf("tries[%d]: got %t, want %t", i, tries[i], want[i])
		}
	}
}
