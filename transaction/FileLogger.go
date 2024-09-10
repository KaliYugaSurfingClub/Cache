package transaction

import (
	"cache/core"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

type Logger struct {
	wg *sync.WaitGroup

	events    chan<- core.Event
	file      *os.File
	currentID uint64
}

func NewFileLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	return &Logger{file: file, wg: &sync.WaitGroup{}}, nil
}

func (tl *Logger) WritePut(key string, value []byte) {
	tl.wg.Add(1)
	tl.events <- core.Event{Type: core.EventPut, Key: key, Value: value}
}

func (tl *Logger) WriteDelete(key string) {
	tl.wg.Add(1)
	tl.events <- core.Event{Type: core.EventDelete, Key: key}
}

func (tl *Logger) Wait() {
	tl.wg.Wait()
}

func (tl *Logger) ReadEvents() (<-chan core.Event, <-chan error) {
	outEvent := make(chan core.Event)
	outError := make(chan error)

	go func() {
		//this goroutine writes to this channel and responsible for closing it
		//after it writes everything
		defer close(outError)
		defer close(outEvent)
		//deadlock without close

		for {
			event, err := readEvent(tl.file)
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				outError <- err
			}

			if tl.currentID > event.ID {
				outError <- fmt.Errorf("transaction numbers out of sequence")
			}

			tl.currentID = event.ID
			outEvent <- event
		}
	}()

	return outEvent, outError
}

func (tl *Logger) Start() <-chan error {
	events := make(chan core.Event)
	errs := make(chan error)

	tl.events = events

	go func() {
		//todo defer close(events, errs)
		//always read from events channel, Somebody who write to this channel is
		//responsible for closing it at the right time
		for e := range events {
			e.ID = tl.currentID
			tl.currentID++

			err := writeEvent(tl.file, e)
			tl.wg.Done()

			if err != nil {
				errs <- err
				//todo finish with context
				return
			}
		}
	}()

	return errs
}

func (tl *Logger) Close() error {
	fmt.Println("close logger")

	tl.Wait()

	if tl.events != nil {
		close(tl.events)
	}

	return tl.file.Close()
}
