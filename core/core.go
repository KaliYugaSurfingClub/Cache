package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
)

var ErrorNoSuchKey = errors.New("no such key")

type EventType byte

const (
	EventDelete EventType = iota
	EventPut
)

type Event struct {
	ID    uint64
	Type  EventType
	Key   string
	Value []byte
}

type TransactionLogger interface {
	WritePut(key string, value []byte)
	WriteDelete(key string)
	ReadEvents() (<-chan Event, <-chan error)
	Start() <-chan error
	Shutdown(ctx context.Context) error
}

type Store struct {
	sync.RWMutex
	data map[string][]byte
	tl   TransactionLogger
}

func NewStore(tl TransactionLogger) *Store {
	return &Store{
		data: make(map[string][]byte),
		tl:   tl,
	}
}

func (s *Store) WithTransactionLogger(tl TransactionLogger) *Store {
	s.tl = tl
	return s
}

func (s *Store) Get(key string) ([]byte, error) {
	s.RLock()
	defer s.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return nil, ErrorNoSuchKey
	}

	return value, nil
}

func (s *Store) Put(key string, value []byte) {
	s.Lock()
	defer s.Unlock()

	////todo debug
	//fmt.Println("write", key)

	s.data[key] = value
	s.tl.WritePut(key, value)
}

func (s *Store) Delete(key string) {
	s.Lock()
	defer s.Unlock()

	delete(s.data, key)
	s.tl.WriteDelete(key)
}

func (s *Store) Start() {
	events, readingErrs := s.tl.ReadEvents()

	//todo maybe it is worth do not use goroutine for reading
	go func() {
		s.Lock()
		defer s.Unlock()

		for event := range events {
			fmt.Println("read event", event)
			switch event.Type {
			case EventPut:
				s.data[event.Key] = event.Value
			case EventDelete:
				delete(s.data, event.Key)
			}
		}
	}()

	//todo if error is critical finish readEvents
	//use context
	go func() {
		for err := range readingErrs {
			log.Println(err)
		}
	}()

	runtimeErrs := s.tl.Start()

	//todo if err is critical shutdown (maybe)
	go func() {
		for err := range runtimeErrs {
			log.Print(err)
		}
	}()
}
