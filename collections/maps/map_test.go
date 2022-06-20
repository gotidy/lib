package maps

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	m1 := map[string]int{"1": 2, "2": 3, "3": 4}
	m2 := map[string]int{"2": 4, "5": 6, "7": 8}
	want := map[string]int{"1": 2, "3": 4}
	if got := Diff(m1, m2); !reflect.DeepEqual(got, want) {
		t.Errorf("Diff() = %v, want %v", got, want)
	}
}

func TestDiffKey(t *testing.T) {
	m := map[string]int{"1": 2, "2": 3, "3": 4}
	keys := []string{"2"}
	want := map[string]int{"1": 2, "3": 4}
	if got := DiffKeys(m, keys); !reflect.DeepEqual(got, want) {
		t.Errorf("DiffKeys() = %v, want %v", got, want)
	}
}

func TestUnion(t *testing.T) {
	m1 := map[string]int{"1": 2, "2": 3, "3": 4}
	m2 := map[string]int{"2": 4, "5": 6, "7": 8}
	want := map[string]int{"1": 2, "2": 3, "3": 4, "5": 6, "7": 8}
	if got := Union(m1, m2); !reflect.DeepEqual(got, want) {
		t.Errorf("Union() = %v, want %v", got, want)
	}
}

func TestIntersect(t *testing.T) {
	m1 := map[string]int{"1": 2, "2": 3, "3": 4}
	m2 := map[string]int{"2": 4, "5": 6, "7": 8}
	want := map[string]int{"2": 3}
	if got := Intersect(m1, m2); !reflect.DeepEqual(got, want) {
		t.Errorf("Intersect() = %v, want %v", got, want)
	}
}

func TestSymmetricDiff(t *testing.T) {
	m1 := map[string]int{"1": 2, "2": 3, "3": 4}
	m2 := map[string]int{"2": 4, "5": 6, "7": 8}
	want := map[string]int{"1": 2, "3": 4, "5": 6, "7": 8}
	if got := SymmetricDiff(m1, m2); !reflect.DeepEqual(got, want) {
		t.Errorf("SymmetricDiff() = %v, want %v", got, want)
	}
}

func TestEach(t *testing.T) {
	m1 := map[string]int{"1": 2, "2": 3, "3": 4}
	m2 := map[string]int{}
	Each(m1, func(k string, v int) { m2[k] = v })
	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("Each() = %v, want %v", m2, m1)
	}
}
