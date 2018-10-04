package service

import (
	"fmt"
	"mailclient/config"
	"mailclient/fetch"
	"mailclient/logger"
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
	PrintMailboxes() error
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
		logger.Debug("Found unprocessed UIDs:", uids)
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
		logger.Error("Error during retrieving channel for UID fetching:", err)
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
		logger.Debug("Current sender: ", from)
	}
	if from == saver.expectedSender {
		logger.Debug("Found possible email for fetching: %+v\n", msg)
		return true
	}
	return false
}
func (saver *checkoutEmails) uidProcessedBefore(uid uint32) bool {
	data, err := saver.dao.FindByUid(uid)
	if err != nil {
		logger.Error("Error during retrieving UID %v, error: %v\n", uid, err)
		// Returning true to stop general processing
		return true
	}
	if data != nil {
		logger.Debug("Found UID processed before:", uid)
		return true
	}
	return false
}

func (saver *checkoutEmails) processUids(uids []uint32) {
	messagesChannel, err := saver.imapClient.GetMessageBodyChannel(inboxName, uids)
	if err != nil {
		logger.Error("Error during retrieving channel for fetching email content:", err)
		return
	}
	metError := false
	count := uint32(0)
	for msg := range messagesChannel {
		if !metError {
			if err := saver.processEmail(msg); err != nil {
				logger.Error("Met error during email processing, so skiping further processing")
				metError = true
			} else {
				count++
			}
			if count%100 == 0 {
				logger.Debug("Saved emails count: ", count)
			}
		}
	}
	logger.Debug("Finally processed %v emails at: %v\n", count, time.Now())
}

func (saver *checkoutEmails) processEmail(msg *imap.Message) error {
	section := &imap.BodySectionName{}
	r := msg.GetBody(section)
	if r == nil {
		logger.Warning("Server didn't returned message body: %+v\n", msg)
		return nil
	}
	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		logger.Error("Cant create email reader for: %+v, due to %v\n", msg, err)
		return err
	}
	mailToSave, ok := saver.emailReader.ReadEmail(mr, msg.Uid)
	if ok {
		logger.Debug("Saving email UID:", msg.Uid)
		err := saver.emailSaver.Save(&mailToSave)
		if err != nil {
			logger.Error("Error saving attached file:", err)
			return err
		}
		err = saver.dao.Save(mailToSave.EmailData)
		if err != nil {
			logger.Error("Error saving email info:", err)
			return err
		}
	}
	// Process each message's part
	return nil
}

func (saver *checkoutEmails) PrintMailboxes() error {
	saver.preProcess()
	defer saver.postProcess()
	boxinfos, err := saver.imapClient.Mailboxes()
	if err != nil {
		logger.Debug("Loading mail boxes error:", err)
		return err
	}
	for info := range boxinfos {
		logger.Debug("Box info:%v\n", info)
	}
	return nil
}
