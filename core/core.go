package core

import (
	"errors"
	"sync"
)

var ErrorNoSuchKey = errors.New("no such key")

type Store struct {
	sync.RWMutex
	data map[string][]byte
}

var store = Store{data: make(map[string][]byte)}

// todo return error?
func Put(key string, value []byte) error {
	store.Lock()
	defer store.Unlock()

	store.data[key] = value
	return nil
}

func Get(key string) ([]byte, error) {
	//произвольное кол-во может удерживать RLock и читать
	//но если кто-то держит Lock то это функция будет ждать пока блокировка на запись закончиться
	store.RLock()
	defer store.RUnlock()

	value, ok := store.data[key]
	if !ok {
		return nil, ErrorNoSuchKey
	}

	return value, nil
}

// todo return error?
func Delete(key string) error {
	store.Lock()
	defer store.Unlock()

	delete(store.data, key)
	return nil
}
