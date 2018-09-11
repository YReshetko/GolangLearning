package dictionary

import (
	"fmt"
)

// Change map experimant
func Experiment4() {
	m := make(map[string]Vertex1)

	// adding element
	m["Google"] = Vertex1{34.34523542, 32.32543254}
	fmt.Println(m)

	// Chnage element
	m["Google"] = Vertex1{50.324534354, 100.3243254}
	fmt.Println(m)

	// Get element
	value := m["Google"]
	fmt.Println(value)

	//remove element
	delete(m, "Google")
	fmt.Println(m)

	// check if exist in map
	val, ok := m["Google"]
	fmt.Println("The Google value:", val, "; exists: ", ok)
}
