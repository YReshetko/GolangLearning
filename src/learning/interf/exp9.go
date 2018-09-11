package interf

import (
	"fmt"
)

type Person struct {
	name string
	age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%v (%v years)", p.name, p.age)
}

/*
Experiment9 - the most popular fmt interface is Stringer,
it's used into almost all methods which print the object somewhere
type Stringer interface {
    String() string
}
*/
func Experiment9() {
	p1 := Person{"Yury", 31}
	fmt.Println(p1)
}
