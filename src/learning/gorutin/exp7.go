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
		go runStackTest()
	}
	//runStackTest()
	fmt.Println("Time spent to retrieve gorutine ids:", time.Now().Sub(startTime))
	//go goExp()
	//go goExp()
	//go goExp()
	//go goExp()
	time.Sleep(10 * time.Second)
}
func runStackTest() {
	testStackTrace()
}

func testStackTrace() {
	b := make([]byte, 1024)
	runtime.Stack(b, true)
	//fmt.Println(string(b))
	n, s := goExp(3)
	fmt.Printf("[GOROUTINE-ID-%v]\t[%s]\n", n, s)
}
func goExp(index int) (uint64, string) {
	b := make([]byte, 1024)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	num := b[:bytes.IndexByte(b, ' ')]
	b = b[bytes.IndexByte(b, ':')+1:]
	count := 0
	var s string
	for bytes.IndexByte(b, ':') > 0 {
		count++
		b = b[bytes.IndexByte(b, ':')+1:]
		m := b[:bytes.IndexByte(b, ':')]
		//fmt.Println(string(m))
		b = b[bytes.IndexByte(b, ':')+1:]
		if count == index {
			s = string(m)
			break
		}
	}

	n, _ := strconv.ParseUint(string(num), 10, 64)
	return n, s
}
