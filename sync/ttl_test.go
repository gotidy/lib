package sync

import (
	"errors"
	"testing"
	"time"

	"github.com/gotidy/lib/ptr"
)

func TestTTLPointerSet_Get(t *testing.T) {
	p := NewTTLPointer[int](time.Second * 1)
	expected := 1
	p.Set(ptr.Of(expected))
	switch got := p.Get(); {
	case p == nil:
		t.Errorf("TTLPointer.Get() expected %d but got nil", expected)
	case *got != expected:
		t.Errorf("TTLPointer.Get() expected %d but got %d", expected, *got)
	}
	time.Sleep(time.Second * 2)
	if got := p.Get(); got != nil {
		t.Errorf("TTLPointer.Get() expected nil but got %d", *got)
	}
}

func TestTTLPointerGetSet(t *testing.T) {
	data := []struct {
		Delay    time.Duration
		Value    int
		Expected int
	}{
		{
			Delay:    0,
			Value:    1,
			Expected: 1,
		},
		{
			Delay:    0,
			Value:    2,
			Expected: 1,
		},
		{
			Delay:    time.Second * 2,
			Value:    3,
			Expected: 3,
		},
	}
	p := NewTTLPointer[int](time.Second * 1)
	updaterCallsCount := 0
	for _, data := range data {
		time.Sleep(data.Delay)
		switch got, err := p.GetSet(func() (*int, error) {
			updaterCallsCount++
			return ptr.Of(data.Value), nil
		}); {
		case err != nil:
			t.Errorf("TTLPointer.GetSet() returns error %s", err.Error())
		case p == nil:
			t.Errorf("TTLPointer.GetSet() expected %d but got nil", data.Expected)
		case *got != data.Expected:
			t.Errorf("TTLPointer.GetSet() expected %d but got %d", data.Expected, *got)
		}
	}
	if updaterCallsCount != 2 {
		t.Errorf("expected 2 calls of TTLPointer.GetSet() but got %d", updaterCallsCount)
	}
}

func TestTTLPointerGetSet_Error(t *testing.T) {
	p := NewTTLPointer[int](time.Second)
	got, err := p.GetSet(func() (*int, error) { return nil, errors.New("") })
	if err == nil {
		t.Errorf("TTLPointer.GetSet() expected error but got nil")
	}
	if got != nil {
		t.Errorf("TTLPointer.GetSet() expected nil but got %v", got)
	}
}
