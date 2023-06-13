package dna_processor

import (
	"errors"
	"github.com/deadpixi/rope"
	"strings"
)

type DnaStorageRope struct {
	offset      int
	savedOffset int
	data        rope.Rope
}

func NewDnaStorageRope(data string) *DnaStorageRope {
	return &DnaStorageRope{data: rope.NewString(data)}
}

// GetChar retrieves the next character from the storage
func (storage *DnaStorageRope) GetChar() byte {
	result := storage.data.Slice(storage.offset, storage.offset+1)[0]
	storage.offset += 1
	return result
}

// UndoGet reverts the effect of the last GetChar call
func (storage *DnaStorageRope) UndoGet() {
	storage.offset -= 1
}

// IsEmpty checks if the storage is empty
func (storage *DnaStorageRope) IsEmpty() bool {
	return storage.offset >= storage.data.Length()
}

func (storage *DnaStorageRope) Index(s string) int {
	_, rightSlice := storage.data.Split(storage.offset)
	idx := strings.Index(rightSlice.String(), s)
	if idx == -1 {
		return -1 // not found
	}
	return idx // adjust by adding the offset
}

func (storage *DnaStorageRope) Skip(n int) {
	storage.offset += n
}

func (storage *DnaStorageRope) String() string {
	_, rightSlice := storage.data.Split(storage.offset)
	return rightSlice.String()
}

func (storage *DnaStorageRope) PrependPrefix(s string) {
	leftRope := rope.NewString(s)
	_, rightRope := storage.data.Split(storage.offset)
	storage.data = leftRope.Append(rightRope)
	storage.offset = 0
}

func (storage *DnaStorageRope) Len() int {
	rest := storage.data.Length() - storage.offset
	if rest > 0 {
		return rest
	}
	return 0
}

func (storage *DnaStorageRope) SaveOffset() {
	storage.savedOffset = storage.offset
}

func (storage *DnaStorageRope) RestoreOffset() {
	storage.offset = storage.savedOffset
}

func (storage *DnaStorageRope) GetString(n int) (string, error) {
	oldOffset := storage.offset
	storage.offset = oldOffset + n

	if storage.offset >= storage.data.Length() {
		return "", errors.New("not match")
	}

	leftRope, _ := storage.data.Split(storage.offset)

	return leftRope.String(), nil
}
