package gorutin

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"time"
)

func Experiment7() {

	startTime := time.Now()
	for i := 0; i < 300; i++ {
		go goExp()
	}

	fmt.Println("Time spent to retrieve gorutine ids:", time.Now().Sub(startTime))
	//go goExp()
	//go goExp()
	//go goExp()
	//go goExp()
	time.Sleep(10 * time.Second)
}

func goExp() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
