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
	MailIterator(box string) (Iterator, error)
	Select(box string, from, to uint32) (chan *imap.Message, error)
	GetMessageChannel(box string, done chan bool) (chan *imap.Message, error)
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
func (cli *imapClient) MailIterator(box string) (Iterator, error) {
	messagesNumber, err := getNumMessagesForBox(box, cli.client)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &mailIterator{messagesNumber, cli.client}, nil
}
func (cli *imapClient) Select(box string, from, to uint32) (chan *imap.Message, error) {
	messagesNumber, err := getNumMessagesForBox(box, cli.client)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if messagesNumber < to {
		to = messagesNumber
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 1)
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

func (cli *imapClient) GetMessageChannel(box string, done chan bool) (chan *imap.Message, error) {
	messagesNumber, err := getNumMessagesForBox(box, cli.client)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	bufferSize := uint32(100)
	messagesOut := make(chan *imap.Message, bufferSize/2)

	bufferCompleted := make(chan error, 1)
	fmt.Printf("Initial messages number: %v\n", messagesNumber)
	go func() {
		defer func() {
			fmt.Println("Closing output channel!!!!")
			close(messagesOut)
		}()
		messagesNumber++
		for {
			if messagesNumber-bufferSize > messagesNumber {
				bufferSize = messagesNumber - 2
				messagesNumber = 1
			} else {
				messagesNumber = messagesNumber - bufferSize - 1
			}
			from := messagesNumber
			to := messagesNumber + bufferSize
			chanSize := bufferSize + 1
			fmt.Printf("Range: %v-%v; Channel size: %v\n", from, to, chanSize)

			messages := make(chan *imap.Message, chanSize)
			select {
			case <-done:
				break
			default:
				seqset := new(imap.SeqSet)
				seqset.AddRange(from, to)
				bufferCompleted <- cli.client.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
			}
			if err := <-bufferCompleted; err == nil {
				for msg := range messages {
					select {
					case <-done:
						fmt.Println("DONE!!!!!!")
						return
					default:
						messagesOut <- msg
					}

				}
			} else {
				log.Fatal(err)
				break
			}
			if messagesNumber == 1 {
				break
			}
		}

	}()
	return messagesOut, nil
}

func getNumMessagesForBox(box string, client *client.Client) (uint32, error) {
	mbox, err := client.Select(box, false)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return mbox.Messages, nil
}
