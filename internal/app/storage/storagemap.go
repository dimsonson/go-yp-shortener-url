package storage

import (
	"fmt"
)

type StorageMap struct {
	IDURL map[string]string
}

var m = StorageMap{
	IDURL: make(map[string]string),
}

func (ms *StorageMap) PutStorage(key string, value string) (err error) {
	if _, ok := m.IDURL[key]; ok {
		return fmt.Errorf("key is already in database")
	}
	m.IDURL[key] = string(value)
	return nil
}

func NewMapStorage(s map[string]string) *StorageMap {
	return &StorageMap{
		IDURL: s,
	}
}

func (ms *StorageMap) GetStorage(key string) (value string, err error) {
	value, ok := m.IDURL[key]
	if !ok {
		return "", fmt.Errorf("key %v not found", key)
	}
	return value, nil
}

func (ms *StorageMap) LenStorage() (lenn int) {
	lenn = len(m.IDURL)
	return lenn
}
