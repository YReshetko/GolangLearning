package fetch

import (
	"log"
	"mailclient/config"

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
	Login() error
	Logout() error
}

func (cli *imapClient) Init(config config.Configuration) {
	cli.config = config
}
func (cli *imapClient) Login() error {
	c, err := client.DialTLS(cli.config.ImapHost+":"+string(cli.config.ImapPort), nil)
	cli.client = c
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("Connected")
	return nil
}
func (cli *imapClient) Logout() error {
	if cli.client != nil {
		cli.client.Logout()
	}
	return nil
}
