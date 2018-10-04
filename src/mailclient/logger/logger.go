package logger

import (
	"bytes"
	"fmt"
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

var logLevel uint
var templateMap map[string]string

func Init() func() {
	templateMap = make(map[string]string)
	f, err := os.OpenFile("log/mfetch.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return func() {
		f.Close()
	}
}
func SetLogLevel(value uint) {
	if value > DEBUG {
		value = 5
	}
	logLevel = value
}
func Debug(msg string, args ...interface{}) {
	if logLevel == DEBUG {
		print("DEBUG", msg, args)
	}
}

func Info(msg string, args ...interface{}) {
	if logLevel >= INFO {
		print("INFO ", msg, args)
	}
}

func Warning(msg string, args ...interface{}) {
	if logLevel >= WARNING {
		print("WARN ", msg, args)
	}
}

func Error(msg string, args ...interface{}) {
	if logLevel >= ERROR {
		print("ERROR", msg, args)
	}
}

func Fatal(msg string, args ...interface{}) {
	if logLevel >= FATAL {
		print("FATAL", msg, args)
	}
}

func print(lvl, msg string, args ...interface{}) {
	n, s := getGoRoutinID(4)
	tmp := templateMap[lvl]
	if tmp == "" {
		tmp = getTemplate(lvl)
		templateMap[lvl] = tmp
	}
	message := fmt.Sprintf(msg, args)
	if message[len(message)-1] != '\n' {
		message = message + "\n"
	}
	log.Printf(tmp, n, s, message)
}
func getTemplate(lvl string) string {
	return "[" + lvl + "] [GOROUTINE-ID-%d][%s] %s"
}
func getGoRoutinID(index int) (uint64, string) {
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
		b = b[bytes.IndexByte(b, ':')+1:]
		if count == index {
			s = string(m)
			break
		}
	}
	n, _ := strconv.ParseUint(string(num), 10, 64)
	return n, s
}
