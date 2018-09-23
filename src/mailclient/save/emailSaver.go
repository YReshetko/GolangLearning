package save

import (
	"bufio"
	"io"
	"mailclient/domain"
	"os"
)

var (
	dao EmailDao = NewDao("acad", "calls")
)

func Save(toSave *domain.EmailToSave) {

	err := saveFile(toSave.Reader, toSave.EmailData.RecordFileName)
	if err == nil {
		dao.save(toSave.EmailData)
	}

}
func saveFile(messageReader io.Reader, fileName string) error {
	// open output file
	fo, err := os.Create("D:/recordStorage/" + fileName)
	if err != nil {
		return err
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	// make a write buffer
	w := bufio.NewWriter(fo)

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := messageReader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := w.Write(buf[:n]); err != nil {
			return err
		}
	}

	if err = w.Flush(); err != nil {
		return err
	}
	return nil
}
