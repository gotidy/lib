package conditions

import (
	"fmt"
	"strconv"
)

func ExampleIf() {
	fmt.Println(If(false, "true", "false"))
	fmt.Println(If(true, "true", "false"))

	// Output:
	// false
	// true
}

func ExampleIfFunc() {
	fmt.Println(IfFunc(false, func() string { return "true" }, F("false")))
	fmt.Println(IfFunc(true, func() string { return "true" }, F("false")))

	// Output:
	// false
	// true
}

func ExampleF1() {
	for i := 0; i < 5; i++ {
		fmt.Println(IfFunc(i != 0, F1(strconv.Itoa, i), F("zero")))
	}

	// Output:
	// zero
	// 1
	// 2
	// 3
	// 4
}
