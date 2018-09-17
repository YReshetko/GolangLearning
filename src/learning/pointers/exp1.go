package main

import (
	"fmt"
)

func main() {
	a := uint32(10)
	increment(&a)
	fmt.Println("a =", a, ", Expected 11")
}

func increment(a *uint32) {
	*a = *a + 1
}
