// Contains conversion helpers.
package conv

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/gotidy/lib/constraints"
)

// StrToBytes converts fast the string to the byte slice.
// Source string will mutate if change bytes in the byte slice.
func StrToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// BytesToStr converts fast the byte slice to the string.
// Result string will mutate if change bytes in the byte slice.
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// SliceTypeCast casts the slice type.
// Use extremely carefully and only in those cases where the types have the same dimension.
// Size of types must be equal.
func SliceTypeCast[T, K any](s []T) []K {
	var t T
	var k K
	if unsafe.Sizeof(t) != unsafe.Sizeof(k) {
		panic(fmt.Errorf("different sizes of %T and %T", t, k))
	}
	return *(*[]K)(unsafe.Pointer(&s))
}

// TypeCast casts the type.
// Use extremely carefully. Size of types must be equal.
func TypeCast[T, K any](s T) K {
	var t T
	var k K
	if unsafe.Sizeof(t) != unsafe.Sizeof(k) {
		panic(fmt.Errorf("different sizes of %T and %T", t, k))
	}
	return *(*K)(unsafe.Pointer(&s))
}

// Itoa returns the string representation of i.
func Itoa[T constraints.Signed](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

// Uitoa returns the string representation of i.
func Uitoa[T constraints.Unsigned](i T) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Atoi returns the string representation of i.
func Atoi[T constraints.Signed](s string) (T, error) {
	var t T
	i, err := strconv.ParseInt(s, 10, int(unsafe.Sizeof(t))*8)
	if err != nil {
		return 0, err
	}
	return T(i), nil
}

// Atoui parses the string to unsigned integer type.
func Atoui[T constraints.Unsigned](s string) (T, error) {
	var t T
	i, err := strconv.ParseUint(s, 10, int(unsafe.Sizeof(t))*8)
	if err != nil {
		return 0, err
	}
	return T(i), nil
}
