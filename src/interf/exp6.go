package interf

import "fmt"

/*
Experiment6 - empty interface (interface{}) describes any type
*/
func Experiment6() {
	var i interface{}
	describe2(i)

	i = 43
	describe2(i)

	i = "hello"
	describe2(i)

}

func describe2(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}
