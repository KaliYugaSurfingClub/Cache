package core

import (
	"errors"
	"sync"
)

var ErrorNoSuchKey = errors.New("no such key")

type Store struct {
	sync.RWMutex
	data map[string]string
}

var store = Store{data: make(map[string]string)}

func Put(key string, value string) error {
	store.Lock()
	defer store.Unlock()

	store.data[key] = value
	return nil
}

func Get(key string) (string, error) {
	//произвольное кол-во может удерживать RLock и читать
	//но если кто-то держит Lock то это функция будет ждать пока блокировка на запись закончиться
	store.RLock()
	defer store.RUnlock()

	value, ok := store.data[key]
	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func Delete(key string) error {
	store.Lock()
	defer store.Unlock()

	delete(store.data, key)
	return nil
}
