package reader

import (
	"fmt"
	"io"
	"strings"
)

func Experiment1() {
	reader := strings.NewReader("ABCDEFGHIJKLMNOPQASTUVWXYZabcdefghijklmnopqrstuvwxyz")
	b := make([]byte, 26)
	for {
		n, err := reader.Read(b)
		fmt.Printf("n=%v, err=%v, b[]=%v\n", n, err, b)
		fmt.Printf("b[:n]=%q\n", b[:n])
		if err == io.EOF {
			break
		}
	}
}
