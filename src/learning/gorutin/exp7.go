package gorutin

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"time"
)

func Experiment7() {
	go goExp()
	go goExp()
	go goExp()
	go goExp()
	go goExp()
	time.Sleep(10 * time.Second)
}

func goExp() {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	fmt.Println("gorutine ID:", n)
}
