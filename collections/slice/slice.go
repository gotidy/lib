package slice

import (
	"fmt"
	"sort"

	"github.com/gotidy/lib/collections/set"
	"github.com/gotidy/lib/constraints"
	"github.com/gotidy/lib/math"
)

// Index returns the index of the first instance of v in s, or -1 if v is not present in s.
func Index[T comparable](s []T, v T) int {
	for i, item := range s {
		if v == item {
			return i
		}
	}

	return -1
}

// Insert the value at the specified position of the slice.
func Insert[T any](s []T, pos int, v T) []T {
	if pos > len(s) || pos < 0 {
		panic(fmt.Sprintf("pos %d is out of range 0..%d", pos, len(s)))
	}
	len := len(s)
	if cap(s) > len {
		s = s[:len+1]
		if pos < len {
			copy(s[pos+1:], s[pos:len])
		}
		s[pos] = v
		return s
	}

	result := make([]T, len+1)
	if pos > 0 {
		copy(result[:pos], s[:pos])
	}
	if pos < len {
		copy(result[pos+1:], s[pos:])
	}
	result[pos] = v
	return result
}

// Reverse items of the slice.
func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < len(s)/2; i++ {
		s[i], s[j] = s[j], s[i]
		j--
	}
}

// Map turns a []T1 to a []T2 using a mapping function.
// This function has two type parameters, T1 and T2.
// This works with slices of any type.
func Map[T1, T2 any](s []T1, f func(T1) T2) []T2 {
	result := make([]T2, len(s))
	for i, v := range s {
		result[i] = f(v)
	}

	return result
}

// MapIndexed turns a []T1 to a []T2 using a mapping function.
// This function has two type parameters, T1 and T2.
// This works with slices of any type.
func MapIndexed[T1, T2 any](s []T1, f func(int, T1) T2) []T2 {
	result := make([]T2, len(s))
	for i, v := range s {
		result[i] = f(i, v)
	}

	return result
}

// Filter filters values from a slice using a filter function.
// It returns a new slice with only the elements of s
// for which f returned true.
func Filter[T any](s []T, f func(T) bool) []T {
	var r []T
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}

	return r
}

// Reduce reduces a []T to a single value using a reduction function.
func Reduce[T, R any](s []T, initializer R, f func(R, T) R) R {
	r := initializer
	for _, v := range s {
		r = f(r, v)
	}

	return r
}

// Diff returns s1 - s2.
func Diff[T comparable](s1, s2 []T) []T {
	if len(s1) == 0 {
		return s1
	}

	var result []T
	s := set.Of(s2...)
	for _, v := range s1 {
		if !s.Has(v) {
			result = append(result, v)
		}
	}
	return result
}

// SymmetricDiff gets the symmetric difference of two sets and gives a set of elements, which are in either of the sets and not in their intersection.
func SymmetricDiff[T comparable](s1, s2 []T) []T {
	var result []T
	set1 := set.Of(s1...)
	set2 := set.Of(s2...)
	for _, v := range s1 {
		if !set2.Has(v) {
			result = append(result, v)
		}
	}
	for _, v := range s2 {
		if !set1.Has(v) {
			result = append(result, v)
		}
	}
	return result
}

// Intersect returns m1 values that keys is contained in m2.
func Intersect[T comparable](s1, s2 []T) []T {
	if len(s1) == 0 || len(s2) == 0 {
		return nil
	}

	var result []T
	if len(s1) < len(s2) {
		s1, s2 = s2, s1
	}
	s := set.Of(s2...)
	for _, v := range s1 {
		if s.Has(v) {
			result = append(result, v)
		}
	}
	return result
}

// Each iterates through map values.
func Each[T any](s []T, f func(v T)) {
	for _, v := range s {
		f(v)
	}
}

// Fold deduplicates items.
func Fold[T comparable](s []T) []T {
	if len(s) == 0 {
		return nil
	}
	m := make(map[T]struct{})
	var result []T
	for _, v := range s {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Group items.
func Group[T any, K comparable](s []T, key func(T) K) map[K][]T {
	m := make(map[K][]T)
	for _, s := range s {
		k := key(s)
		m[k] = append(m[k], s)
	}
	return m
}

// GroupOrder items.
func GroupOrder[T any, K comparable](s []T, key func(T) K, less func(s []T, i, j int) bool) map[K][]T {
	m := Group(s, key)
	for _, s := range m {
		sort.SliceStable(s, func(i, j int) bool { return less(s, i, j) })
	}
	return m
}

// Clone slice.
func Clone[T any](s []T) []T {
	return append([]T(nil), s...)
}

// Min returns the minimum value of the slice.
func Min[T constraints.Ordered](s ...T) T {
	if len(s) == 0 {
		panic("count of items must be 1 or more")
	}
	if len(s) == 1 {
		return s[0]
	}
	return Reduce(s[1:], s[0], math.Min[T])
}

// Max returns the maximum value of the slice.
func Max[T constraints.Ordered](s ...T) T {
	if len(s) == 0 {
		panic("count of items must be 1 or more")
	}
	if len(s) == 1 {
		return s[0]
	}
	return Reduce(s[1:], s[0], math.Max[T])
}

// New allocate fast slice of pointers of specified type.
func New[T any](size int) []*T {
	switch size {
	case 0:
		return nil
	case 1:
		return []*T{new(T)}
	default:
		t := make([]*T, size)
		tt := make([]T, size)
		for i := range t {
			t[i] = &tt[i]
		}
		return t
	}
}

// NewInit allocate fast the slice of pointers of specified type and initializes it.
func NewInit[T any](size int, init func(i int, t *T)) []*T {
	switch size {
	case 0:
		return nil
	case 1:
		p := new(T)
		init(0, p)
		return []*T{p}
	default:
		t := make([]*T, size)
		tt := make([]T, size)
		for i := range t {
			p := &tt[i]
			init(i, p)
			t[i] = p
		}
		return t
	}
}

// // NewUnsafe allocate fast slice of pointers of specified type.
// func NewUnsafe[T any](size int) []*T {
// 	switch size {
// 	case 0:
// 		return nil
// 	default:
// 		var vpT *T
// 		var vT T
// 		ptrSize := int(unsafe.Sizeof(vpT)) * size
// 		typeSize := int(unsafe.Sizeof(vT)) * size

// 		b := make([]byte, ptrSize+typeSize) // allocate memory for slice of pointers and slice of data.

// 		firstPtrItem := (**T)(unsafe.Pointer(&b[0]))
// 		firstTypeItem := (*T)(unsafe.Pointer(&b[ptrSize])) // Data is stored after an slice of pointers.

// 		t := unsafe.Slice(firstPtrItem, size)
// 		tt := unsafe.Slice(firstTypeItem, size)
// 		for i := range t {
// 			t[i] = &tt[i]
// 		}
// 		return t
// 	}
// }

// Batch split the slice to batches and call the callback function with the every batch.
func Batch[T any](s []T, size int, f func([]T) error) error {
	l := len(s)
	for low := 0; low < l; low += size {
		high := low + size
		if high > l {
			high = l
		}
		err := f(s[low:high])
		if err != nil {
			return err
		}
	}
	return nil
}

// ConvertNumbers converts slices of numbers of one number type to slices of another numbers type.
func ConvertNumbers[T, V constraints.Number](t []T) []V {
	v := make([]V, len(t))
	for i, t := range t {
		v[i] = V(t)
	}
	return v
}

// Contains checks that the slice contains the specified value.
func Contains[T comparable](s T, in []T) bool {
	for _, v := range in {
		if s == v {
			return true
		}
	}
	return false
}
