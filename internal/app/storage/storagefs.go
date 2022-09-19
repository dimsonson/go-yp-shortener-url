package storage

import (
	"fmt"
)

type StorageFs struct {
	IdURL map[string]string `json:"idurl,omitempty"`
}

var d = StorageFs{
	IdURL: make(map[string]string),
}

func (ms *StorageFs) PutStorage(key string, value string) (err error) {
	// запись в JSON

	if _, ok := d.IdURL[key]; ok {
		return fmt.Errorf("key is already in database")
	}
	d.IdURL[key] = string(value)
	return nil
}

func NewFsStorage(s map[string]string) *StorageFs {
	// загрузка базы из JSON

	return &StorageFs{
		IdURL: s,
	}
}

func (ms *StorageFs) GetStorage(key string) (value string, err error) {
	value, ok := d.IdURL[key]
	if !ok {
		return "", fmt.Errorf("key %v not found", key)
	}
	return value, nil
}

func (ms *StorageFs) LenStorage() (lenn int) {
	lenn = len(d.IdURL)
	return lenn
}
