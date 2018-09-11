package dictionary

import "fmt"

type Vertex1 struct {
	Lat, Long float64
}

var m1 = map[string]Vertex1{
	"Bel Labs": Vertex1{
		20.328746, 76.238746,
	},
	"Google": Vertex1{
		174.2389493, 47.239856,
	},
}

/*
	Map initialization with initial values
*/
func Experiment2() {
	fmt.Println(m1)

	// what if we take value from map and change that
	value := m1["Bel Labs"]
	value.Lat = 0
	fmt.Println(value)
	fmt.Println(m1)
}
