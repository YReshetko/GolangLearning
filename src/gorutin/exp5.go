package gorutin

import (
	"fmt"
)

/*
select allows to wait several commands which sent to channels
select is blocked when one of the cases can not be executed
when 2 or more cases can be executed the case is chosen randomly
*/
func Experiment5() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("Quit")
			return
		}
	}
}
