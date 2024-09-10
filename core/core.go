package core

import (
	"errors"
	"sync"
)

var ErrorNoSuchKey = errors.New("no such key")

type EventType byte

const (
	EventDelete EventType = iota
	EventPut
)

type Event struct {
	Sequence uint64
	Type     EventType
	Key      string
	Value    []byte
}

type transactionLogger interface {
	WritePut(key string, value []byte)
	WriteDelete(key string)
	ReadEvents() (<-chan Event, <-chan error)
	Close() error
	Start()
}

type Store struct {
	sync.RWMutex
	data map[string][]byte
	tl   transactionLogger
}

func NewStore() *Store {
	return &Store{
		data: make(map[string][]byte),
		tl:   &ZeroLogger{},
	}
}

// todo maybe refactor
func (s *Store) WithTransactionLogger(tl transactionLogger) *Store {
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

	s.data[key] = value
	s.tl.WritePut(key, value)
}

func (s *Store) Delete(key string) {
	s.Lock()
	defer s.Unlock()

	delete(s.data, key)
	s.tl.WriteDelete(key)
}

func (s *Store) Restore() error {
	events, errs := s.tl.ReadEvents()

	var ok = true
	var err error = nil
	var event Event

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

	s.tl.Start()

	return err
}

type ZeroLogger struct{}

func (tl *ZeroLogger) WritePut(key string, value []byte)        {}
func (tl *ZeroLogger) WriteDelete(key string)                   {}
func (tl *ZeroLogger) ReadEvents() (<-chan Event, <-chan error) { return nil, nil }
func (tl *ZeroLogger) ErrCh() <-chan error                      { return nil }
func (tl *ZeroLogger) Start()                                   {}
func (tl *ZeroLogger) Close() error                             { return nil }
