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
	config config.HostConfig
	client *client.Client
}

func NewImapClient(config config.HostConfig) ImapClient {
	return &imapClient{config: config}
}

type ImapClient interface {
	Connect() error
	Login() error
	Logout() error
	Mailboxes() (chan *imap.MailboxInfo, error)
	GetMessageChannel(box string, done chan bool) (chan *imap.Message, error)
	GetMessageBodyChannel(box string, uids []uint32) (chan *imap.Message, error)
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

func (cli *imapClient) GetMessageChannel(box string, done chan bool) (chan *imap.Message, error) {
	messagesNumber, err := getNumMessagesForBox(box, cli.client)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	bufferSize := uint32(100)
	messagesOut := make(chan *imap.Message, bufferSize/2)
	fmt.Printf("Initial messages number: %v\n", messagesNumber)
	go startFetching(messagesOut, done, cli, NewEnvelopFetchManager(cli.client.Fetch, messagesNumber, bufferSize))
	return messagesOut, nil
}

func (cli *imapClient) GetMessageBodyChannel(box string, uids []uint32) (chan *imap.Message, error) {
	bufferSize := uint32(2)
	messagesOut := make(chan *imap.Message, bufferSize/2)
	fmt.Printf("Initial uids number: %v\n", len(uids))
	go startFetching(messagesOut, nil, cli, NewBodyFetchManager(cli.client.UidFetch, uids, bufferSize))
	return messagesOut, nil
}

func startFetching(messagesOut chan *imap.Message, done chan bool, cli *imapClient, fetchManager FetchManager) {
	defer func() {
		fmt.Println("Closing output channel!!!!")
		close(messagesOut)
	}()
	bufferCompleted := make(chan error, 1)
	for fetchManager.HasNext() {
		chanSize := fetchManager.BufferSize() + 1
		messages := make(chan *imap.Message, chanSize)
		fmt.Println("Messages chan size", chanSize)
		select {
		case <-done:
			fmt.Println("Can read from done!!!!")
			return
		default:
			fmt.Println("Start fetching emails")
			seqset := fetchManager.NextSequenceSet()
			fmt.Printf("Surrent subseq:%v\n", seqset)
			bufferCompleted <- fetchManager.FetchFunction()(seqset, fetchManager.FetchItems(), messages)
			fmt.Println("Complete fetching emails")
		}
		if err := <-bufferCompleted; err == nil {
			fmt.Println("Start redirecting")
			redirectMessages(messages, messagesOut)
			fmt.Println("Complete redirecting")
		} else {
			fmt.Println("Got an error")
			log.Fatal(err)
			break
		}
	}
}

func redirectMessages(from, to chan *imap.Message) {
	for msg := range from {
		to <- msg
	}
}

func getNumMessagesForBox(box string, client *client.Client) (uint32, error) {
	mbox, err := client.Select(box, false)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return mbox.Messages, nil
}
