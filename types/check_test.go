package types

import (
	"fmt"
	"testing"
)

func ExampleSameType() {
	type Obj1 struct{}
	type Obj2 struct{}

	fmt.Println(SameType(Obj1{}, Obj1{}))
	fmt.Println(SameType(Obj1{}, Obj2{}))

	// Output:
	// true
	// false
}

func ExampleIsType() {
	type Obj1 struct{}
	type Obj2 struct{}

	fmt.Println(IsType[Obj1](Obj1{}))
	fmt.Println(IsType[Obj2](Obj1{}))

	// Output:
	// true
	// false
}

func BenchmarkIsNil(b *testing.B) {
	type Obj struct {
	}

	var obj *Obj

	b.Run("is_nil_nil_reflect", func(b *testing.B) {
		var a any = obj
		for i := 0; i < b.N; i++ {
			_ = IsNil(a)
		}
	})

	b.Run("is_nil_nil_native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if obj == nil { //nolint
			}
		}
	})

	obj = &Obj{}

	b.Run("is_nil_reflect", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = IsNil(obj)
		}
	})

	b.Run("is_nil_native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if obj == nil { //nolint
			}
		}
	})
}

func ExampleIsNil() {
	type Obj struct{}

	var obj *Obj
	var a any = obj
	fmt.Println(IsNil(a))

	obj = &Obj{}
	a = obj
	fmt.Println(IsNil(a))

	// Output:
	// true
	// false
}
