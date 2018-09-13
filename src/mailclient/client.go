package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"mailclient/config"
	"mailclient/fetch"
	"os"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
)

func main() {
	var config config.Configuration
	exec(config.Load())
	log(config)
	client := fetch.NewImapClient()
	client.Init(config)
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
	//iterateByEmails(client)
	getMessagestByChannel(client)

}

func getMessagestByChannel(client fetch.ImapClient) {
	done := make(chan bool, 2)
	messagesChannel, err := client.GetMessageChannel("INBOX", done)
	if err != nil {
		fmt.Println("Loading mail channel error:", err)
		os.Exit(1)
	}

	count := 0
	for msg := range messagesChannel {
		count++
		fmt.Printf("Email sender %+v:\n", *msg.Envelope.Sender[0])
		fmt.Println("Mail ID:", msg.Envelope.MessageId)
		fmt.Println("Mail subject:", msg.Envelope.Subject)
		fmt.Println("Mail body:", msg.Body)
		if count > 101 {
			done <- true
			break
		}
	}
	fmt.Println("Emails number: ", count)
}

func iterateByEmails(client fetch.ImapClient) {
	iterator, err := client.MailIterator("INBOX")
	if err != nil {
		fmt.Println("Cant retrieve mail iterator:", err)
		os.Exit(1)
	}
	for ; ; iterator.HasNext() {
		msgChan, err := iterator.Next()
		if err != nil {
			fmt.Println("Cant retrieve an mail:", err)
			os.Exit(1)
		}
		msg := <-msgChan
		fmt.Println("----------------------------------")
		//fmt.Println("Mail ID:", msg.Envelope.MessageId)
		//fmt.Println("Mail subject:", msg.Envelope.Subject)
		fmt.Println("Mail body:", msg.Body)

		section := &imap.BodySectionName{}
		r := msg.GetBody(section)
		if r == nil {
			fmt.Println("Server didn't returned message body")
			continue
		}

		// Create a new mail reader
		mr, err := mail.CreateReader(r)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Process each message's part
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}

			switch h := p.Header.(type) {
			case mail.TextHeader:
				// This is the message's text (can be plain-text or HTML)
				b, _ := ioutil.ReadAll(p.Body)
				fmt.Println("Got text: %v", string(b))
			case mail.AttachmentHeader:
				// This is an attachment
				filename, _ := h.Filename()
				fmt.Println("Got attachment: %v", filename)
			}
		}
	}
}

func getMessagestFrom(client fetch.ImapClient) {
	messages, err := client.Select("INBOX", 1, 1)
	if err != nil {
		fmt.Println("Loading emails error:", err)
		os.Exit(1)
	}
	for msg := range messages {
		fmt.Println("Mail ID:", msg.Envelope.MessageId)
		fmt.Println("Mail subject:", msg.Envelope.Subject)
		fmt.Println("Mail body:", msg.Body)
	}
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
