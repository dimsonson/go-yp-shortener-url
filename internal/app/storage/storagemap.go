package storage

import (
	"fmt"
)

// структура хранилища в памяти
type StorageMap struct {
	IDURL map[string]string
}

// метод записи в хранилище в памяти
func (ms *StorageMap) PutToStorage(key string, value string) (err error) {
	if _, ok := ms.IDURL[key]; ok {
		return fmt.Errorf("key is already in database")
	}
	ms.IDURL[key] = string(value)
	return nil
}

// конструктор хранилища в памяти
func NewMapStorage(s map[string]string) *StorageMap {
	return &StorageMap{
		IDURL: s,
	}
}

// метод получения id:url из хранилища в памяти
func (ms *StorageMap) GetStorage(key string) (value string, err error) {
	// метод получения записи из хранилища
	value, ok := ms.IDURL[key]
	if !ok {
		return "", fmt.Errorf("key %v not found", key)
	}
	return value, nil
}

// метод определения длинны хранилища
func (ms *StorageMap) LenStorage() (lenn int) {
	lenn = len(ms.IDURL)
	return lenn
}
