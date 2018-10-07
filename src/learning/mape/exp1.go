package mape

import (
	"fmt"
)

type mapKey struct {
	Field1 string
	field2 string
}

func Experiment1() {
	//testWithPointerKeys()
	testWithValueKeys()
}
func testWithValueKeys() {
	fmt.Println("Test map with values key")
	myMap := make(map[mapKey]string)

	key1 := mapKey{"key1_field1", "key1_field2"}
	key2 := mapKey{"key2_field1", "key2_field2"}

	myMap[key1] = "value1"
	myMap[key2] = "value2"

	fmt.Println("Check with the same keys")
	fmt.Println(myMap[key1])
	fmt.Println(myMap[key2])

	key1.field2 = "1"
	key2.field2 = "2"
	fmt.Println("Check after changed private key fields")
	fmt.Println(myMap[key1])
	fmt.Println(myMap[key2])

	fmt.Println("Check with new created keys")
	fmt.Println("v1:", myMap[mapKey{"key1_field1", "key1_field2"}])
	fmt.Println("v2:", myMap[mapKey{"key2_field1", "key2_field2"}])
}
func testWithPointerKeys() {
	fmt.Println("Test map with pointers key")
	myMap := make(map[*mapKey]string)

	key1 := &mapKey{"key1_field1", "key1_field2"}
	key2 := &mapKey{"key2_field1", "key2_field2"}

	myMap[key1] = "value1"
	myMap[key2] = "value2"

	fmt.Println("Check with the same keys")
	fmt.Println(myMap[key1])
	fmt.Println(myMap[key2])

	fmt.Println("Check with new created keys")
	fmt.Println("v1:", myMap[&mapKey{"key1_field1", "key1_field2"}])
	fmt.Println("v2:", myMap[&mapKey{"key2_field1", "key2_field2"}])
}
