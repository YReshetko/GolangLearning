package fetch

import (
	imap "github.com/emersion/go-imap"
)

type fetchFunc func(*imap.SeqSet, []imap.FetchItem, chan *imap.Message) error

type fetchManager struct {
	fetch      fetchFunc
	items      []imap.FetchItem
	buffersize uint32
}

type envelopFetchManager struct {
	fetchManager
	messagesNumber uint32
}

type bodyFetchManager struct {
	fetchManager
	uids         []uint32
	currentIndex int
}

/*
FetchManager - controles methods and parameters of email fetching
	FetchFunction - fetch function from go-imap client (Fetch(), FetchUid())
	FetchItems - Array of fetch items Uid, Envelop, Body, etc
	BufferSize - size of the buffer for fetching several emails
	NextSequenceSet - imap.SeqSet of next emails portion
	HasNext - true if we can fetch more emails
*/
type FetchManager interface {
	FetchFunction() fetchFunc
	FetchItems() []imap.FetchItem
	BufferSize() uint32
	NextSequenceSet() *imap.SeqSet
	HasNext() bool
}

/*
NewEnvelopFetchManager - creates fetch manager for email envelop fetching
*/
func NewEnvelopFetchManager(fetch fetchFunc, messagesNumber uint32, buffersize uint32) FetchManager {
	return &envelopFetchManager{
		fetchManager{
			fetch, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid},
			buffersize,
		},
		messagesNumber,
	}
}

/*
NewBodyFetchManager - creates fetch manager for email body fetching
*/
func NewBodyFetchManager(fetch fetchFunc, uids []uint32, buffersize uint32) FetchManager {
	section := &imap.BodySectionName{}
	return &bodyFetchManager{
		fetchManager{
			fetch, []imap.FetchItem{section.FetchItem(), imap.FetchUid},
			buffersize,
		},
		uids,
		0,
	}
}

func (manager *fetchManager) FetchFunction() fetchFunc {
	return manager.fetch
}
func (manager *fetchManager) FetchItems() []imap.FetchItem {
	return manager.items
}
func (manager *fetchManager) BufferSize() uint32 {
	return manager.buffersize
}

func (manager *envelopFetchManager) NextSequenceSet() *imap.SeqSet {
	from, to := manager.recalculateMessageRange()
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)
	return seqset
}
func (manager *envelopFetchManager) HasNext() bool {
	return manager.messagesNumber > 1
}
func (manager *envelopFetchManager) recalculateMessageRange() (uint32, uint32) {
	if manager.messagesNumber-manager.buffersize > manager.messagesNumber {
		manager.buffersize = manager.messagesNumber - 2
		manager.messagesNumber = 1
	} else {
		manager.messagesNumber = manager.messagesNumber - manager.buffersize - 1
	}
	return manager.messagesNumber, manager.messagesNumber + manager.buffersize
}

func (manager *bodyFetchManager) NextSequenceSet() *imap.SeqSet {
	currentIndex := manager.currentIndex
	newIndex := currentIndex + int(manager.buffersize)
	slice := make([]uint32, int(manager.buffersize))
	if newIndex < len(manager.uids) {
		slice = manager.uids[currentIndex:newIndex]
		manager.currentIndex = newIndex
		manager.buffersize = uint32(newIndex) - uint32(currentIndex)
	} else {
		slice = manager.uids[currentIndex:]
		manager.currentIndex = len(manager.uids)
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(slice...)
	return seqset
}
func (manager *bodyFetchManager) HasNext() bool {
	currentIndex := manager.currentIndex
	newIndex := currentIndex + int(manager.buffersize)
	if newIndex >= len(manager.uids) {
		manager.buffersize = uint32(len(manager.uids)) - uint32(currentIndex)
	}
	return manager.currentIndex < len(manager.uids)
}
