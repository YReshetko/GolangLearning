package fetch

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mailclient/config"
	"mailclient/domain"
	"regexp"
	"time"

	"github.com/emersion/go-message/mail"
	"github.com/grokify/gotilla/time/timeutil"
)

const (
	ATTACHED_FILE_NAME = "attached-file-name"
	WHO_CALLS          = "who-calls"
	INPUT_NUMBER       = "input-number"
	PARTICIPANT        = "participant"
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

type EmailReader interface {
	ReadEmail(reader *mail.Reader, uid uint32) (domain.EmailToSave, bool)
}

func NewEmailReader(regexpConfig config.MailStructure) EmailReader {
	regExpMap := make(map[string]*regexp.Regexp)
	regExpMap[ATTACHED_FILE_NAME] = regexp.MustCompile(regexpConfig.FileNameRegExp)
	regExpMap[WHO_CALLS] = regexp.MustCompile(regexpConfig.WhoCallsRegExp)
	regExpMap[INPUT_NUMBER] = regexp.MustCompile(regexpConfig.InputNumberRegExp)
	regExpMap[PARTICIPANT] = regexp.MustCompile(regexpConfig.ParticipantRegExp)
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
			fmt.Println(err)
		}
		switch h := p.Header.(type) {
		case mail.TextHeader:
			// This is the message's text (can be plain-text or HTML)
			b, _ := ioutil.ReadAll(p.Body)
			partText := string(b)
			if !foundTextData && emailReader.matchText(partText) {
				// TODO extract groups of matched text and save it in structure
				emailData.WhoCalls = extractField(emailReader.regExpMap[WHO_CALLS], partText)
				emailData.Participant = extractField(emailReader.regExpMap[PARTICIPANT], partText)
				emailData.InputNumber = extractField(emailReader.regExpMap[INPUT_NUMBER], partText)
				foundTextData = true
			}
		case mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			if !foundAttachedFile && emailReader.matchFileName(filename) {
				emailData.RecordFileName = filename
				emailData.Date = getDateByFileName(emailReader.regExpMap[ATTACHED_FILE_NAME], filename)
				var buffer bytes.Buffer
				buffer.ReadFrom(p.Body)
				emailToSave.Buffer = &buffer
				foundAttachedFile = true
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
		time, err := time.Parse(timeutil.SQLTimestamp, fmt.Sprintf("%s-%s-%s %s:%s:%s", year, month, day, hours, minutes, seconds))
		if err != nil {
			fmt.Println("Error happened during file day-time parsing:", err)
		}
		return time
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
	ok := true
	ok1 := ok && emailReader.regExpMap[WHO_CALLS].MatchString(text)
	ok2 := ok && emailReader.regExpMap[INPUT_NUMBER].MatchString(text)
	ok3 := ok && emailReader.regExpMap[PARTICIPANT].MatchString(text)
	fmt.Printf("Email text: %v\n", text)
	fmt.Printf("Match value: %v, %v, %v\n", ok1, ok2, ok3)
	return ok1 && ok2 && ok3
}
func (emailReader *emailReader) matchFileName(text string) bool {
	return emailReader.regExpMap[ATTACHED_FILE_NAME].MatchString(text)
}
