package iters

import (
	"iter"

	"github.com/gotidy/lib/collections/set"
)

// Filter filters values from the the sequence ov values using a filter function.
func Filter[V any](seq iter.Seq[V], f func(item V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if f(v) && !yield(v) {
				return
			}
		}
	}
}

// Filter2 filters values from the the sequence ov key-values pairs using a filter function.
func Filter2[K, V any](seq iter.Seq2[K, V], f func(k K, v V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if f(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

// Map converts the sequence of values to the sequence of values another type using a mapping function.
func Map[V1, V2 any](seq iter.Seq[V1], f func(v V1) V2) iter.Seq[V2] {
	return func(yield func(V2) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Map converts the sequence of values to the sequence of values another type using a mapping function.
func Map2[K1, K2, V1, V2 any](seq iter.Seq2[K1, V1], f func(k K1, v V1) (K2, V2)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range seq {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// MapValues converts the sequence of key-value pairs using a values mapping function.
func MapValues[K, V1, V2 any](seq iter.Seq2[K, V1], f func(item V1) V2) iter.Seq2[K, V2] {
	return func(yield func(K, V2) bool) {
		for k, v := range seq {
			if !yield(k, f(v)) {
				return
			}
		}
	}
}

// MapKeys converts the sequence of key-value pairs using a keys mapping function.
func MapKeys[K1, K2, V any](seq iter.Seq2[K1, V], f func(item K1) K2) iter.Seq2[K2, V] {
	return func(yield func(K2, V) bool) {
		for k, v := range seq {
			if !yield(f(k), v) {
				return
			}
		}
	}
}

// NotNil skips nil values in the sequence.
func NotNil[V any, P *V](seq iter.Seq[P]) iter.Seq[P] {
	return Filter(seq, func(p P) bool {
		return p != nil
	})
}

// NotNilValues skips nil values in the sequence.
func NotNilValues[K, V any, P *V](seq iter.Seq2[K, P]) iter.Seq2[K, P] {
	return Filter2(seq, func(_ K, p P) bool {
		return p != nil
	})
}

// NotEmpty skips zero values in the sequence.
func NotEmpty[V comparable](seq iter.Seq[V]) iter.Seq[V] {
	return Filter(seq, func(v V) bool {
		var zero V
		return v != zero
	})
}

// NotEmptyValues skips zero values in the sequence.
func NotEmptyValues[K any, V comparable](seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return Filter2(seq, func(k K, v V) bool {
		var zero V
		return v != zero
	})
}

// WithKeys converts the sequence of values to the sequence of key-value pairs by adding key to the sequence.
func WithKeys[K, V any](seq iter.Seq[V], f func(item V) K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for v := range seq {
			if !yield(f(v), v) {
				return
			}
		}
	}
}

// ToSeq2 converts the sequence of of individual values to the sequence of key-value pairs.
func ToSeq2[T, K, V any](seq iter.Seq[T], f func(item T) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Fold push items from one sequence to another sequence skipping duplicates.
func Fold[V comparable](seq iter.Seq[V]) iter.Seq[V] {
	return FoldFunc(seq, func(v V) V { return v })
}

// Fold push items from one sequence to another sequence skipping duplicates.
func FoldFunc[K comparable, V any](seq iter.Seq[V], foldKey func(V) K) iter.Seq[V] {
	return func(yield func(V) bool) {
		m := make(map[K]struct{})
		for v := range seq {
			key := foldKey(v)
			if _, ok := m[key]; ok {
				continue
			}
			m[key] = struct{}{}
			if !yield(v) {
				return
			}
		}
	}
}

// Fold2 push items from one sequence to another sequence skipping duplicates.
func Fold2[K comparable, V any](seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m := make(map[K]struct{})
		for k, v := range seq {
			if _, ok := m[k]; ok {
				continue
			}
			m[k] = struct{}{}
			if !yield(k, v) {
				return
			}
		}
	}
}

// Fold2Func push items from one sequence to another sequence skipping duplicates.
func Fold2Func[F comparable, K, V any](seq iter.Seq2[K, V], foldKey func(K, V) F) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m := make(map[F]struct{})
		for k, v := range seq {
			key := foldKey(k, v)
			if _, ok := m[key]; ok {
				continue
			}
			m[key] = struct{}{}
			if !yield(k, v) {
				return
			}
		}
	}
}

// Reduce reduces a sequence to a single value using a reduction function.
func Reduce[T, R any](seq iter.Seq[T], initializer R, f func(R, T) R) R {
	r := initializer
	for v := range seq {
		r = f(r, v)
	}

	return r
}

// Values convert Seq2 to a Seq by returning the values of the sequence.
func Values[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}

// Keys convert Seq2 to a Seq by returning the keys of the sequence.
func Keys[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range seq {
			if !yield(k) {
				return
			}
		}
	}
}

// Contains checks that the sequence contains the specified value.
func Contains[T comparable](s T, in iter.Seq[T]) bool {
	for v := range in {
		if s == v {
			return true
		}
	}
	return false
}

// Equal compare two sequences. Slow.
func Equal[T comparable](s1, s2 iter.Seq[T]) bool {
	next1, stop1 := iter.Pull(s1)
	next2, stop2 := iter.Pull(s2)
	for {
		v1, ok1 := next1()
		v2, ok2 := next2()
		if ok1 != ok2 || (ok1 && ok2 && v1 != v2) {
			stop1()
			stop2()
			return false
		}
		if !ok1 {
			stop1()
			stop2()
			return true
		}
	}
}

// Merge sequences of values into one.
func Merge[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Merge2 sequences of key-value pairs into one.
func Merge2[K, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, seq := range seqs {
			for k, v := range seq {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Diff returns s1 - s2.
func Diff[V comparable](s1, s2 iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		s := set.OfSeq(s2)
		for v := range s1 {
			if !s.Has(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Diff2 returns s1 - s2.
func Diff2[K comparable, V any](s1, s2 iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		s := set.OfSeq(Keys(s2))
		for k, v := range s1 {
			if !s.Has(k) {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// SymmetricDiff gets the symmetric difference of two sets and gives a set of elements, which are in either of the sets and not in their intersection.
func SymmetricDiff[V comparable](s1, s2 iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		set1 := set.OfSeq(s1)
		set2 := set.OfSeq(s2)
		for v := range s1 {
			if !set2.Has(v) {
				if !yield(v) {
					return
				}
			}
		}
		for v := range s2 {
			if !set1.Has(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// SymmetricDiff2 gets the symmetric difference of two sets and gives a set of elements, which are in either of the sets and not in their intersection.
func SymmetricDiff2[K comparable, V any](s1, s2 iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		set1 := set.OfSeq(Keys(s1))
		set2 := set.OfSeq(Keys(s2))
		for k, v := range s1 {
			if !set2.Has(k) {
				if !yield(k, v) {
					return
				}
			}
		}
		for k, v := range s2 {
			if !set1.Has(k) {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Intersect returns m1 values that keys is contained in m2.
func Intersect[V comparable](s1, s2 iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		s := set.OfSeq(s2)
		for v := range s1 {
			if s.Has(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Intersect2 returns m1 values that keys is contained in m2.
func Intersect2[K comparable, V any](s1, s2 iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		s := set.OfSeq(Keys(s2))
		for k, v := range s1 {
			if s.Has(k) {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Union sequences into one. Duplicates is removed. Order of items in sequence is preserved.
func Union[V comparable](ss ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, s := range ss {
			set := set.Set[V]{}
			for v := range s {
				if set.TryAdd(v) {
					if !yield(v) {
						return
					}
				}
			}
		}
	}
}

// Union sequences into one. Duplicates is removed. Order of items in sequence is preserved.
func Union2[K comparable, V any](ss ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, s := range ss {
			set := set.Set[K]{}
			for k, v := range s {
				if set.TryAdd(k) {
					if !yield(k, v) {
						return
					}
				}
			}
		}
	}
}

// Count values.
func Count[V any](s iter.Seq[V], f func(V) bool) int {
	var count int
	for v := range s {
		if f(v) {
			count++
		}
	}
	return count
}

// Count2 counts values.
func Count2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) int {
	var count int
	for k, v := range s {
		if f(k, v) {
			count++
		}
	}
	return count
}

// Group group sequence by key.
func Group[K comparable, V any](seq iter.Seq2[K, V]) map[K][]V {
	m := make(map[K][]V)
	for k, v := range seq {
		m[k] = append(m[k], v)
	}
	return m
}

// GroupFunc group sequence by key.
func GroupFunc[K comparable, V any](seq iter.Seq[V], key func(V) K) map[K][]V {
	return Group(WithKeys(seq, key))
}

// One return a sequence with one element.

func One[V any](v V) iter.Seq[V] {
	return func(yield func(V) bool) {
		_ = yield(v)
	}
}

// One return a sequence with one element.
func One2[K, V any](k K, v V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		_ = yield(k, v)
	}
}
