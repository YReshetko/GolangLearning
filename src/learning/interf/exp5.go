package interf

import "fmt"

type I3 interface {
	method()
}

/*
Experiment5 - interface with nil value doesn't contain type so we cant invoke interface method
and this is a cause of Runtime Error
*/
func Experiment5() {
	var i I3
	describe1(i)
	i.method() // Runtime error
}

func describe1(i I2) {
	fmt.Printf("(%v, %T)\n", i, i)
}
