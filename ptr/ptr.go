package ptr

import "github.com/gotidy/lib/constraints"

// Of returns pointer to value.
func Of[T any](v T) *T {
	return &v
}

// OfOrNil returns pointer to value or nil if value is a zero value.
func OfOrNil[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

// Value returns the value of the pointer passed in or the default value if the pointer is nil.
func Value[T any](v *T) T {
	if v == nil {
		return Zero[T]()
	}
	return *v
}

// ValueDef returns the value of the pointer passed in or default value if the pointer is nil.
func ValueDef[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}

// Zero returns a value representing the zero value for the specified type.
func Zero[T any]() T {
	var zero T
	return zero
}

// Equal compare two pointers.
func Equal[T comparable](v1, v2 *T) bool {
	return v1 == v2 || v1 != nil && v2 != nil && *v1 == *v2
}

// ConvertNumber  the pointer of number of one type to another type.
func ConvertNumber[T, V constraints.Number](t *T) *V {
	if t == nil {
		return nil
	}
	v := V(*t)
	return &v
}
