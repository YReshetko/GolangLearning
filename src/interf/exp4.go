package interf

import (
	"fmt"
)

type I2 interface {
	method()
}

type T2 struct {
	S string
}

func (t *T2) method() {
	if t == nil {
		fmt.Println("<nil>")
		return
	}
	fmt.Println(t.S)
}

func Experiment4() {
	var i I2
	var t *T2

	i = t

	describe(i)
	i.method()

	i = &T2{"Hello"}

	describe(i)
	i.method()
}
func describe(i I2) {
	fmt.Printf("(%v, %T)\n", i, i)
}
