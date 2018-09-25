package save

import (
	"bufio"
	"bytes"
	"mailclient/domain"
	"os"
)

var (
//dao = NewDao("acad", "calls")
)

type s struct {
	storagePath string
}

type EmailSaver interface {
	Save(toSave *domain.EmailToSave) error
}

func EmailSaverInstance(storagePath string) EmailSaver {
	return &s{storagePath}
}

func (s *s) Save(toSave *domain.EmailToSave) error {
	err := s.saveFile(toSave.Buffer, toSave.EmailData.RecordFileName)
	if err == nil {
		//s.dao.save(toSave.EmailData)
	}
	return err
}
func (s *s) saveFile(buffer *bytes.Buffer, fileName string) error {
	// open output file
	fo, err := os.Create(s.storagePath + fileName)
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
	buffer.WriteTo(w)

	if err = w.Flush(); err != nil {
		return err
	}
	return nil
}
