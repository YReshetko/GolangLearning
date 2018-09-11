package gorutin

import (
	"fmt"
	"time"
)

/*
select {
case i := <-c:
    // use i to do some action
default:
	// recieving from c is blocked so we can send some information to another chennal
	// used in non-blocking algorithms
}
*/
func Experiment6() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Print("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)

		}
	}
}
