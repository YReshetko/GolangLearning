package fetch

import imap "github.com/emersion/go-imap"

type fetchFunc func(*imap.SeqSet, []imap.FetchItem, chan *imap.Message) error

type fetchManager struct {
	fetch fetchFunc
	items []imap.FetchItem
}

type envelopFetchManager struct {
	fetchManager
	messagesNumber uint32
	buffersize     uint32
}

type bodyFetchManager struct {
	fetchManager
}

type FetchManager interface {
	FetchFunction() fetchFunc
	FetchItems() []imap.FetchItem
	BufferSize() uint32
	NextSequenceSet() *imap.SeqSet
	HasNext() bool
}

func NewEnvelopFetchManager(fetch fetchFunc, messagesNumber uint32, buffersize uint32) FetchManager {
	return &envelopFetchManager{
		fetchManager{
			fetch, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid},
		},
		messagesNumber,
		buffersize,
	}
}

func (manager *envelopFetchManager) FetchFunction() fetchFunc {
	return manager.fetch
}
func (manager *envelopFetchManager) FetchItems() []imap.FetchItem {
	return manager.items
}
func (manager *envelopFetchManager) BufferSize() uint32 {
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
