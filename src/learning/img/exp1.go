package img

import (
	"fmt"
	"image"
)

func Experiment1() {
	m := image.NewRGBA(image.Rect(0, 0, 100, 100))
	fmt.Println(m.Bounds())
	fmt.Println(m.At(0, 0))
}
