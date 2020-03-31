package survivalclient

import (
	"errors"
	"fmt"
	"sync"
)

type BufferToPrint struct {
	textBuffer []byte
	lock       sync.Mutex
}

func (buff *BufferToPrint) Write(p []byte) (int, error) {

	buff.lock.Lock()
	defer func() { buff.lock.Unlock() }()
	currentLen := len(buff.textBuffer)
	expectedLen := currentLen + len(p)
	buff.textBuffer = append(buff.textBuffer, p...)
	if newLen := len(buff.textBuffer); newLen != expectedLen {
		return newLen - currentLen, errors.New("error while appending bytes")
	}
	return expectedLen - currentLen, nil

}

func (buff *BufferToPrint) Print() {
	buff.lock.Lock()
	defer func() { buff.lock.Unlock() }()
	fmt.Print(string(buff.textBuffer))
}
