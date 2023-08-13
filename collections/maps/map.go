package maps

import (
	"github.com/gotidy/lib/collections/set"
)

// Has returns true if m1 contains key.
func Has[K comparable, V any](m map[K]V, key K) bool {
	_, exists := m[key]
	return exists
}

// Clone map.
func Clone[K comparable, V any](m map[K]V) map[K]V {
	result := make(map[K]V, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

// Diff returns m1 - m2.
func Diff[K comparable, V1, V2 any](m1 map[K]V1, m2 map[K]V2) map[K]V1 {
	result := make(map[K]V1)
	for k, v := range m1 {
		if !Has(m2, k) {
			result[k] = v
		}
	}
	return result
}

// DiffKeys returns m1 - m2.
func DiffKeys[K comparable, V any](m map[K]V, keys []K) map[K]V {
	result := make(map[K]V)
	s := set.Of(keys...)
	for k, v := range m {
		if !s.Has(k) {
			result[k] = v
		}
	}
	return result
}

// SymmetricDiff gets the symmetric difference of two sets and gives a set of elements, which are in either of the sets and not in their intersection.
func SymmetricDiff[K comparable, V any](m1, m2 map[K]V) map[K]V {
	result := Clone(m1)
	for k, v := range m2 {
		if Has(m1, k) {
			delete(result, k)
		} else {
			result[k] = v
		}
	}
	return result
}

// Union returns m1 + m2. m1 values are preferred.
func Union[K comparable, V any](m1, m2 map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m2 {
		result[k] = v
	}
	for k, v := range m1 {
		result[k] = v
	}
	return result
}

// Merge maps, last values overwrite first.
func Merge[K comparable, T any](items ...map[K]T) map[K]T {
	result := make(map[K]T, len(items[0])*len(items))
	for _, item := range items {
		for k, t := range item {
			result[k] = t
		}
	}
	return result
}

// Intersect returns m1 values that keys is contained in m2.
func Intersect[K comparable, V any](m1, m2 map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m1 {
		if Has(m2, k) {
			result[k] = v
		}
	}
	return result
}

// Each iterates through map values.
func Each[K comparable, V any](m map[K]V, f func(k K, v V)) {
	for k, v := range m {
		f(k, v)
	}
}

// Keys returns keys of the map.
func Keys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// Values returns values of the map.
func Values[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

// Filter filters values from a map using a filter function.
// It returns a new map with only the elements of s
// for which f returned true.
func Filter[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if f(k, v) {
			result[k] = v
		}
	}
	return result
}

// Map turns a map[K]V to a []T using a mapping function.
func Map[K comparable, V, T any](m map[K]V, f func(K, V) T) []T {
	result := make([]T, 0, len(m))
	for k, v := range m {
		result = append(result, f(k, v))
	}
	return result
}

// Reduce reduces a map[K]V to a single value using a reduction function.
func Reduce[K comparable, V, R any](m map[K]V, initializer R, f func(R, K, V) R) R {
	r := initializer
	for k, v := range m {
		r = f(r, k, v)
	}
	return r
}

// Append one map to another map.
func Append[K comparable, V any](dst, src map[K]V) map[K]V {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
