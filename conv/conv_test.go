package conv

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBytesToString(t *testing.T) {
	tests := [][]byte{
		[]byte("hello"),
		[]byte(""),
		nil,
	}
	for _, tt := range tests {
		if got := BytesToStr(tt); got != string(tt) {
			t.Errorf("BytesToString() = %v, want %v", got, tt)
		}
	}
}

func TestStringToBytes(t *testing.T) {
	tests := []string{
		"hello",
		"",
	}
	for _, tt := range tests {
		want := []byte(tt)
		if got := StrToBytes(tt); !(reflect.DeepEqual(got, want) || len(got) == 0 && len(want) == 0) {
			t.Errorf("StringToBytes() = %v, want %v", got, want)
		}
	}
}

func ExampleSliceTypeCast() {
	s := []uint64{10, 5, 1, 0}
	fmt.Println(SliceTypeCast[uint64, int64](s))

	// Output:
	// [10 5 1 0]
}

func TestSliceTypeCast_Panic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("want panic")
		}
	}()
	s := []uint64{10, 5, 1, 0}
	fmt.Println(SliceTypeCast[uint64, int16](s))
}

func ExampleTypeCast() {
	var s uint64 = 10
	fmt.Println(TypeCast[uint64, int64](s))

	// Output:
	// 10
}

func TestTypeCast_Panic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("want panic")
		}
	}()
	var s uint64 = 10
	fmt.Println(TypeCast[uint64, int16](s))
}
