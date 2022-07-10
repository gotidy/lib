package ptr

// Of returns pointer to value.
func Of[T any](v T) *T {
	return &v
}

// Value returns the value of the pointer passed in or the default value if the pointer is nil.
func Value[T any](v *T) T {
	if v == nil {
		return Zero[T]()
	}
	return *v
}

// ValueDef returns the value of the int pointer passed in or default value if the pointer is nil.
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
