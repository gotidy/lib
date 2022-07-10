package math

import (
	"fmt"
)

func ExampleMin() {
	fmt.Println(Min(1, 10))
	fmt.Println(Min(11, 10))
	// Output:
	// 1
	// 10
}

func ExampleMax() {
	fmt.Println(Max(1, 10))
	fmt.Println(Max(11, 10))
	// Output:
	// 10
	// 11
}

func ExampleMustBetween() {
	fmt.Println(MustBetween(1, 10, 100))
	fmt.Println(MustBetween(11, 10, 100))
	fmt.Println(MustBetween(110, 10, 100))
	// Output:
	// 10
	// 11
	// 100
}

func ExampleBetween() {
	fmt.Println(Between(1, 10, 100))
	fmt.Println(Between(110, 10, 100))
	fmt.Println(Between(11, 10, 100))
	// Output:
	// false
	// false
	// true
}
