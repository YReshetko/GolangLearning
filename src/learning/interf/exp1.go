package interf

import (
	"fmt"
	"math"
)

type Abser interface {
	Abs() float64
}

type MyFloat float64

type vertex struct {
	X, Y float64
}

func (v MyFloat) Abs() float64 {
	if v < 0 {
		return float64(-v)
	}
	return float64(v)
}

func (v *vertex) Abs() float64 {
	//func (v vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func Experiment1() {
	var a Abser
	f := MyFloat(-math.Sqrt2)
	v := vertex{3, 4}

	a = f  // a MyFloat implements Abser
	a = &v // a *Vertex implements Abser

	// In the following line, v is a Vertex (not *Vertex)
	// and does NOT implement Abser.
	// a = v // WILL NOT BE COMPILED!!!!!!!!!!!

	fmt.Println(a.Abs())
}
