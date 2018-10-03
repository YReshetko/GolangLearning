package logger

import (
	"bytes"
	"log"
	"os"
	"runtime"
	"strconv"
)

const (
	FATAL   = 1
	ERROR   = 2
	WARNING = 3
	INFO    = 4
	DEBUG   = 5
)

func Init() func() {
	f, err := os.OpenFile("log/mfetch.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return func() {
		f.Close()
	}
}
func getGoRoutinID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
