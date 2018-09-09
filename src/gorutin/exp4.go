package gorutin

import (
	"fmt"
)

/*
close(ch) - closes channel - can be closed only by producer
v, ok := <-ch - returns 2 values: value and status if channel closed

the status is used in val := range ch to determine when range could be stopped
*/

func Experiment4() {
	ch := make(chan int, 100)
	go fibonachchi(cap(ch), ch)
	for val := range ch {
		fmt.Println(val)
	}
}

func fibonachchi(n int, ch chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		ch <- x
		x, y = y, x+y
	}
	close(ch)
}
