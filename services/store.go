package services

import (
	"cache/services/transactionLogger"
	"errors"
	"sync"
)

var ErrorNoSuchKey = errors.New("no such key")

type Store struct {
	sync.RWMutex
	data map[string]string
}

var store = Store{data: make(map[string]string)}

func Get(key string) (string, error) {
	store.RLock()
	defer store.RUnlock()

	value, ok := store.data[key]
	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func Put(key string, value string, tl transactionLogger.TransactionLogger) error {
	store.Lock()
	defer store.Unlock()

	tl.WritePut(key, value)

	store.data[key] = value
	return nil
}

func Delete(key string, tl transactionLogger.TransactionLogger) error {
	store.Lock()
	defer store.Unlock()

	tl.WriteDelete(key)

	delete(store.data, key)
	return nil
}
