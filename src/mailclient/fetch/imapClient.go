package fetch

import (
	"fmt"
	"log"
	"mailclient/config"
	"strconv"

	"github.com/emersion/go-imap"

	"github.com/emersion/go-imap/client"
)

type imapClient struct {
	config config.Configuration
	client *client.Client
}

func NewImapClient() ImapClient {
	return &imapClient{}
}

type ImapClient interface {
	Init(config config.Configuration)
	Connect() error
	Login() error
	Logout() error
	Mailboxes() (chan *imap.MailboxInfo, error)
	Select(box string, from, to uint32) (chan *imap.Message, error)
}

func (cli *imapClient) Init(config config.Configuration) {
	cli.config = config
}
func (cli *imapClient) Connect() error {
	server := cli.config.ImapHost + ":" + strconv.Itoa(cli.config.ImapPort)
	fmt.Printf("Server: %s", server)
	c, err := client.DialTLS(server, nil)
	cli.client = c
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("Connected")
	return nil
}

func (cli *imapClient) Login() error {
	if err := cli.client.Login(cli.config.ClientEmail, cli.config.ClientPassword); err != nil {
		return err
	}
	log.Println("Logged")
	return nil
}
func (cli *imapClient) Logout() error {
	if cli.client != nil {
		cli.client.Logout()
		log.Println("Logout")
	}
	return nil
}

func (cli *imapClient) Mailboxes() (chan *imap.MailboxInfo, error) {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- cli.client.List("", "*", mailboxes)
	}()
	if err := <-done; err != nil {
		log.Fatal(err)
		return nil, err
	}
	return mailboxes, nil
}

func (cli *imapClient) Select(box string, from, to uint32) (chan *imap.Message, error) {
	mbox, err := cli.client.Select(box, false)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if mbox.Messages < to {
		to = mbox.Messages
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- cli.client.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()
	if err := <-done; err != nil {
		log.Fatal(err)
		return nil, err
	}

	return messages, nil
}
