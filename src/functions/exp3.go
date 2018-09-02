package functions

import (
	"fmt"
	"sync"
)

/*
	Closure func cacnge the valiable which is determined in scope closureFunc
*/
var wg sync.WaitGroup

func Experiment3() {
	pos, neg := closureFunc(), closureFunc()
	wg.Add(1)
	go corutFunc("Pos index value: %d\n", 1, pos)
	wg.Add(1)
	go corutFunc("Neg index value: %d\n", -2, neg)
	wg.Wait()
}

func corutFunc(message string, index int, fn func(int) int) {
	for i := 0; i < 10; i++ {
		fmt.Printf(message, fn(index*i))
	}
	wg.Done()
}
