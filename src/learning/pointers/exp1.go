package main

import (
	"fmt"
)

type MyStruct struct {
	field1 string
	Field2 string
}

func main() {
	intExperiment()
	stringExperimant()
	structExperiment()
}

func structExperiment() {
	value := MyStruct{"Hello", "world"}
	fmt.Printf("Struct:%+v\n", value)
	fmt.Printf("Struct:%+v\n", &value.field1)
	changeStruct(&value)
	fmt.Printf("Struct:%+v\n", value)
	fmt.Printf("Struct:%+v\n", &value.field1)
}

func changeStruct(value *MyStruct) {
	value.field1 = "Good bye"
	value.Field2 = "Space"
}

func stringExperimant() {
	str := "Hello"
	changeString(&str)
	fmt.Println(str)
}

func changeString(str *string) {
	a := "hello"
	*str = a
}

func intExperiment() {
	a := uint32(10)
	increment(&a)
	fmt.Println("a =", a, ", Expected 11")
}
func increment(a *uint32) {
	*a = *a + 1
}
