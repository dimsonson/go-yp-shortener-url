package storage

import (
	"context"
	"fmt"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

// структура хранилища в памяти
type StorageMap struct {
	UserID map[string]string
	IDURL  map[string]string
	DelURL map[string]bool
}

// метод записи в хранилище в памяти
func (ms *StorageMap) StoragePut(ctx context.Context, key string, value string, userid string) (existKey string, err error) {

	ms.IDURL[key] = string(value)
	ms.UserID[key] = userid
	ms.DelURL[key] = false
	existKey = key
	return existKey, err
}

// конструктор хранилища в памяти
func NewMapStorage(u map[string]string, s map[string]string, d map[string]bool) *StorageMap {
	return &StorageMap{
		UserID: u,
		IDURL:  s,
		DelURL: d,
	}
}

// метод получения id:url из хранилища в памяти
func (ms *StorageMap) StorageGet(ctx context.Context, key string) (value string, del bool, err error) {
	// метод получения записи из хранилища
	value, ok := ms.IDURL[key]
	if !ok {
		return "", false, fmt.Errorf("key %v not found", key)
	}
	del = ms.DelURL[key]
	return value, del, nil
}

// метод определения длинны хранилища
func (ms *StorageMap) StorageLen(ctx context.Context) (lenn int) {
	lenn = len(ms.IDURL)
	return lenn
}

// метод отбора URLs по UserID
func (ms *StorageMap) StorageURLsByUserID(ctx context.Context, userid string) (userURLs map[string]string, err error) {

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

func (ms *StorageMap) StorageLoadFromFile() {

}

func (ms *StorageMap) StorageOkPing(ctx context.Context) (bool, error) {

	return true, nil
}

func (ms *StorageMap) StorageConnectionClose() {

}

// метод пакетной записи id:url в хранилище
func (ms *StorageMap) StoragePutBatch(ctx context.Context, dc settings.DecodeBatchJSON, userid string) (dcCorr settings.DecodeBatchJSON, err error) {
	// итерируем по слайсу
	for _, v := range dc {
		// записываем в хранилице userid, id, URL
		ms.IDURL[v.ShortURL] = userid
		ms.UserID[v.ShortURL] = v.OriginalURL
		ms.DelURL[v.ShortURL] = false
	}
	return dc, err
}

func (ms *StorageMap) StorageDeleteURL(key string, userid string) {
	ms.IDURL[key] = userid
	ms.DelURL[key] = true
}
