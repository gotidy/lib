package math

import (
	"fmt"
	"testing"
)

func ExampleMin(t *testing.T) {
	fmt.Println(Min(1, 10))
	fmt.Println(Min(11, 10))
	// Output:
	// 1
	// 10
}

func ExampleMax(t *testing.T) {
	fmt.Println(Max(1, 10))
	fmt.Println(Max(11, 10))
	// Output:
	// 10
	// 11
}

func ExampleMustBetween(t *testing.T) {
	fmt.Println(MustBetween(1, 10, 100))
	fmt.Println(MustBetween(11, 10, 100))
	fmt.Println(MustBetween(110, 10, 100))
	// Output:
	// 10
	// 11
	// 100
}

func ExampleBetween(t *testing.T) {
	fmt.Println(Between(1, 10, 100))
	fmt.Println(Between(110, 10, 100))
	fmt.Println(Between(11, 10, 100))
	// Output:
	// false
	// false
	// true
}
