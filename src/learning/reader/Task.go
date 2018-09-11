package reader

import "golang.org/x/tour/reader"

type MyReader struct{}

func (reader MyReader) Read(b []byte) (int, error) {
	var i int
	for i = range b {
		b[i] = 'A' // symbol "A"
	}
	return i, nil
}

func Task() {
	reader.Validate(MyReader{})
}
