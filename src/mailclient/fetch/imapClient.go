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

func NewImapClient(config config.Configuration) ImapClient {
	return &imapClient{config: config}
}

type ImapClient interface {
	Connect() error
	Login() error
	Logout() error
	Mailboxes() (chan *imap.MailboxInfo, error)
	GetMessageChannel(box string, done chan bool) (chan *imap.Message, error)
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

func startFetching(messagesOut chan *imap.Message, done chan bool, cli *imapClient, fetchManager FetchManager) {
	defer func() {
		fmt.Println("Closing output channel!!!!")
		close(messagesOut)
	}()
	bufferCompleted := make(chan error, 1)
	//section := &imap.BodySectionName{}
	//fetchItems := []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid}
	//fetchItems := []imap.FetchItem{imap.FetchEnvelope, section.FetchItem(), imap.FetchUid}
	//fetchItems := []imap.FetchItem{imap.FetchEnvelope}
	//fetchItems := []imap.FetchItem{section.FetchItem()}
	//var seqset *imap.SeqSet
	//messagesNumber++
	for fetchManager.HasNext() {
		//from, to := recalculateMessageRange(&messagesNumber, &bufferSize)
		chanSize := fetchManager.BufferSize() + 1
		//fmt.Printf("Range: %v-%v; Channel size: %v; messagesNumber: %v; bufferSize: %v\n", from, to, chanSize, messagesNumber, bufferSize)
		messages := make(chan *imap.Message, chanSize)
		select {
		case <-done:
			break
		default:
			//seqset = new(imap.SeqSet)
			//seqset.AddRange(from, to)
			bufferCompleted <- fetchManager.FetchFunction()(fetchManager.NextSequenceSet(), fetchManager.FetchItems(), messages)
		}
		needsContinue := true
		if err := <-bufferCompleted; err == nil {
			needsContinue = redirectMessages(messages, messagesOut, done)
		} else {
			log.Fatal(err)
			break
		}
		if !needsContinue {
			break
		}
	}

}

func redirectMessages(from, to chan *imap.Message, done chan bool) bool {
	needsContinue := true
	for msg := range from {
		select {
		case <-done:
			fmt.Println("DONE!!!!!!")
			needsContinue = false
			break
		default:
			to <- msg
		}
	}
	return needsContinue
}

func getNumMessagesForBox(box string, client *client.Client) (uint32, error) {
	mbox, err := client.Select(box, false)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return mbox.Messages, nil
}
