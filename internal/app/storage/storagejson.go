package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// структура хранилища
type StorageJSON struct {
	IDURL    map[string]string `json:"idurl,omitempty"`
	pathName string
}

// метод записи id:url в хранилище
func (ms *StorageJSON) PutToStorage(key string, value string) (err error) {
	if value, ok := ms.IDURL[key]; ok {
		return fmt.Errorf("key %s is already in database", value)
	}
	ms.IDURL[key] = value
	// открываем файл
	sfile, err := os.OpenFile(ms.pathName, os.O_WRONLY, 0777)
	if err != nil {
		log.Println("storage file opening error: ", err)
		return err
	}
	defer sfile.Close()
	// кодирование в JSON
	js, err := json.Marshal(&ms.IDURL)
	if err != nil {
		log.Println("JSON marshalling from struct error: ", err)
		return err
	}
	// запись в файл
	sfile.Write(js)
	return nil
}

// конструктор нового хранилища JSON
func NewJSONStorage(s map[string]string, p string) *StorageJSON {
	// загрузка базы из JSON
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
		err = json.Unmarshal(b, &s)
		if err != nil {
			log.Println("JSON unmarshalling to struct error:", err)
		}
	}
	return &StorageJSON{
		IDURL:    s,
		pathName: p,
	}
}

// метод получения записи из хранилища
func (ms *StorageJSON) GetFromStorage(key string) (value string, err error) {
	value, ok := ms.IDURL[key]
	if !ok {
		return "", fmt.Errorf("key %v not found", key)
	}
	return value, nil
}

// метод определения длинны хранилища
func (ms *StorageJSON) LenStorage() (lenn int) {
	lenn = len(ms.IDURL)
	return lenn
}
