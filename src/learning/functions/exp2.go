package functions

import (
	"fmt"
)

func closureFunc() func(int) int {
	sum := 0
	return func(i int) int {
		sum = sum + i
		return sum
	}
}

/*
	Closure func cacnge the valiable which is determined in scope closureFunc
*/
func Experiment2() {
	pos, neg := closureFunc(), closureFunc()
	for i := 0; i < 10; i++ {
		fmt.Printf("Pos sum value: %d; Neg sum value: %d\n", pos(i), neg(-2*i))
	}
}
