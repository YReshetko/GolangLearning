package interf

import (
	"fmt"
)

type I1 interface {
	method()
}

type T1 struct {
	S string
}

type MyFloat1 float64

func (t *T1) method() {
	fmt.Println(t.S)
}

func (fl MyFloat1) method() {
	fmt.Println(fl)
}

/*
Experiment3 - The value of interface is a pair (value, type) where:
- value is a alue of paticular type
- type is real type of value (structure, alias to type etc.)
*/
func Experiment3() {
	var myType I1

	myType = &T1{"Hello world"}
	declaration(myType)
	myType.method()

	myType = MyFloat1(1.234)
	declaration(myType)
	myType.method()
}

func declaration(value I1) {
	fmt.Printf("(%v, %T)\n", value, value)
}
