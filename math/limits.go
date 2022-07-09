package math

import "github.com/gotidy/lib/constraints"

// Min returns minimum value.
func Min[T constraints.Ordered](v1, v2 T) T {
	if v1 < v2 {
		return v1
	}
	return v2
}

// Max returns maximum value.
func Max[T constraints.Ordered](v1, v2 T) T {
	if v1 < v2 {
		return v2
	}
	return v1
}

// Between returns true if value is between max and min.
func Between[T constraints.Ordered](v, min, max T) bool {
	return v >= min && v <= max
}

// MustBetween returns true if value is between max and min.
func MustBetween[T constraints.Ordered](v, min, max T) T {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
