package gorutin

import (
	"fmt"
	"time"
)

func Experiment1() {
	go say("world")
	say("Hello")
}

func say(word string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Microsecond)
		fmt.Println(word)
	}

}
