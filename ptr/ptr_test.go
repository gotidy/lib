package ptr

import (
	"fmt"
	"reflect"
	"testing"
)

func equal(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, actual %#v", expected, actual)
	}
}

func TestOf(t *testing.T) {
	equal(t, int(10), *Of(10))
}

func TestValue(t *testing.T) {
	equal(t, int(10), Value(Of(10)))
	equal(t, int(0), Value((*int)(nil)))
}

func TestValueDef(t *testing.T) {
	equal(t, int(10), ValueDef(Of(10), 0))
	equal(t, int(5), ValueDef(nil, 5))
}

func ExampleEqual() {
	v1 := 1
	v2 := 1
	v3 := 10

	fmt.Println(Equal((*int)(nil), (*int)(nil)))
	fmt.Println(Equal(&v1, &v1))
	fmt.Println(Equal((*int)(nil), &v1))
	fmt.Println(Equal(&v1, &v2))
	fmt.Println(Equal(&v1, &v3))

	// Output:
	// true
	// true
	// false
	// true
	// false
}

func ExampleCoalesce() {
	fmt.Println(Value(Coalesce[int](nil, Of(1), Of(2))))
	fmt.Println(Value(Coalesce[int](Of(2), Of(1), nil)))
	fmt.Println(Value(Coalesce[int](Of(2))))
	fmt.Println(Coalesce[int](nil, nil))
	fmt.Println(Coalesce[int]())

	// Output:
	// 1
	// 2
	// 2
	// <nil>
	// <nil>
}
