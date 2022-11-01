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
}

// метод записи в хранилище в памяти
func (ms *StorageMap) PutToStorage(ctx context.Context, key string, value string, userid string) (existKey string, err error) {
	// получаем значение iserid из контекста
	//userid := ctx.Value(settings.CtxKeyUserID).(string)
	ms.IDURL[key] = string(value)
	ms.UserID[key] = userid
	existKey = key
	return existKey, err
}

// конструктор хранилища в памяти
func NewMapStorage(u map[string]string, s map[string]string) *StorageMap {
	return &StorageMap{
		UserID: u,
		IDURL:  s,
	}
}

// метод получения id:url из хранилища в памяти
func (ms *StorageMap) GetFromStorage(ctx context.Context, key string) (value string, err error) {
	// метод получения записи из хранилища
	value, ok := ms.IDURL[key]
	if !ok {
		return "", fmt.Errorf("key %v not found", key)
	}
	return value, nil
}

// метод определения длинны хранилища
func (ms *StorageMap) LenStorage(ctx context.Context) (lenn int) {
	lenn = len(ms.IDURL)
	return lenn
}

// метод отбора URLs по UserID
func (ms *StorageMap) URLsByUserID(ctx context.Context, userid string) (userURLs map[string]string, err error) {
	// получаем значение iserid из контекста
	// userid := ctx.Value(settings.CtxKeyUserID).(string)
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
func (ms *StorageMap) UserIDExist(ctx context.Context, userid string) bool {
	// цикл по map поиск значения без ключа
	for _, v := range ms.UserID {
		if v == userid {
			return true
		}
	}
	return false
}

func (ms *StorageMap) StorageOkPing(ctx context.Context) (bool, error) {

	return true, nil
}

func (ms *StorageMap) StorageConnectionClose() {

}

// метод пакетной записи id:url в хранилище
func (ms *StorageMap) PutBatchToStorage(ctx context.Context, dc settings.DecodeBatchJSON, userid string) (dcCorr settings.DecodeBatchJSON, err error) {
	// userid := ctx.Value(settings.CtxKeyUserID).(string)
	for _, v := range dc {
		// записываем в хранилице userid, id, URL
		ms.IDURL[v.ShortURL] = userid
		ms.UserID[v.ShortURL] = v.OriginalURL
	}
	return dc, err
}
