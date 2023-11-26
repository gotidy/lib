package slice

import (
	"fmt"
	"sort"

	"github.com/gotidy/lib/collections/set"
	"github.com/gotidy/lib/constraints"
	"github.com/gotidy/lib/math"
	"github.com/gotidy/lib/ptr"
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
func Insert[T any](s []T, pos int, v ...T) []T {
	if pos > len(s) || pos < 0 {
		panic(fmt.Sprintf("pos %d is out of range 0..%d", pos, len(s)))
	}
	length := len(s)
	if cap(s) >= length+len(v) {
		s = s[:length+len(v)]
		if pos < length {
			copy(s[pos+len(v):], s[pos:length])
		}
		copy(s[pos:], v)
		return s
	}

	result := make([]T, length+len(v))
	if pos > 0 {
		copy(result[:pos], s[:pos])
	}
	if pos < length {
		copy(result[pos+len(v):], s[pos:])
	}
	copy(result[pos:], v)
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

// MapFilter turns a []T1 to a []T2 using a mapping function,
// Values is not placed to result slice when the mapping function return false.
// The resulting slice may have a smaller size than the original.
func MapFilter[T1, T2 any](s []T1, f func(T1) (T2, bool)) []T2 {
	result := make([]T2, 0, len(s))
	for _, v := range s {
		if r, ok := f(v); ok {
			result = append(result, r)
		}
	}

	return result
}

// MapNotNil turns a []*T1 to a []T2 using a mapping function, exclude nil values.
// This works with slices of any type.
// The resulting slice may have a smaller size than the original.
func MapNotNil[T1, T2 any](s []*T1, f func(*T1) T2) []T2 {
	result := make([]T2, 0, len(s))
	for _, v := range s {
		if v != nil {
			result = append(result, f(v))
		}
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

// FilterMap first filters values from a slice using a filter function and
// then map them.
// It returns a new slice with only the mapped elements of s
// for which filter returned true.
func FilterMap[T1, T2 any](s []T1, filter func(T1) bool, mapper func(T1) T2) []T2 {
	var r []T2
	for _, v := range s {
		if filter(v) {
			r = append(r, mapper(v))
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

// Union slices into one. Duplicates is removed. Order of items in slices is preserved.
func Union[T comparable](ss ...[]T) []T {
	if len(ss) == 0 {
		return nil
	}
	if len(ss) == 1 {
		return ss[0]
	}
	var result []T
	var count, max int
	for _, s := range ss {
		l := len(s)
		if l == 0 {
			continue
		}
		count++
		if l > max {
			max = l
			result = s
		}
	}
	if count == 1 {
		return result
	}

	set := make(set.Set[T], max)
	result = make([]T, 0, max)
	for _, s := range ss {
		for _, item := range s {
			if set.TryAdd(item) {
				result = append(result, item)
			}
		}
	}

	return result
}

// Merge slices into one.
func Merge[T any](ss ...[]T) []T {
	if len(ss) == 0 {
		return nil
	}

	l := 0
	for _, s := range ss {
		l += len(s)
	}

	res := make([]T, 0, l)
	for _, s := range ss {
		res = append(res, s...)
	}
	return res
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

// FoldFunc deduplicates items by keys.
func FoldFunc[T, K comparable](s []T, key func(T) K) []T {
	if len(s) == 0 {
		return nil
	}
	m := make(map[K]struct{})
	result := make([]T, 0, len(s))
	for _, v := range s {
		k := key(v)
		if _, ok := m[k]; !ok {
			m[k] = struct{}{}
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
func NewInit[T any](size int, init func(i int, item *T)) []*T {
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

// NewInitFilter allocate fast the slice of pointers of specified type and initializes it.
func NewInitFilter[T any](size int, init func(i int, item *T) bool) []*T {
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
		outputIndex := 0
		for i := range t {
			p := &tt[outputIndex]
			if ok := init(i, p); ok {
				t[outputIndex] = p
				outputIndex++
			}
		}
		return t[:outputIndex]
	}
}

// NewFrom allocate fast the slice of pointers of specified type and initializes it.
func NewFrom[K, T any](source []K, init func(dst *T, src K)) []*T {
	switch size := len(source); size {
	case 0:
		return nil
	case 1:
		p := new(T)
		init(p, source[0])
		return []*T{p}
	default:
		t := make([]*T, size)
		tt := make([]T, size)
		for i := range t {
			p := &tt[i]
			init(p, source[i])
			t[i] = p
		}
		return t
	}
}

// NewFromFilter allocate fast the slice of pointers of specified type and initializes it.
// It skips output if the init function return false.
func NewFromFilter[K, T any](source []K, init func(dst *T, src K) bool) []*T {
	switch size := len(source); size {
	case 0:
		return nil
	case 1:
		p := new(T)
		init(p, source[0])
		return []*T{p}
	default:
		t := make([]*T, size)
		tt := make([]T, size)
		outputIndex := 0
		for i := range t {
			p := &tt[i]
			if init(p, source[i]) {
				t[outputIndex] = p
				outputIndex++
			}
		}
		return t[:outputIndex]
	}
}

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

// Equal compare two slices.
func Equal[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, s := range s1 {
		if s != s2[i] {
			return false
		}
	}
	return true
}

// Append not nil elements.
func AppendNotNil[T any](slice []*T, elems ...*T) []*T {
	if len(elems) == 0 {
		return slice
	}

	i := 0
	for i < len(elems) {
		if elems[i] != nil {
			i++
			continue
		}
		if i > 0 {
			slice = append(slice, elems[:i]...)
		}
		if len(elems) == 1 {
			elems = nil
			break
		}
		elems = elems[i+1:]
		i = 0
	}
	return append(slice, elems...)
}

// ProcessNotNil process not nil elements.
func ProcessNotNil[T any](s []*T, f func(*T) error) error {
	for _, item := range s {
		if item == nil {
			continue
		}
		if err := f(item); err != nil {
			return err
		}
	}
	return nil
}

// ToMap convert slice to map.
func ToMap[I any, K comparable, T any](items []I, fn func(item I) (K, T)) map[K]T {
	result := make(map[K]T, len(items))
	for _, item := range items {
		k, t := fn(item)
		result[k] = t
	}
	return result
}

// FromMap extract slice from map.
func FromMap[I any, K comparable, T any](items map[K]T, fn func(key K, value T) I) []I {
	result := make([]I, 0, len(items))
	for key, value := range items {
		result = append(result, fn(key, value))
	}
	return result
}

// FindFirst item.
func FindFirst[T any](items []T, filter func(item T) bool) (T, bool) {
	for _, item := range items {
		if filter(item) {
			return item, true
		}
	}
	return ptr.Zero[T](), false
}

// Count values.
func Count[T any](s []T, f func(T) bool) int {
	var count int
	for _, item := range s {
		if f(item) {
			count++
		}
	}
	return count
}

// Remove the element at position. It returns the same slice with reduced size.
func Remove[T any](s []T, pos int) []T {
	return RemoveMany(s, pos, 1)
}

// Remove the elements at position. It returns the same slice with reduced size.
func RemoveMany[T any](s []T, pos int, amount int) []T {
	if pos > len(s) || pos < 0 {
		panic(fmt.Sprintf("pos %d is out of range 0..%d", pos, len(s)))
	}
	if amount < 0 || amount+pos > len(s) {
		panic(fmt.Sprintf("amount %d is out of range 0..%d", amount, len(s)-pos))
	}
	copy(s[pos:], s[pos+amount:])
	return s[:len(s)-amount]
}

// Remove elements from the slice. It returns the same slice with reduced size.
func RemoveFunc[T any](s []T, remove func(item T) bool) []T {
	if s == nil {
		return nil
	}

	var current int
	for i := 0; i < len(s); i++ {
		if remove(s[i]) {
			continue
		}
		if current != i {
			s[current] = s[i]
		}
		current++
	}
	return s[:current]
}

// HasDuplicates checks if the slice has deduplicates.
func HasDuplicates[T comparable](s []T) bool {
	if len(s) < 2 {
		return false
	}
	m := make(map[T]struct{})
	for _, v := range s {
		if _, ok := m[v]; ok {
			return true
		}
		m[v] = struct{}{}
	}
	return false
}

// HasDuplicatesFunc checks if the slice has deduplicates.
func HasDuplicatesFunc[T, K comparable](s []T, key func(T) K) bool {
	if len(s) < 2 {
		return false
	}
	m := make(map[K]struct{})
	for _, v := range s {
		k := key(v)
		if _, ok := m[k]; ok {
			return true
		}
		m[k] = struct{}{}
	}
	return false
}

// Grow slice capacity.
func Grow[T any](s []T, capacity int) []T {
	if capacity < 0 {
		panic("capacity cannot be negative")
	}

	if capacity < len(s) || capacity <= cap(s) {
		return s
	}

	return append(make([]T, 0, capacity), s...)
}

// MergeSortedTo merges two sorted slices to one sorted. Dst may be overwritten if capacity is enough.
// s1 or s2 may be used as dst (for collecting values from many slices to one).
// If limit is zero then result size is not limited.
func MergeSortedTo[T any](dst, src1, src2 []T, less func(v1, v2 T) bool, limit int) []T {
	size1 := len(src1)
	size2 := len(src2)
	i1 := len(src1) - 1
	i2 := len(src2) - 1
	switch {
	case limit == 0:
		limit = size1 + size2
	case limit < 0:
		panic("capacity cannot be negative")
	default:
		limit = math.Min(limit, size1+size2)
	}
	dst = Grow(dst, size1+size2)[:size1+size2]
	for k := size1 + size2 - 1; k >= 0; k-- {
		if i1 >= 0 && i2 >= 0 && !less(src1[i1], src2[i2]) {
			dst[k] = src1[i1]
			i1--
		} else if i2 >= 0 {
			dst[k] = src2[i2]
			i2--
		}
	}

	return dst[:limit]
}

// MergeSorted merges two sorted slices to one new sorted.
// If limit is zero then result size is not limited.
func MergeSorted[T any](s1 []T, s2 []T, less func(v1, v2 T) bool, limit int) []T {
	size1 := len(s1)
	size2 := len(s2)
	switch {
	case limit == 0:
		limit = size1 + size2
	case limit < 0:
		panic("capacity cannot be negative")
	default:
		limit = math.Min(limit, size1+size2)
	}
	dst := make([]T, limit)
	for k, i1, i2 := 0, 0, 0; k < limit; k++ {
		if i1 < size1 && i2 < size2 && less(s1[i1], s2[i2]) {
			dst[k] = s1[i1]
			i1++
		} else if i2 < size2 {
			dst[k] = s2[i2]
			i2++
		}
	}

	return dst
}

// Sort slices ascending.
func Sort[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
}

// Sort slices descending.
func SortDesc[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool { return s[i] > s[j] })
}

// FitIndex fits the index into slice range.
// If the slice is empty the result will be -1.
func FitIndex[T any](index int, s []T) int {
	return math.MustBetween(index, 0, len(s)-1)
}

// Last returns last element. It panics if slice has zero length.
func Last[T constraints.Ordered](s []T) T {
	return s[len(s)-1]
}

// LastExists returns last element.
func LastExists[T constraints.Ordered](s []T) (element T, exists bool) {
	if len(s) == 0 {
		return ptr.Zero[T](), false
	}
	return s[len(s)-1], true
}
