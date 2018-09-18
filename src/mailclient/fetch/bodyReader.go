package fetch

import (
	"fmt"
	"io"
	"io/ioutil"
	"mailclient/domain"
	"regexp"

	"github.com/emersion/go-message/mail"
)

type emailReader struct {
	emailBodyRegexp    *regexp.Regexp
	attachedFileRegexp *regexp.Regexp
}

type EmailReader interface {
	ReadEmail(reader mail.Reader, uid uint32) (domain.EmailToSave, bool)
}

func NewEmailReader(bodyRegexp, fileRegexp string) EmailReader {
	return &emailReader{
		regexp.MustCompile(bodyRegexp),
		regexp.MustCompile(fileRegexp),
	}
}

func (emailReader *emailReader) ReadEmail(reader mail.Reader, uid uint32) (domain.EmailToSave, bool) {
	emailToSave := domain.EmailToSave{}
	emailData := domain.EmailData{Uid: uid}
	foundTextData := false
	foundAttachedFile := false
	for {
		p, err := reader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		switch h := p.Header.(type) {
		case mail.TextHeader:
			// This is the message's text (can be plain-text or HTML)
			b, _ := ioutil.ReadAll(p.Body)
			partText := string(b)
			if emailReader.matchText(partText) {
				// TODO extract groups of matched text and save it in structure
				foundTextData = true
			}
		case mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			if emailReader.matchFileName(filename) {
				emailData.RecordFileName = filename
				emailToSave.Reader = p.Body
				foundAttachedFile = true
			}

		}
	}
	emailToSave.EmailData = emailData
	return emailToSave, foundAttachedFile && foundTextData
}

func (emailReader *emailReader) matchText(text string) bool {
	return emailReader.emailBodyRegexp.MatchString(text)
}
func (emailReader *emailReader) matchFileName(text string) bool {
	return emailReader.attachedFileRegexp.MatchString(text)
}
