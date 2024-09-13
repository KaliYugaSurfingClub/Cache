package transaction

import (
	"cache/core"
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
)

type FileLogger struct {
	wg *sync.WaitGroup

	events    chan<- core.Event
	file      *os.File
	currentID uint64
}

type BackgroundLogger struct{}

func (tl *BackgroundLogger) WritePut(key string, value []byte)             {}
func (tl *BackgroundLogger) WriteDelete(key string)                        {}
func (tl *BackgroundLogger) ReadEvents() (<-chan core.Event, <-chan error) { return nil, nil }
func (tl *BackgroundLogger) Start() <-chan error                           { return nil }
func (tl *BackgroundLogger) Shutdown(ctx context.Context) error            { return nil }

func NewLogger(filename string) (core.TransactionLogger, error) {
	if filename == "" {
		return &BackgroundLogger{}, nil
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	return &FileLogger{file: file, wg: &sync.WaitGroup{}}, nil
}

// todo what will happens if i will write at the shutdown time
func (tl *FileLogger) WritePut(key string, value []byte) {
	tl.wg.Add(1)
	tl.events <- core.Event{Type: core.EventPut, Key: key, Value: value}
}

func (tl *FileLogger) WriteDelete(key string) {
	tl.wg.Add(1)
	tl.events <- core.Event{Type: core.EventDelete, Key: key}
}

func (tl *FileLogger) Wait() {
	tl.wg.Wait()
}

func (tl *FileLogger) ReadEvents() (<-chan core.Event, <-chan error) {
	outEvent := make(chan core.Event)
	outError := make(chan error)

	go func() {
		//this goroutine writes to this channel and responsible for closing it
		//after it writes everything
		defer close(outError)
		defer close(outEvent)

		for {
			event, err := readEvent(tl.file)
			////todo debug
			//time.Sleep(1 * time.Minute)

			if errors.Is(err, ErrEmptyFile) {
				return
			}

			fmt.Println(event)

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

func (tl *FileLogger) Start() <-chan error {
	//buffer 16 means that 16 handlers can send event and do not wait when logger write event to file
	//todo remove literal
	events := make(chan core.Event, 16)
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

func (tl *FileLogger) Shutdown(ctx context.Context) error {
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
