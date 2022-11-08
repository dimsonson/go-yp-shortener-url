package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
)

// структура хранилища
type StorageFile struct {
	UserID   map[string]string `json:"iserid,omitempty"` // shorturl:userid
	IDURL    map[string]string `json:"idurl,omitempty"`  // shorturl:URL
	DelURL   map[string]bool   `json:"_"`                // shorturl:deleted_url
	pathName string
}

// метод записи id:url в хранилище
func (ms *StorageFile) StoragePut(ctx context.Context, key string, value string, userid string) (existKey string, err error) {

	// записываем в хранилице userid, id, URL
	ms.IDURL[key] = value
	ms.UserID[key] = userid
	ms.DelURL[key] = false
	existKey = key
	// открываем файл
	sfile, err := os.OpenFile(ms.pathName, os.O_WRONLY, 0777)
	if err != nil {
		log.Println("storage file opening error: ", err)
		return //err
	}
	defer sfile.Close()
	// кодирование в JSON
	js, err := json.Marshal(&ms)
	if err != nil {
		log.Println("JSON marshalling from struct error: ", err)
		return //err
	}
	// запись в файл
	sfile.Write(js)
	return existKey, err
}

// конструктор нового хранилища JSON
func NewFileStorage(u map[string]string, s map[string]string, d map[string]bool, p string) *StorageFile {

	return &StorageFile{
		UserID:   u,
		IDURL:    s,
		DelURL:   d,
		pathName: p,
	}
}

// метод получения записи из хранилища
func (ms *StorageFile) StorageGet(ctx context.Context, key string) (value string, del bool, err error) {
	value, ok := ms.IDURL[key]
	if !ok {
		return "", false, fmt.Errorf("key %v not found", key)
	}
	del = ms.DelURL[key]
	return value, del, nil
}

// метод определения длинны хранилища
func (ms *StorageFile) StorageLen(ctx context.Context) (lenn int) {
	lenn = len(ms.IDURL)
	return lenn
}

// метод отбора URLs по UserID
func (ms *StorageFile) StorageURLsByUserID(ctx context.Context, userid string) (userURLs map[string]string, err error) {

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

func (ms *StorageFile) StorageLoadFromFile() {
	// загрузка базы из JSON
	p := ms.pathName
	_, pathOk := os.Stat(filepath.Dir(p))
	if os.IsNotExist(pathOk) {
		os.MkdirAll(filepath.Dir(p), 0777)
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

func (ms *StorageFile) StorageOkPing(ctx context.Context) (bool, error) {

	return true, nil
}

func (ms *StorageFile) StorageConnectionClose() {

}

// метод пакетной записи id:url в хранилище
func (ms *StorageFile) StoragePutBatch(ctx context.Context, dc models.BatchRequest, userid string) (dcCorr models.BatchRequest, err error) {
	// итерируем по слайсу
	for _, v := range dc {
		// записываем в хранилице userid, id, URL, del
		ms.IDURL[v.ShortURL] = userid
		ms.UserID[v.ShortURL] = v.OriginalURL
		ms.DelURL[v.ShortURL] = false
	}
	return dc, err
}

func (ms *StorageFile) StorageDeleteURL(key string, userid string) (err error) {
	ms.IDURL[key] = userid
	ms.DelURL[key] = true
	return nil
}
