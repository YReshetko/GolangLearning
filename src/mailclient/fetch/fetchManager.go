package fetch

import (
	"fmt"

	imap "github.com/emersion/go-imap"
)

/*
Fetch Items example
section := &imap.BodySectionName{}
fetchItems := []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid}
fetchItems := []imap.FetchItem{imap.FetchEnvelope, section.FetchItem(), imap.FetchUid}
fetchItems := []imap.FetchItem{imap.FetchEnvelope}
fetchItems := []imap.FetchItem{section.FetchItem()}
*/
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
			buffersize,
		},
		messagesNumber,
	}
}

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
	fmt.Println("Next sequence set: ", slice)
	fmt.Printf("Seqset: %+v\n", seqset)
	return seqset
}
func (manager *bodyFetchManager) HasNext() bool {
	currentIndex := manager.currentIndex
	newIndex := currentIndex + int(manager.buffersize)
	if newIndex >= len(manager.uids) {
		manager.buffersize = uint32(len(manager.uids)) - uint32(currentIndex)
		fmt.Println("Recalculated buffer size:", manager.buffersize)
	}
	return manager.currentIndex < len(manager.uids)
}
