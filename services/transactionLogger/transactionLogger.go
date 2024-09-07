package transactionLogger

import (
	"fmt"
	"os"
)

type TransactionLogger interface {
	WritePut(key, value string)
	WriteDelete(key string)
}

type EventType byte

const (
	EventDelete EventType = iota
	EventPut
)

type Event struct {
	Sequence uint64
	Type     EventType
	Key      string
	Value    string
}

type FileTransactionLogger struct {
	events       chan<- Event
	errs         <-chan error
	file         *os.File
	lastSequence uint64
}

func NewFileTransactionLogger(filename string) (*FileTransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	return &FileTransactionLogger{file: file}, nil
}

func (tl *FileTransactionLogger) Start() {
	events := make(chan Event)
	tl.events = events

	errs := make(chan error)
	tl.errs = errs

	go func() {
		for e := range events {
			tl.lastSequence++

			_, err := fmt.Fprintf(
				tl.file, "%d\t%d\t%s\t%s\n",
				tl.lastSequence, e.Type, e.Key, e.Value,
			)
			if err != nil {
				errs <- err
				return
			}
		}
	}()
}

func (tl *FileTransactionLogger) WritePut(key, value string) {
	tl.events <- Event{Type: EventPut, Key: key, Value: value}
}

func (tl *FileTransactionLogger) WriteDelete(key string) {
	tl.events <- Event{Type: EventDelete, Key: key}
}

func (tl *FileTransactionLogger) ErrCh() <-chan error {
	return tl.errs
}
