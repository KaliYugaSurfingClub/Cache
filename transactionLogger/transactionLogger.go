package transactionLogger

import (
	"bufio"
	"cache/core"
	"fmt"
	"os"
)

type TransactionLogger interface {
	WritePut(key, value string)
	WriteDelete(key string)
	ReadEvents() error
	Start()
	ErrCh() <-chan error
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

func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	return &FileTransactionLogger{file: file}, nil
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

func (tl *FileTransactionLogger) Start() {
	//todo buffer 16
	events := make(chan Event)
	tl.events = events

	//todo buffer 1
	errs := make(chan error)
	tl.errs = errs

	go func() {
		//always read from events channel, Somebody who write to this channel is
		//responsible for closing it at the right time
		for e := range events {
			//todo first log with id = 1 and maybe current instead last
			tl.lastSequence++

			_, err := fmt.Fprintf(
				tl.file, "%d\t%d\t%s\t%s\n",
				tl.lastSequence, e.Type, e.Key, e.Value,
			)
			//todo this errors do not handled
			if err != nil {
				errs <- err
				return
			}
		}
	}()
}

func (tl *FileTransactionLogger) ReadEvents() error {
	scanner := bufio.NewScanner(tl.file)

	for scanner.Scan() {
		line := scanner.Text()

		var e Event
		//todo EOF error if last log is delete log
		_, err := fmt.Sscanf(
			line, "%d\t%d\t%s\t%s",
			&e.Sequence, &e.Type, &e.Key, &e.Value,
		)

		if err != nil {
			return err
		}

		if tl.lastSequence >= e.Sequence {
			return err
		}

		tl.lastSequence = e.Sequence
		switch e.Type {
		case EventPut:
			err = core.Put(e.Key, e.Value)
		case EventDelete:
			err = core.Delete(e.Key)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
