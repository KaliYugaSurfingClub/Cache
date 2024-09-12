package transaction

import (
	"cache/core"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
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

// todo what will happens if i will write at the shutdown time
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

		for {
			event, err := readEvent(tl.file)

			//todo debug
			time.Sleep(1 * time.Minute)

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
		defer close(errs)
		//always read from events channel, Somebody who write to this channel is
		//responsible for closing it at the right time
		for e := range events {
			e.ID = tl.currentID
			tl.currentID++

			if err := writeEventTo(tl.file, e); err != nil {
				errs <- err
			}

			tl.wg.Done()
		}
	}()

	return errs
}

func (tl *Logger) Shutdown(ctx context.Context) error {
	errs := make(chan error)

	go func() {
		defer close(errs)

		tl.Wait()

		if tl.events != nil {
			close(tl.events)
		}

		if err := tl.file.Close(); err != nil {
			errs <- err
		}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("shutdown logger was cancelled: %w", ctx.Err())
	case err := <-errs:
		return fmt.Errorf("shutdown logger was faild: %w", err)
	default:
		return nil
	}
}
