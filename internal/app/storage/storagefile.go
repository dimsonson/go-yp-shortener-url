package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

// структура хранилища
type StorageFile struct {
	UserID   map[string]string `json:"iserid,omitempty"` // shorturl:userid
	IDURL    map[string]string `json:"idurl,omitempty"`  // shorturl:URL
	DelURL   map[string]bool   `json:"_"`                // shorturl:deleted_url
	pathName string
}

// метод записи id:url в хранилище
func (ms *StorageFile) PutToStorage(ctx context.Context, key string, value string, userid string) (existKey string, err error) {
	// получаем значение iserid из контекста
	// userid := ctx.Value(settings.CtxKeyUserID).(string)
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
func (ms *StorageFile) GetFromStorage(ctx context.Context, key string) (value string, err error) {
	value, ok := ms.IDURL[key]
	if !ok {
		return "", fmt.Errorf("key %v not found", key)
	}
	return value, nil
}

// метод определения длинны хранилища
func (ms *StorageFile) LenStorage(ctx context.Context) (lenn int) {
	lenn = len(ms.IDURL)
	return lenn
}

// метод отбора URLs по UserID
func (ms *StorageFile) URLsByUserID(ctx context.Context, userid string) (userURLs map[string]string, err error) {
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

func (ms *StorageFile) LoadFromFileToStorage() {
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
func (ms *StorageFile) PutBatchToStorage(ctx context.Context, dc settings.DecodeBatchJSON, userid string) (dcCorr settings.DecodeBatchJSON, err error) {
	//userid := ctx.Value(settings.CtxKeyUserID).(string)
	for _, v := range dc {
		// записываем в хранилице userid, id, URL
		ms.IDURL[v.ShortURL] = userid
		ms.UserID[v.ShortURL] = v.OriginalURL
		ms.DelURL[v.ShortURL] = false

	}
	return dc, err
}
