package functions

import (
	"fmt"
)

func fibonacci() func() int {
	first := 0
	second := 1
	return func() int {
		new := first + second
		first = second
		second = new
		return new
	}
}

func Task() {
	f := fibonacci()
	for result := f(); result < 1000000000; result = f() {
		fmt.Println(result)
	}
}
