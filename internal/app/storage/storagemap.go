package storage

import (
	"fmt"
)

// структура хранилища в памяти
type StorageMap struct {
	IDURL map[string]map[string]string
}

// метод записи в хранилище в памяти
func (ms *StorageMap) PutToStorage(userid string, key string, value string) (err error) {
	if _, ok := ms.IDURL[key]; ok {
		return fmt.Errorf("key is already in database")
	}
	ms.IDURL[userid][key] = value
	return nil
}

// конструктор хранилища в памяти
func NewMapStorage(s map[string]map[string]string) *StorageMap {
	return &StorageMap{
		IDURL: s,
	}
}

// метод получения id:url из хранилища в памяти
func (ms *StorageMap) GetFromStorage(key string) (value string, err error) {
	// метод получения записи из хранилища
	for _, v := range ms.IDURL {
		if value, ok := v[key]; ok {
			return value, nil
		}
	}
	return "", fmt.Errorf("key %v not found", key)
}

// метод определения длинны хранилища
func (ms *StorageMap) LenStorage() (lenn int) {
	lenn = len(ms.IDURL)
	return lenn
}
