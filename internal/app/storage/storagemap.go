package storage

import (
	"fmt"
)

type StorageMap struct {
	IdURL map[string]string
}

var m = StorageMap{
	IdURL: make(map[string]string),
}

func (ms *StorageMap) PutStorage(key string, value string) (err error) {
	if _, ok := m.IdURL[key]; ok {
		return fmt.Errorf("key is already in database")
	}
	m.IdURL[key] = string(value)
	return nil
}

func NewMapStorage(s map[string]string) *StorageMap {
	return &StorageMap{
		IdURL: s,
	}
}

func (ms *StorageMap) GetStorage(key string) (value string, err error) {
	value, ok := m.IdURL[key]
	if !ok {
		return "", fmt.Errorf("key %v not found", key)
	}
	return value, nil
}

func (ms *StorageMap) LenStorage() (lenn int) {
	lenn = len(m.IdURL)
	return lenn
}
