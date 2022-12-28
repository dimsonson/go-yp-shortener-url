// storage пакет хранилища.
package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
)

// StorageFile структура файлового хранилища.
type StorageFile struct {
	UserID   map[string]string `json:"iserid,omitempty"` // shorturl:userid
	IDURL    map[string]string `json:"idurl,omitempty"`  // shorturl:URL
	DelURL   map[string]bool   `json:"_"`                // shorturl:deleted_url
	pathName string
	mu       sync.RWMutex
}

// Put метод записи id:url в файловое хранилище.
func (ms *StorageFile) Put(ctx context.Context, key string, value string, userid string) (existKey string, err error) {
	// записываем в хранилице userid, id, URL
	ms.mu.Lock()
	ms.IDURL[key] = value
	ms.UserID[key] = userid
	ms.DelURL[key] = false
	existKey = key
	// открываем файл
	sfile, err := os.OpenFile(ms.pathName, os.O_WRONLY, 0777)
	if err != nil {
		log.Println("storage file opening error: ", err)
		return "", err
	}
	defer sfile.Close()
	// кодирование в JSON
	js, err := json.Marshal(ms)
	if err != nil {
		log.Println("JSON marshalling from struct error: ", err)
		return //err
	}
	// запись в файл
	_, err = sfile.Write(js)
	if err != nil {
		return "", err
	}
	ms.mu.Unlock()
	return existKey, err
}

// NewFileStorage конструктор нового файлового хранилища.
func NewFileStorage(u map[string]string, s map[string]string, d map[string]bool, p string) *StorageFile {
	return &StorageFile{
		UserID:   u,
		IDURL:    s,
		DelURL:   d,
		pathName: p,
	}
}

// Get метод получения записи из файлового хранилища.
func (ms *StorageFile) Get(ctx context.Context, key string) (value string, del bool, err error) {
	value, ok := ms.IDURL[key]
	if !ok {
		return "", false, fmt.Errorf("key %v not found", key)
	}
	del = ms.DelURL[key]
	return value, del, nil
}

// Len метод определения длинны файлового хранилища.
func (ms *StorageFile) Len(ctx context.Context) (lenn int) {
	lenn = len(ms.IDURL)
	return lenn
}

// GetBatch метод отбора URLs по UserID.
func (ms *StorageFile) GetBatch(ctx context.Context, userid string) (userURLs map[string]string, err error) {

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

// Load метод загрузки хранилища в кеш при инциализации.
func (ms *StorageFile) Load() {
	// загрузка базы из JSON
	p := ms.pathName
	_, pathOk := os.Stat(filepath.Dir(p))
	if os.IsNotExist(pathOk) {
		err := os.MkdirAll(filepath.Dir(p), 0777)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("folder %s created\n", filepath.Dir(p))
	}
	sfile, err := os.OpenFile(p, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal("file creating error: ", err)
	}
	defer sfile.Close()

	fileInfo, _ := os.Stat(p)
	if fileInfo.Size() != 0 {
		b, err := io.ReadAll(sfile)
		if err != nil {
			log.Println("file storage reading error:", err)
		}
		err = json.Unmarshal(b, &ms)
		if err != nil {
			log.Println("JSON unmarshalling to struct error:", err)
		}
	}
}

// Ping метод проверки доступности SQL хранилища.
func (ms *StorageFile) Ping(ctx context.Context) (bool, error) {
	return true, nil
}

// Close метод закрытия соединения доступности SQL хранилища.
func (ms *StorageFile) Close() {
}

// PutBatch метод пакетной записи id:url в хранилище.
func (ms *StorageFile) PutBatch(ctx context.Context, dc models.BatchRequest, userid string) (dcCorr models.BatchRequest, err error) {
	// итерируем по слайсу
	for _, v := range dc {
		// записываем в хранилице userid, id, URL, del
		ms.IDURL[v.ShortURL] = userid
		ms.UserID[v.ShortURL] = v.OriginalURL
		ms.DelURL[v.ShortURL] = false
	}
	return dc, err
}

// Delete метод пометки записи в файловом хранилище как удаленной.
func (ms *StorageFile) Delete(key string, userid string) (err error) {
	ms.IDURL[key] = userid
	ms.DelURL[key] = true
	return nil
}
