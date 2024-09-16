package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrorNoSuchKey = errors.New("no such key")

type TransactionLogger interface {
	WriteEvent(t EventType, key string, value string)
	ReadEvents() (<-chan Event, <-chan error)
	Start() <-chan error
	Shutdown(ctx context.Context) error
}

type Store struct {
	sync.RWMutex
	data map[string]string
	tl   TransactionLogger
}

func NewStore(tl TransactionLogger) *Store {
	return &Store{
		data: make(map[string]string),
		tl:   tl,
	}
}

func (s *Store) WithTransactionLogger(tl TransactionLogger) *Store {
	s.tl = tl
	return s
}

func (s *Store) Get(key string) (string, error) {
	s.RLock()
	defer s.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func (s *Store) Put(key string, value string) {
	s.Lock()
	defer s.Unlock()

	s.data[key] = value
	s.tl.WriteEvent(EventPut, key, value)
}

func (s *Store) Delete(key string) {
	s.Lock()
	defer s.Unlock()

	delete(s.data, key)
	s.tl.WriteEvent(EventDelete, key, "")
}

func (s *Store) Restore() error {
	var err error

	s.Lock()
	defer s.Unlock()

	events, errs := s.tl.ReadEvents()
	ok, event := true, Event{}

	for ok && err == nil {
		select {
		case err, ok = <-errs:
		case event, ok = <-events:
			switch event.Type {
			case EventPut:
				s.data[event.Key] = event.Value
			case EventDelete:
				delete(s.data, event.Key)
			}
		}
	}

	fmt.Println("store", s.data)

	return err
}
