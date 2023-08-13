package ptr

import (
	"fmt"

	"github.com/gotidy/lib/constraints"
)

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

// EqualValue compare the value of pointer with the value.
func EqualValue[T comparable](p *T, v T) bool {
	return p != nil && *p == v
}

// ConvertNumber  the pointer of number of one type to another type.
func ConvertNumber[T, V constraints.Number](t *T) *V {
	if t == nil {
		return nil
	}
	v := V(*t)
	return &v
}

// Copy returns a copy of the pointered data.
func Copy[T any](t *T) *T {
	r := *t
	return &r
}

// Empty checks that value is nil or equal zero value.
func Empty[T comparable](t *T) bool {
	return t == nil || *t == Zero[T]()
}

// Str return value or empty string if value if nil.
func Str[T any](v *T) string {
	if v == nil {
		return ""
	}
	i := interface{}(v)
	if i, ok := i.(interface{ String() string }); ok {
		return i.String()
	}
	return fmt.Sprintf("%v", *v)
}

// StrDef return value or default string if value if nil.
func StrDef[T any](v *T, def string) string {
	if v == nil {
		return def
	}
	i := interface{}(v)
	if i, ok := i.(interface{ String() string }); ok {
		return i.String()
	}
	return fmt.Sprintf("%v", *v)
}
