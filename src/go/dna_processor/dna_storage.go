package dna_processor

import (
	"fmt"
	"strings"
)

type DnaStorage interface {
	GetChar() byte
	UndoGet()
	IsEmpty() bool
	Index(string) int
	Skip(n int)
	String() string
	PrependPrefix(s string)
	Len() int
	SaveOffset()
	RestoreOffset()
}

type SimpleDnaStorage struct {
	offset      int
	savedOffset int
	data        string
}

func NewSimpleDnaStorage(data string) *SimpleDnaStorage {
	return &SimpleDnaStorage{data: data}
}

// GetChar retrieves the next character from the storage
func (storage *SimpleDnaStorage) GetChar() byte {
	result := storage.data[storage.offset]
	storage.offset += 1
	return result
}

// UndoGet reverts the effect of the last GetChar call
func (storage *SimpleDnaStorage) UndoGet() {
	storage.offset -= 1
}

// IsEmpty checks if the storage is empty
func (storage *SimpleDnaStorage) IsEmpty() bool {
	return storage.offset >= len(storage.data)
}

func (storage *SimpleDnaStorage) Index(s string) int {
	idx := strings.Index(storage.data[storage.offset:], s)
	if idx == -1 {
		return -1 // not found
	}
	return idx // adjust by adding the offset
}

func (storage *SimpleDnaStorage) Skip(n int) {
	storage.offset += n
}

func (storage *SimpleDnaStorage) String() string {
	return storage.data[storage.offset:]
}

func (storage *SimpleDnaStorage) PrependPrefix(s string) {
	storage.data = s + storage.data[storage.offset:]
	storage.offset = 0
}

func (storage *SimpleDnaStorage) Len() int {
	rest := len(storage.data) - storage.offset
	if rest > 0 {
		return rest
	}
	return 0
}

func (storage *SimpleDnaStorage) SaveOffset() {
	fmt.Print("saved offset: ", storage.savedOffset)
	storage.savedOffset = storage.offset
}

func (storage *SimpleDnaStorage) RestoreOffset() {
	fmt.Print("restored offset: replace", storage.offset, " with ", storage.savedOffset)
	storage.offset = storage.savedOffset
}
