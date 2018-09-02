package interf

import (
	"fmt"
)

type I interface {
	M()
}

type T struct {
	S string
}

func (t T) M() {
	fmt.Println(t.S)
}

func Experiment2() {
	var value I = T{"Hello world"}

	m := make(map[string]I)
	m["Greating"] = T{"Hello"}
	m["Persone"] = T{"Yury"}

	value.M()

	for key, value := range m {
		fmt.Printf("Key:%s, value:%s", key, value)
	}
}
