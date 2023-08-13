package types

import (
	"reflect"
)

// SameType checks that values has the same type. It uses the reflection.
func SameType(v1, v2 any) bool {
	return reflect.ValueOf(v1).Type() == reflect.ValueOf(v2).Type()
}

// IsType checks type.
func IsType[T any](v any) bool {
	_, ok := v.(T)
	return ok
}

// IsNil checks that the interface contains a nil value.
func IsNil(i any) bool {
	if i == nil {
		return true
	}

	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		return true
	}
	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}
