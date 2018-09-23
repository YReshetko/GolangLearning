package main

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

func main() {
	var config config.Configuration
	exec(config.Load())
	log(config)
	client := fetch.NewImapClient(config.HostConfiguration)
	exec(client.Connect())
	exec(client.Login())
	defer client.Logout()

	boxinfos, err := client.Mailboxes()
	if err != nil {
		fmt.Println("Loading mail boxes error:", err)
		os.Exit(1)
	}
	for info := range boxinfos {
		fmt.Println(info)
	}
	getMessagestByChannel(client, config)

}

func getMessagestByChannel(client fetch.ImapClient, config config.Configuration) {
	done := make(chan bool, 2)
	messagesChannel, err := client.GetMessageChannel("INBOX", done)
	if err != nil {
		fmt.Println("Loading mail channel error:", err)
		os.Exit(1)
	}

	count := 0
	uidsToProcess := make([]uint32, 0, 100)
	for msg := range messagesChannel {
		count++
		if needsProcessing(msg, config.EmailStructure.ExpectedSender) {
			if uidProcessedBefore(msg) {
				done <- true
			} else {
				uidsToProcess = append(uidsToProcess, msg.Uid)
			}

		}
		// TODO Remove as redundand
		if count > 100 {
			done <- true
			break
		}
	}
	time.Sleep(2 * time.Second)
	fmt.Printf("Found uids: %v\n", uidsToProcess)
	if len(uidsToProcess) > 0 {
		messagesChannel, err := client.GetMessageBodyChannel("INBOX", uidsToProcess)
		if err != nil {
			fmt.Println("Loading mail channel error:", err)
			os.Exit(1)
		}
		for msg := range messagesChannel {
			processEmail(msg, config.EmailStructure)
		}
	} else {
		fmt.Println("No emails to fetch")
	}

	fmt.Println("Emails number: ", count)
}

func needsProcessing(msg *imap.Message, expectedSender string) bool {
	var from string
	if msg.Envelope != nil {
		from = msg.Envelope.From[0].MailboxName
	}
	if from == expectedSender {
		fmt.Println("Found expected email:")
		fmt.Printf("Email sender: %+v\n", msg.Envelope)
		fmt.Printf("Email uid: %+v\n", msg.Uid)
		fmt.Printf("From: %+v\n", msg.Envelope.From[0].MailboxName)
		fmt.Println("Mail subject:", msg.Envelope.Subject)
		return true
	}
	return false
}
func uidProcessedBefore(msg *imap.Message) bool {
	// TODO implement checker via DB
	return false
}

func processEmail(msg *imap.Message, config config.MailStructure) error {
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
	bodyReader := fetch.NewEmailReader(config)
	fmt.Println("Start reading email")
	mailToSave, ok := bodyReader.ReadEmail(mr, msg.Uid)
	if ok {
		fmt.Println("found email to save")
		fmt.Printf("The structure of email to save: %+v\n", mailToSave)
		save.Save(&mailToSave)
	}
	// Process each message's part
	return nil
}

func exec(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func log(object interface{}) {
	fmt.Printf("Logging object: %+v\n", object)
}
