package reader

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	reader io.Reader
}

func (rot13 rot13Reader) Read(b []byte) (int, error) {
	i, err := rot13.reader.Read(b)
	for index, val := range b {
		if isInUpperCase(val) {
			b[index] = rotate13(val, 'N')
		}
		if isInLowerCase(val) {
			b[index] = rotate13(val, 'n')
		}

	}
	return i, err
}
func rotate13(b, midValue byte) byte {
	if b < midValue {
		return b + 13
	} else {
		return b - 13
	}
}
func isInUpperCase(b byte) bool {
	return b <= 'Z' && b >= 'A'
}
func isInLowerCase(b byte) bool {
	return b <= 'z' && b >= 'a'
}

func Task2() {
	sr := strings.NewReader("Lbh penpxrq gur pbqr!")
	rot := rot13Reader{sr}
	io.Copy(os.Stdout, &rot)
}
