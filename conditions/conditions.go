package conditions

// If returns the value depend to the condition.
func If[T any](condition bool, whenTrue, whenFalse T) T {
	if condition {
		return whenTrue
	}
	return whenFalse
}

// IfFunc returns the value depend to the condition.
func IfFunc[T any](condition bool, whenTrue func() T, whenFalse func() T) T {
	if condition {
		return whenTrue()
	}
	return whenFalse()
}

// F wrap the value with a func.
func F[T any](t T) func() T {
	return func() T { return t }
}

// F1 wrap the function with a func.
func F1[R, T any](f func(t T) R, t T) func() R {
	return func() R { return f(t) }
}

// F2 wrap the function with a func.
func F2[R, T1, T2 any](f func(t1 T1, t2 T2) R, t1 T1, t2 T2) func() R {
	return func() R { return f(t1, t2) }
}

// F3 wrap the function with a func.
func F3[R, T1, T2, T3 any](f func(t1 T1, t2 T2, t3 T3) R, t1 T1, t2 T2, t3 T3) func() R {
	return func() R { return f(t1, t2, t3) }
}

// F4 wrap the function with a func.
func F4[R, T1, T2, T3, T4 any](f func(t1 T1, t2 T2, t3 T3, t4 T4) R, t1 T1, t2 T2, t3 T3, t4 T4) func() R {
	return func() R { return f(t1, t2, t3, t4) }
}

// F5 wrap the function with a func.
func F5[R, T1, T2, T3, T4, T5 any](f func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) R, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) func() R {
	return func() R { return f(t1, t2, t3, t4, t5) }
}
