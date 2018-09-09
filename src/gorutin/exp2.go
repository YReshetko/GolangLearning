package gorutin

import (
	"fmt"
)

func Experiment2() {
	arr := []int{1, -4, 5, 8, -2, 4}

	ch := make(chan int)
	go sum(arr[:len(arr)/2], ch)
	go sum(arr[len(arr)/2:], ch)

	// Default implementation
	//x, y := <-ch, <-ch // receive values from chennal
	//fmt.Println(x, y, x+y)

	// modification
	// Alternative implementation which shows
	//How many times we write into chennal so many times we must read from channel
	sum := 0
	for i := 0; i < 2; i++ {
		sum += <-ch
	}
	fmt.Println(sum)

}

func sum(arr []int, ch chan int) {
	sum := 0
	for _, val := range arr {
		sum += val
	}
	ch <- sum
}
