package fetch

import (
	"fmt"
	"log"
	"mailclient/config"
	"strconv"

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
func (cli *imapClient) Logout() error {
	if cli.client != nil {
		cli.client.Logout()
	}
	return nil
}
