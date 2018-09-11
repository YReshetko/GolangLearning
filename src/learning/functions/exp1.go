package functions

import (
	"fmt"
	"math"
)

func compute(fn func(float64, float64) float64) float64 {
	return fn(20.32432, 30.2312)
}

/*
Experiment1 - functiona as a value
*/
func Experiment1() {
	hypot := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}

	fmt.Println(hypot(5, 10))

	fmt.Println(compute(hypot))
	fmt.Println(compute(math.Pow))
	fmt.Println(compute(func(x, y float64) float64 {
		return x + y
	}))
}
