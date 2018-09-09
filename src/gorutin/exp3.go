package gorutin

import (
	"fmt"
)

/*
Buffered channel
- blocking write to channel when buffer is full
- blocking read from channel when buffer is empty
*/
func Experiment3() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	//fatal error: all goroutines are asleep - deadlock!
	// as buffer is full and the program can't continue execution
	//ch <- 3
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
