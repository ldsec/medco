package survivalserver

import (
	"errors"
	"fmt"
	"time"

	"go.dedis.ch/onet/v3/log"
)

//ExceptionHandler joins multiple parallel threads that can each returns an error
type ExceptionHandler struct {
	errorChannel        chan error
	bufferSize          int
	endOfProcessChannel chan struct{}
}

//NewExceptionHandler ExceptionHandler constructor
func NewExceptionHandler(bufferSize int) (res *ExceptionHandler, err error) {
	if bufferSize < 1 {
		err = fmt.Errorf(`The size of the buffered channel must be at least 1, provided %d`, bufferSize)
		return
	}
	if bufferSize > errorChanMaxSize {
		err = fmt.Errorf(`%d exceeds max size of buffered channel %d`, bufferSize, errorChanMaxSize)
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

//PushError push the error if th buffered channel is not full, does nothing else
func (handler *ExceptionHandler) PushError(err error) {
	select {
	case handler.errorChannel <- err:
		log.Lvl2("pushed an error")
	default:
		log.Lvl2("error channels already full")
	}
}

//Finished indicates that all executions in threads are finished
func (handler *ExceptionHandler) Finished() {
	handler.endOfProcessChannel <- struct{}{}
}

//WaitEndSignal waits until either the
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
