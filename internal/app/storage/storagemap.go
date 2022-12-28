package storage

import (
	"context"
	"fmt"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
)

// StorageMap структура хранилища в памяти.
type StorageMap struct {
	UserID map[string]string
	IDURL  map[string]string
	DelURL map[string]bool
}

// Put метод записи в хранилище в памяти.
func (ms *StorageMap) Put(ctx context.Context, key string, value string, userid string) (existKey string, err error) {

	ms.IDURL[key] = string(value)
	ms.UserID[key] = userid
	ms.DelURL[key] = false
	existKey = key
	return existKey, err
}

// NewMapStorage конструктор хранилища в памяти.
func NewMapStorage(u map[string]string, s map[string]string, d map[string]bool) *StorageMap {
	return &StorageMap{
		UserID: u,
		IDURL:  s,
		DelURL: d,
	}
}

// Get метод получения id:url из хранилища в памяти.
func (ms *StorageMap) Get(ctx context.Context, key string) (value string, del bool, err error) {
	// метод получения записи из хранилища
	value, ok := ms.IDURL[key]
	if !ok {
		return "", false, fmt.Errorf("key %v not found", key)
	}
	del = ms.DelURL[key]
	return value, del, nil
}

// Len метод определения длинны хранилища.
func (ms *StorageMap) Len(ctx context.Context) (lenn int) {
	lenn = len(ms.IDURL)
	return lenn
}

// GetBatch метод отбора URLs по UserID.
func (ms *StorageMap) GetBatch(ctx context.Context, userid string) (userURLs map[string]string, err error) {

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

// Load метод загрузки хранилища в кеш при инциализации файлового хранилища.
func (ms *StorageMap) Load() {

}
// Ping метод проверки доступности SQL хранилища.
func (ms *StorageMap) Ping(ctx context.Context) (bool, error) {
	return true, nil
}
// Close метод закрытия соединения доступности SQL хранилища.
func (ms *StorageMap) Close() {

}

// PutBatch метод пакетной записи id:url в хранилище в памяти.
func (ms *StorageMap) PutBatch(ctx context.Context, dc models.BatchRequest, userid string) (dcCorr models.BatchRequest, err error) {
	// итерируем по слайсу
	for _, v := range dc {
		// записываем в хранилице userid, id, URL
		ms.IDURL[v.ShortURL] = userid
		ms.UserID[v.ShortURL] = v.OriginalURL
		ms.DelURL[v.ShortURL] = false
	}
	return dc, err
}
// Delete метод пометки записи в хранилище в памяти как удаленной.
func (ms *StorageMap) Delete(key string, userid string) (err error) {
	ms.IDURL[key] = userid
	ms.DelURL[key] = true
	return nil
}
