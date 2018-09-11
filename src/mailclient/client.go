package main

import (
	"fmt"
	"mailclient/config"
	"mailclient/fetch"
	"os"
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

	messages, err := client.Select("INBOX", 1, 2)
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
