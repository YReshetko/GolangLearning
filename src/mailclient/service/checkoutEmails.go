package service

import (
	"fmt"
	"log"
	"mailclient/config"
	"mailclient/fetch"
	"mailclient/save"
	"sort"
	"time"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
)

const inboxName = "INBOX"

type checkoutEmails struct {
	imapClient     fetch.ImapClient
	emailReader    fetch.EmailReader
	emailSaver     save.EmailSaver
	dao            save.EmailDao
	inProgerss     bool
	collectionName string
	expectedSender string
}

/*
EmailService - runs full email fetching saving into DB and local storage
*/
type EmailService interface {
	Process() error
	PrintMailboxes()
}

/*
NewEmailFetcher - create new email fetcher by checkoutEmails impl
*/
func NewEmailFetcher(config config.Configuration, dao save.EmailDao) EmailService {
	return &checkoutEmails{
		imapClient:     fetch.NewImapClient(config.HostConfiguration),
		emailReader:    fetch.NewEmailReader(config.EmailStructure),
		dao:            dao,
		emailSaver:     save.NewEmailSaver(config.StorageConfiguration.LocalStorageBasePath),
		inProgerss:     false,
		collectionName: config.StorageConfiguration.CollectionName,
		expectedSender: config.EmailStructure.ExpectedSender,
	}
}

type checkoutError struct {
	message string
}

func (err checkoutError) Error() string {
	return fmt.Sprintf("Error in fetch service %s", err.message)
}

func (saver *checkoutEmails) Process() error {
	if saver.inProgerss {
		return checkoutError{"The service is in progress already!"}
	}
	defer saver.postProcess()
	if err := saver.preProcess(); err == nil {
		uids := saver.findUnprocessedEmailUids()
		log.Println("Found unprocessed UIDs:", uids)
		if len(uids) > 0 {
			time.Sleep(2 * time.Second)
			saver.processUids(uids)
		}
	} else {
		return err
	}
	return nil
}

func (saver *checkoutEmails) preProcess() error {
	saver.inProgerss = true
	if err := saver.imapClient.Connect(); err != nil {
		return err
	}
	if err := saver.imapClient.Login(); err != nil {
		return err
	}
	return nil
}
func (saver *checkoutEmails) postProcess() error {
	saver.inProgerss = false
	saver.dao = nil
	if err := saver.imapClient.Logout(); err != nil {
		return err
	}
	return nil
}

func (saver *checkoutEmails) findUnprocessedEmailUids() []uint32 {
	done := make(chan bool, 2)
	messagesChannel, err := saver.imapClient.GetMessageEnvelopChannel(inboxName, done)
	if err != nil {
		done <- false
		log.Println("Error during retrieving channel for UID fetching:", err)
		return make([]uint32, 0)
	}

	count := 0
	uidsToProcess := make([]uint32, 0, 100)
	isComplated := false
	for msg := range messagesChannel {
		count++
		if saver.needsProcessing(msg) {
			if saver.uidProcessedBefore(msg.Uid) {
				if !isComplated {
					isComplated = true
					done <- true
				}
			} else {
				uidsToProcess = append(uidsToProcess, msg.Uid)
			}
		}
	}
	if len(uidsToProcess) > 0 {
		sort.Slice(uidsToProcess, func(i, j int) bool {
			return uidsToProcess[i] < uidsToProcess[j]
		})
	}
	return uidsToProcess
}

func (saver *checkoutEmails) needsProcessing(msg *imap.Message) bool {
	var from string
	if msg.Envelope != nil {
		from = msg.Envelope.From[0].MailboxName
	}
	if from == saver.expectedSender {
		log.Printf("Found possible email for fetching: %+v\n", msg)
		return true
	}
	return false
}
func (saver *checkoutEmails) uidProcessedBefore(uid uint32) bool {
	if data := saver.dao.FindByUid(uid); data != nil {
		log.Println("Found UID processed before:", uid)
		return true
	}
	return false
}

func (saver *checkoutEmails) processUids(uids []uint32) {
	messagesChannel, err := saver.imapClient.GetMessageBodyChannel(inboxName, uids)
	if err != nil {
		log.Println("Error during retrieving channel for fetching email content:", err)
		return
	}
	metError := false
	count := uint32(0)
	for msg := range messagesChannel {
		if !metError {
			if err := saver.processEmail(msg); err != nil {
				log.Printf("Met error during email processing, so skiping further processing")
				metError = true
			} else {
				count++
			}
			if count%100 == 0 {
				log.Println("Saved emails count: ", count)
			}
		}
	}
	log.Printf("Finally processed %v emails at: %v\n", count, time.Now())
}

func (saver *checkoutEmails) processEmail(msg *imap.Message) error {
	section := &imap.BodySectionName{}
	r := msg.GetBody(section)
	if r == nil {
		log.Printf("Server didn't returned message body: %+v\n", msg)
		return nil
	}
	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		log.Printf("Cant create email reader for: %+v, due to %v\n", msg, err)
		return err
	}
	mailToSave, ok := saver.emailReader.ReadEmail(mr, msg.Uid)
	if ok {
		log.Println("Saving email UID:", msg.Uid)
		err := saver.emailSaver.Save(&mailToSave)
		if err != nil {
			log.Println("Error saving attached file:", err)
			return err
		}
		err = saver.dao.Save(mailToSave.EmailData)
		if err != nil {
			log.Println("Error saving email info:", err)
			return err
		}
	}
	// Process each message's part
	return nil
}

func (saver *checkoutEmails) PrintMailboxes() {
	saver.preProcess()
	defer saver.postProcess()
	boxinfos, err := saver.imapClient.Mailboxes()
	if err != nil {
		log.Println("Loading mail boxes error:", err)
		return
	}
	for info := range boxinfos {
		log.Println(info)
	}
}
