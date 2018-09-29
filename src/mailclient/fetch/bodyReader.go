package fetch

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mailclient/config"
	"mailclient/domain"
	"regexp"
	"time"

	"github.com/emersion/go-message/mail"
	"github.com/grokify/gotilla/time/timeutil"
)

const (
	attachedFileName = "attached-file-name"
	whoCalls         = "who-calls"
	inputNumber      = "input-number"
	participant      = "participant"
	callLength       = "call-length"

	dateTimePattern = "%s-%s-%s %s:%s:%s"
)

var monthes = map[string]string{
	"Jan": "01",
	"Feb": "02",
	"Mar": "03",
	"Apr": "04",
	"May": "05",
	"Jun": "06",
	"Jul": "07",
	"Aug": "08",
	"Sep": "09",
	"Oct": "10",
	"Nov": "11",
	"Dec": "12",
}

type emailReader struct {
	regExpMap map[string]*regexp.Regexp
}

/*
EmailReader - read particular email, match to ecpected patterns and returns EmailToSave structure
*/
type EmailReader interface {
	ReadEmail(reader *mail.Reader, uid uint32) (domain.EmailToSave, bool)
}

/*
NewEmailReader - creates new EmailReader entity
*/
func NewEmailReader(regexpConfig config.MailStructure) EmailReader {
	regExpMap := make(map[string]*regexp.Regexp)
	regExpMap[attachedFileName] = regexp.MustCompile(regexpConfig.FileNameRegExp)
	regExpMap[whoCalls] = regexp.MustCompile(regexpConfig.WhoCallsRegExp)
	regExpMap[inputNumber] = regexp.MustCompile(regexpConfig.InputNumberRegExp)
	regExpMap[participant] = regexp.MustCompile(regexpConfig.ParticipantRegExp)
	regExpMap[callLength] = regexp.MustCompile(regexpConfig.CallLengthRegExp)
	return &emailReader{regExpMap}
}

func (emailReader *emailReader) ReadEmail(reader *mail.Reader, uid uint32) (domain.EmailToSave, bool) {
	emailToSave := domain.EmailToSave{}
	emailData := domain.EmailData{Uid: uid}
	foundTextData := false
	foundAttachedFile := false
	for {
		p, err := reader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Error reading next part of email:", err)
		}
		switch h := p.Header.(type) {
		case mail.TextHeader:
			// This is the message's text (can be plain-text or HTML)
			b, _ := ioutil.ReadAll(p.Body)
			partText := string(b)
			if !foundTextData && emailReader.matchText(partText) {
				// TODO extract groups of matched text and save it in structure
				emailData.WhoCalls = extractField(emailReader.regExpMap[whoCalls], partText)
				emailData.Participant = extractField(emailReader.regExpMap[participant], partText)
				emailData.InputNumber = extractField(emailReader.regExpMap[inputNumber], partText)
				emailData.Duration = extractField(emailReader.regExpMap[callLength], partText)
				foundTextData = true
			}
		case mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			if !foundAttachedFile && emailReader.matchFileName(filename) {
				emailData.RecordFileName = filename
				emailData.Date = getDateByFileName(emailReader.regExpMap[attachedFileName], filename)
				var buffer bytes.Buffer
				_, err := buffer.ReadFrom(p.Body)
				if err != nil {
					log.Printf("Error buffering attached file:%s; error:%v\n", filename, err)
				} else {
					emailToSave.Buffer = &buffer
					foundAttachedFile = true
				}
			}
		}
	}
	emailToSave.EmailData = emailData
	return emailToSave, foundAttachedFile && foundTextData
}

func getDateByFileName(regExp *regexp.Regexp, filename string) time.Time {
	slice := regExp.FindAllStringSubmatch(filename, -1)[0]
	// prepare time in format"2006-01-02 15:04:05"
	if len(slice) == 7 {
		day := addLeadingZerroz(slice[1], 2)
		month := monthes[slice[2]]
		year := slice[3]
		hours := addLeadingZerroz(slice[4], 2)
		minutes := addLeadingZerroz(slice[5], 2)
		seconds := addLeadingZerroz(slice[6], 2)
		parsedTime, err := time.Parse(timeutil.SQLTimestamp, fmt.Sprintf(dateTimePattern, year, month, day, hours, minutes, seconds))
		if err != nil {
			log.Printf("Error happened during file day-time parsing: %v; So, returning current time\n", err)
			return time.Now()
		}
		return parsedTime
	}
	return time.Now()
}

func addLeadingZerroz(value string, expectedLength int) string {
	result := value
	if len(value) < expectedLength {
		for i := len(value); i < expectedLength; i++ {
			result = "0" + result
		}
	}
	return result
}

func extractField(regExp *regexp.Regexp, text string) string {
	return regExp.FindAllStringSubmatch(text, -1)[0][1]
}
func (emailReader *emailReader) matchText(text string) bool {
	ok := emailReader.regExpMap[whoCalls].MatchString(text)
	ok = ok && emailReader.regExpMap[inputNumber].MatchString(text)
	ok = ok && emailReader.regExpMap[participant].MatchString(text)
	return ok
}
func (emailReader *emailReader) matchFileName(text string) bool {
	return emailReader.regExpMap[attachedFileName].MatchString(text)
}
