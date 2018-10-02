package fetch

import (
	"fmt"
	"log"
	"mailclient/config"
	"strconv"
	"time"

	"github.com/emersion/go-imap"

	"github.com/emersion/go-imap/client"
)

const (
	defaultEmailEnvelopBufferSize = uint32(100)
	defaultEmailBodyBufferSize    = uint32(10)
)

type imapClient struct {
	config config.HostConfig
	client *client.Client
}
type imapClientError struct {
	msg string
}

func (err imapClientError) Error() string {
	return fmt.Sprintf("IMAP Error: %s", err.msg)
}

type ErrorEmailFetching struct {
	message string
}

func (err ErrorEmailFetching) Error() string {
	return fmt.Sprintf("Fetch email error: %s", err.message)
}

/*
ImapClient - access to IMAP server
*/
type ImapClient interface {
	Connect() error
	Login() error
	Logout() error
	Mailboxes() (chan *imap.MailboxInfo, error)
	GetMessageEnvelopChannel(box string, done chan bool) (chan *imap.Message, error)
	GetMessageBodyChannel(box string, uids []uint32) (chan *imap.Message, error)
}

/*
NewImapClient - creates new ImapClient
*/
func NewImapClient(config config.HostConfig) ImapClient {
	return &imapClient{config: config}
}

func (cli *imapClient) Connect() error {
	server := cli.config.ImapHost + ":" + strconv.Itoa(cli.config.ImapPort)
	log.Printf("Connectiong to IMAP server: %s\n", server)
	c, err := client.DialTLS(server, nil)
	cli.client = c
	if err != nil {
		log.Println("Error during connecting to IMAP server", err)
		return err
	}
	log.Println("Connected to IMAP server")
	return nil
}

func (cli *imapClient) Login() error {
	if err := cli.client.Login(cli.config.ClientEmail, cli.config.ClientPassword); err != nil {
		return err
	}
	log.Println("Logged on IMAP server")
	return nil
}
func (cli *imapClient) Logout() error {
	if cli.client != nil {
		cli.client.Logout()
		log.Println("Logout on IMAP server")
		cli.client.Close()
		log.Println("Closed IMAP session")
	} else {
		log.Println("Logouting when the client does not exist")
	}
	return nil
}

func (cli *imapClient) Mailboxes() (chan *imap.MailboxInfo, error) {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		if cli.client != nil {
			done <- cli.client.List("", "*", mailboxes)
		} else {
			done <- imapClientError{"Problem with connecting to mailbox"}
		}

	}()
	if err := <-done; err != nil {
		log.Println("Error during fetching mailboxes:", err)
		return nil, err
	}
	return mailboxes, nil
}

func (cli *imapClient) GetMessageEnvelopChannel(box string, done chan bool) (chan *imap.Message, error) {
	messagesNumber, err := getNumMessagesForBox(box, cli.client)
	if err != nil {
		log.Println("Error retrieving number of messages from IMAP server:", err)
		return nil, err
	}
	bufferSize := defaultEmailEnvelopBufferSize
	messagesOut := make(chan *imap.Message, bufferSize/2)
	log.Printf("Initial messages number: %v\n", messagesNumber)
	go startFetching(messagesOut, done, cli, NewEnvelopFetchManager(cli.client.Fetch, messagesNumber, bufferSize))
	return messagesOut, nil
}

func (cli *imapClient) GetMessageBodyChannel(box string, uids []uint32) (chan *imap.Message, error) {
	bufferSize := defaultEmailBodyBufferSize
	messagesOut := make(chan *imap.Message, bufferSize/2)
	log.Printf("Initial uids number: %v\n", len(uids))
	go startFetching(messagesOut, nil, cli, NewBodyFetchManager(cli.client.UidFetch, uids, bufferSize))
	return messagesOut, nil
}

func startFetching(messagesOut chan *imap.Message, done chan bool, cli *imapClient, fetchManager FetchManager) {
	defer func() {
		log.Println("Closing major output channel")
		close(messagesOut)
	}()
	var fetchError error
	for fetchManager.HasNext() {
		chanSize := fetchManager.BufferSize() + 1
		messages := make(chan *imap.Message, chanSize)
		log.Println("Buffer size to fetch new portion:", chanSize)
		select {
		case <-done:
			return
		default:
			seqset := fetchManager.NextSequenceSet()
			log.Printf("Current subseq:%v, start fetching new portion\n", seqset)
			fetchError = nonBlockingFetch(fetchManager.FetchFunction(), seqset, fetchManager.FetchItems(), messages)
			log.Println("Complete fetching emails")
		}
		if fetchError == nil {
			redirectMessages(messages, messagesOut)
		} else {
			log.Println("Error during fetching emails from IMAP server:", fetchError)
			return
		}
	}
}

func nonBlockingFetch(fetchF fetchFunc, seqset *imap.SeqSet, items []imap.FetchItem, messages chan *imap.Message) error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- fetchF(seqset, items, messages)
	}()
	defaultTimeout := 20
	timeouts := defaultTimeout
	for {
		select {
		case err := <-errChan:
			return err
		default:
			timeouts--
			if timeouts == 0 {
				log.Printf("Long waiting time: %v sec. for fetching emails set: %v\n", defaultTimeout, seqset)
				return ErrorEmailFetching{"Stop fetching due to timeout"}
			}
			time.Sleep(time.Second)

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
		return 0, err
	}
	return mbox.Messages, nil
}
