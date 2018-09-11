package dictionary

import (
	"fmt"
)

type Vertex struct {
	Lat, Long float64
}

var m map[string]Vertex

/**
Base example of empty map initialization

Additionaly check how we can retrieve non existing value and compare with empty value
*/
func Experiment1() {
	nilVertex := Vertex{0, 0}
	m = make(map[string]Vertex)
	m["Some building"] = Vertex{10.234324, 134.124235}
	fmt.Println(m["Some building"])
	fmt.Println(m["Non existing key"])
	ifNilVertex := nilVertex == m["Non existing key"]
	fmt.Printf("Is nil vertex: %t\n", ifNilVertex)

	unk := m["Non existing key"]
	var returnedNilVertex *Vertex = &unk
	(*returnedNilVertex).Lat = 100.3432453

	fmt.Println(m["Non existing key"])
	ifNilVertex = nilVertex == m["Non existing key"]
	fmt.Printf("Is nil vertex: %t\n", ifNilVertex)
	fmt.Printf("nil vertex: %+v", nilVertex)

}
