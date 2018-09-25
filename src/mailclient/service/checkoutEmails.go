package service

import (
	"fmt"
	"mailclient/config"
	"mailclient/fetch"
	"mailclient/save"
	"os"
	"time"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
)

type checkoutEmails struct {
	imapClient     fetch.ImapClient
	emailReader    fetch.EmailReader
	dbAccess       save.DBAccess
	emailSaver     save.EmailSaver
	dao            save.EmailDao
	inProgerss     bool
	collectionName string
	expectedSender string
}

type EmailService interface {
	Process() error
}

func NewEmailFetcher(config config.Configuration) EmailService {
	return &checkoutEmails{
		imapClient:     fetch.NewImapClient(config.HostConfiguration),
		emailReader:    fetch.NewEmailReader(config.EmailStructure),
		dbAccess:       save.NewDBAccess(config.StorageConfiguration.DbHost, config.StorageConfiguration.DbPort, config.StorageConfiguration.DbName),
		emailSaver:     save.EmailSaverInstance(config.StorageConfiguration.LocalStorageBasePath),
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
	} else {
		if err := saver.preProcess(); err == nil {
			defer saver.postProcess()
			uids := saver.findUnprocessedEmailUids()
			if len(uids) > 0 {
				time.Sleep(2 * time.Second)
				saver.processUids(uids)
			}
		} else {
			return err
		}
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
	if ok := saver.dbAccess.StartSession(); !ok {
		return checkoutError{"Can't start db session!"}
	}
	saver.dao = save.NewDao(saver.dbAccess.GetCollection(saver.collectionName))
	return nil
}
func (saver *checkoutEmails) postProcess() error {
	saver.inProgerss = false
	saver.dao = nil
	if ok := saver.dbAccess.CloseSession(); !ok {
		return checkoutError{"Can't close db session!"}
	}
	if err := saver.imapClient.Logout(); err != nil {
		return err
	}
	return nil
}

func (saver *checkoutEmails) findUnprocessedEmailUids() []uint32 {
	done := make(chan bool, 2)
	messagesChannel, err := saver.imapClient.GetMessageChannel("INBOX", done)
	if err != nil {
		fmt.Println("Loading mail channel error:", err)
		os.Exit(1)
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
		// TODO Remove as redundand
		if count > 100 && !isComplated {
			done <- true
			break
		}
	}
	return uidsToProcess
}

func (saver *checkoutEmails) needsProcessing(msg *imap.Message) bool {
	var from string
	if msg.Envelope != nil {
		from = msg.Envelope.From[0].MailboxName
	}
	if from == saver.expectedSender {
		fmt.Println("Found expected email:")
		fmt.Printf("Email sender: %+v\n", msg.Envelope)
		fmt.Printf("Email uid: %+v\n", msg.Uid)
		fmt.Printf("From: %+v\n", msg.Envelope.From[0].MailboxName)
		fmt.Println("Mail subject:", msg.Envelope.Subject)
		return true
	}
	return false
}
func (saver *checkoutEmails) uidProcessedBefore(uid uint32) bool {
	if data := saver.dao.FindByUid(uid); data != nil {
		return true
	}
	return false
}

func (saver *checkoutEmails) processUids(uids []uint32) {
	messagesChannel, err := saver.imapClient.GetMessageBodyChannel("INBOX", uids)
	if err != nil {
		fmt.Println("Loading mail channel error:", err)
		os.Exit(1)
	}
	for msg := range messagesChannel {
		saver.processEmail(msg)
	}
}

func (saver *checkoutEmails) processEmail(msg *imap.Message) error {
	section := &imap.BodySectionName{}
	r := msg.GetBody(section)
	if r == nil {
		fmt.Println("Server didn't returned message body")
		return nil
	}

	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Create email reader")
	fmt.Println("Start reading email")
	mailToSave, ok := saver.emailReader.ReadEmail(mr, msg.Uid)
	if ok {
		fmt.Println("found email to save")
		fmt.Printf("The structure of email to save: %+v\n", mailToSave)
		saver.emailSaver.Save(&mailToSave)
		saver.dao.Save(mailToSave.EmailData)
	}
	// Process each message's part
	return nil
}
