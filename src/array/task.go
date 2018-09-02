package array

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	toReturn := make([][]uint8, dy)
	for i := range toReturn {
		slice := make([]uint8, dx)
		for j := range slice {
			slice[j] = uint8(i*2 + j*3)
		}
		toReturn[i] = slice
	}
	return toReturn
}

func Experiment2() {
	pic.Show(Pic)
	//fmt.Println(Pic(10, 10))
}
