package methods

import (
	"fmt"
	"math"
)

type vertex struct {
	Lat, Long float64
}

func (v vertex) Abs() float64 {
	return math.Sqrt(v.Lat*v.Lat + v.Long*v.Long)
}

/*
Experiment1 - apply methods to structures
*/
func Experiment1() {
	vert := vertex{9.0, 12.0}
	fmt.Printf("vertex %+v has abs %d\n", vert, vert.Abs())
}
