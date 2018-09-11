package methods

import (
	"fmt"
	"math"
)

/*
MyFloat - the same as float 64
*/
type MyFloat float64

/*
Abs - the equal to math.Abs function
*/
func (v MyFloat) Abs() MyFloat {
	if v < 0 {
		return MyFloat(-v)
	}
	return MyFloat(v)
}

/*
Experiment2 - methods can not be only applied to structures. The main point is that type has to be declared in the same package wher method
*/
func Experiment2() {
	val := MyFloat(-math.Sqrt2)
	fmt.Println(val)
	fmt.Println(math.Abs(float64(val)))
	fmt.Println(val.Abs())

}
