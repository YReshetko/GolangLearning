package interf

import (
	"fmt"
)

/*
Experiment7 - cast empty interface to particular type
*/
func Experiment7() {
	var i interface{}

	i = "Hello"
	s := i.(string)
	fmt.Println(s)

	s1, ok := i.(string)
	fmt.Println(s1, ok)

	f, ok := i.(float64)
	fmt.Println(f, ok)

	f1 := i.(float64) //panic
	fmt.Println(f1)
}
