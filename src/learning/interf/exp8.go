package interf

import "fmt"

func do(val interface{}) {
	switch v := val.(type) {
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}

/*
Experiment8 - check type into switch
switch v := i.(type) {
case T:
    // здесь v имеет тип T
case S:
    // здесь v имеет тип S
default:
    // нет совпадения; здесь v имеет такой же тип, что и i
}
*/
func Experiment8() {
	do(21)
	do("hello")
	do(true)
	do(21.0)
}
