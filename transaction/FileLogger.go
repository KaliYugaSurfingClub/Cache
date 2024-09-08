package transaction

import (
	"fmt"
	"os"
)

type FileLogger struct {
	events       chan<- Event
	errs         <-chan error
	file         *os.File
	lastSequence uint64
}

func NewFileLogger(filename string) (Logger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	return &FileLogger{file: file}, nil
}

func (tl *FileLogger) WritePut(key string, value []byte) {
	tl.events <- Event{Type: EventPut, Key: key, Value: value}
}

func (tl *FileLogger) WriteDelete(key string) {
	tl.events <- Event{Type: EventDelete, Key: key}
}

func (tl *FileLogger) ErrCh() <-chan error {
	return tl.errs
}

func (tl *FileLogger) Start() {
	events := make(chan Event) //todo buffer 16
	tl.events = events

	errs := make(chan error) //todo buffer 1
	tl.errs = errs

	go func() {
		//always read from events channel, Somebody who write to this channel is
		//responsible for closing it at the right time
		for e := range events {
			tl.lastSequence++ //todo first log with id = 1 and maybe current instead last
			e.Sequence = tl.lastSequence

			eventBytes, err := encodeEvent(&e)
			if err != nil {
				errs <- err
			}

			fmt.Println(eventBytes)

			_, err = tl.file.Write(eventBytes)

			if err != nil { //todo this errors do not handled
				errs <- err
				return
			}
		}
	}()
}

func (tl *FileLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event)
	//todo buffer 1
	outError := make(chan error)

	go func() {
		//this goroutine writes to this channel and responsible for closing it
		//after it writes everything
		defer close(outError)
		defer close(outEvent)
		//deadlock without close

		//todo do not read all file read to delim and send event
		eventsBytes, err := os.ReadFile("logs.bin")
		if err != nil {
			outError <- fmt.Errorf("error reading logs.bin: %s", err)
			return
		}

		decodeEvents(eventsBytes, outEvent, outError)
	}()

	return outEvent, outError
}
