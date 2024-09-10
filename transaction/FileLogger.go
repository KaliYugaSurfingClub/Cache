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

	events       chan<- core.Event
	errs         <-chan error
	file         *os.File
	lastSequence uint64
}

func NewFileLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	return &Logger{file: file, wg: &sync.WaitGroup{}}, nil
}

// todo only one method write
func (tl *Logger) WritePut(key string, value []byte) {
	tl.wg.Add(1)
	tl.events <- core.Event{Type: core.EventPut, Key: key, Value: value}
}

func (tl *Logger) WriteDelete(key string) {
	tl.wg.Add(1)
	tl.events <- core.Event{Type: core.EventDelete, Key: key}
}

func (tl *Logger) ErrCh() <-chan error {
	return tl.errs
}

func (tl *Logger) Wait() {
	tl.wg.Wait()
}

func (tl *Logger) Start() {
	events := make(chan core.Event)
	errs := make(chan error)

	tl.events = events
	tl.errs = errs

	go func() {
		//always read from events channel, Somebody who write to this channel is
		//responsible for closing it at the right time
		for e := range events {
			tl.lastSequence++ //todo first log with id = 1 and maybe do current instead last
			e.Sequence = tl.lastSequence

			//todo catch critical errors and stop non-critical continue
			//todo HARD
			err := writeEvent(tl.file, e)

			if err != nil {
				errs <- err
				return
			}

			tl.wg.Done()
		}
	}()
}

func (tl *Logger) Close() error {
	tl.Wait()

	if tl.events != nil {
		close(tl.events)
	}

	return tl.file.Close()
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
				return
			}

			if tl.lastSequence >= event.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			tl.lastSequence = event.Sequence
			outEvent <- event
		}
	}()

	return outEvent, outError
}
