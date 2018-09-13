package fetch

import (
	"log"

	"github.com/emersion/go-imap"

	"github.com/emersion/go-imap/client"
)

type mailIterator struct {
	currentEmailIndex uint32
	client            *client.Client
}

type Iterator interface {
	HasNext() bool
	Next() (chan *imap.Message, error)
}

func (iter *mailIterator) HasNext() bool {
	if iter.currentEmailIndex > 0 {
		return true
	}
	return false
}

func (iter *mailIterator) Next() (chan *imap.Message, error) {
	seqset := new(imap.SeqSet)
	iter.currentEmailIndex--
	seqset.AddRange(iter.currentEmailIndex, iter.currentEmailIndex+1)

	messages := make(chan *imap.Message, 2)
	done := make(chan error, 1)
	go func() {
		section := &imap.BodySectionName{}
		//done <- iter.client.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		done <- iter.client.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
	}()
	if err := <-done; err != nil {
		log.Fatal(err)
		return nil, err
	}

	return messages, nil
}
