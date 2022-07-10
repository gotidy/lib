package ptr

import (
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
