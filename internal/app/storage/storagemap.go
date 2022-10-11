package storage

import (
	"fmt"
)

// структура хранилища в памяти
type StorageMap struct {
	UserID map[string]int 
	IDURL  map[string]string
}

// метод записи в хранилище в памяти
func (ms *StorageMap) PutToStorage(userid int, key string, value string) (err error) {
	if _, ok := ms.IDURL[key]; ok {
		return fmt.Errorf("key is already in database")
	}
	ms.IDURL[key] = string(value)
	ms.UserID[key] = userid
	return nil
}

// конструктор хранилища в памяти
func NewMapStorage(u map[string]int, s map[string]string) *StorageMap {
	return &StorageMap{
		UserID: u,
		IDURL:  s,
	}
}

// метод получения id:url из хранилища в памяти
func (ms *StorageMap) GetFromStorage(key string) (value string, err error) {
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

// метод отбора URLs по UserID
func (ms *StorageMap) URLsByUserID(userid int) (userURLs map[string]string, err error) {
	userURLs = make(map[string]string)
	for k, v := range ms.UserID {
		if v == userid {
			userURLs[k] = ms.IDURL[k]
		}
	}
	if len(userURLs) == 0 {
		err = fmt.Errorf("userid not found in the storage")
	}
	return userURLs, err
}

func (ms *StorageMap) LoadFromFileToStorage() {

}

// посик userid в хранилице
func (ms *StorageMap) UserIDExist(userid int) bool {
	// цикл по map поиск значения без ключа
	for _, v := range ms.UserID {
		if v == userid {
			return true
		}
	}
	return false
}
