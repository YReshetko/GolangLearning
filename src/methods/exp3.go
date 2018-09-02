package methods

import (
	"fmt"
	"math"
)

type vertex2 struct {
	X, Y float64
}

func (v vertex2) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *vertex2) Scale(value int) {
	v.X = v.X * float64(value)
	v.Y = v.Y * float64(value)
}

func (v vertex2) Scale2(value int) {
	v.X = v.X * float64(value)
	v.Y = v.Y * float64(value)
}

/*
Experiment3 - shows different when method declared for type or pointer to type
- when method is declared on type, the it works with copy of value of type
- when method is declared ot type pointer then it changes directly the value on which method was called
*/
func Experiment3() {
	vert := vertex2{2.10, 1.20}
	vert.Scale(2)
	//vert.Scale2(2)
	fmt.Printf("Vertex: %+v\n", vert)
	fmt.Printf("Abs: %g\n", vert.Abs())
}
