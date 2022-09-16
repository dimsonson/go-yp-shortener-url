package storage

import (
	"fmt"
)

type Storage struct {
	Type string
}

var DB = make(map[string]string)

func (ms *Storage) PutStorage(key string, value string) (err error) {
	if _, ok := DB[key]; ok {
		return fmt.Errorf("key is already in database")
	}
	DB[key] = string(value)
	return nil
}

func NewMapStorage(s string) *Storage {
	return &Storage{
		Type: s,
	}
}

func (ms *Storage) GetStorage(key string) (value string, err error) {
	value, ok := DB[key]
	if !ok {
		return "", fmt.Errorf("key %v not found", key)
	}
	return value, nil
}

func (ms *Storage) LenStorage() (lenn int) {
	lenn = len(DB)
	return lenn
}
