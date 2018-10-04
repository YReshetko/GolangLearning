package save

import (
	"bufio"
	"bytes"
	"mailclient/domain"
	"mailclient/logger"
	"os"
)

type saver struct {
	storagePath string
}

/*
EmailSaver - save buffer from EmailToSave into fule
*/
type EmailSaver interface {
	Save(toSave *domain.EmailToSave) error
}

/*
NewEmailSaver - creates new EmailSaver which contains storage path
*/
func NewEmailSaver(storagePath string) EmailSaver {
	return &saver{storagePath}
}

func (s *saver) Save(toSave *domain.EmailToSave) error {
	return s.saveFile(toSave.Buffer, toSave.EmailData.RecordFileName)
}
func (s *saver) saveFile(buffer *bytes.Buffer, fileName string) error {
	fullPath := s.storagePath + fileName
	fo, err := os.Create(fullPath)
	if err != nil {
		logger.Error("Error creating file %s: %v", fullPath, err)
		return err
	}
	defer func() {
		if err := fo.Close(); err != nil {
			logger.Error("Error closing file %s: %v", fullPath, err)
		}
	}()
	w := bufio.NewWriter(fo)
	buffer.WriteTo(w)
	if err = w.Flush(); err != nil {
		logger.Error("Error flushing bytes into file %s: %v", fullPath, err)
		return err
	}
	return nil
}
