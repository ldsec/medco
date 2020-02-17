package survivalserver

import (
	"errors"
	"fmt"
	"time"

	"go.dedis.ch/onet/v3/log"
)

//for the entire survival query, maybe wrap this around a query anin a structure that is created for each new survival query..
const MaxSize = 1024

type ExceptionHandler struct {
	errorChannel        chan error
	bufferSize          int
	endOfProcessChannel chan struct{}
}

func NewExceptionHandler(bufferSize int) (res *ExceptionHandler, err error) {
	if bufferSize < 1 {
		err = fmt.Errorf(`The size of the buffered channel must be at least 1, provided %d`, bufferSize)
		return
	}
	if bufferSize > MaxSize {
		err = fmt.Errorf(`%d exceeds max size of buffered channel %d`, bufferSize, MaxSize)
		return
	}

	errorChan := make(chan error, bufferSize)
	endOfProcessChan := make(chan struct{})
	res = &ExceptionHandler{
		errorChannel:        errorChan,
		bufferSize:          bufferSize,
		endOfProcessChannel: endOfProcessChan,
	}
	return

}

//var errorChannel = make(chan error)
//var endOfProcessChannel = make(chan bool)

func (handler *ExceptionHandler) PushError(err error) {
	select {
	case handler.errorChannel <- err:
		log.Lvl2("pushed an error")
	default:
		log.Lvl2("error channels already full")
	}
}
func (handler *ExceptionHandler) Finished() {
	handler.endOfProcessChannel <- struct{}{}
}

func (handler *ExceptionHandler) WaitEndSignal(timeoutInSeconds int) (err error) {
	select {
	case err = <-handler.errorChannel:
		return
	case <-handler.endOfProcessChannel:
		return
	case <-time.After(time.Duration(timeoutInSeconds) * time.Second):
		err = errors.New("Survival query Timeout")
		return
	}

}
