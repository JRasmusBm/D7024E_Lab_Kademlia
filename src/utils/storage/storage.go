package storage

import (
	"fmt"
	"sync"
)

type Storage interface {
	Read(key string, ch chan string, errCh chan error)
	Write(key, data string)
}

type RealStorage struct {
	sync.RWMutex
	Data map[string]string
}

func (store RealStorage) Read(key string, ch chan string, errCh chan error) {
	store.RLock()
	value, ok := store.Data[key]
	store.RUnlock()
	if !ok {
		errCh <- fmt.Errorf("%s is not a valid key", key)
		return
	}
	ch <- value
}

func (store RealStorage) Write(key, data string) {
	store.Lock()
	store.Data[key] = data
	store.Unlock()
}
