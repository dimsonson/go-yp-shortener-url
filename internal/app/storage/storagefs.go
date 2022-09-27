package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type StorageFs struct {
	IDURL    map[string]string `json:"-"`
	Id       string            `json:"id,omitempty"`
	URL      string            `json:"url,omitempty"`
	pathName string            `json:"-"`
}

var d = StorageFs{
	IDURL: make(map[string]string),
	Id:       "",            
	URL:      "",           
	pathName: "",
}

// без рефакторинга
func (ms *StorageFs) PutStorage(key string, value string) (err error) {
	if _, ok := ms.IDURL[key]; ok {
		return fmt.Errorf("key is already in database")
	}
	ms.IDURL[key] = string(value)
	ms.Id = key
	ms.URL = value

	// запись в JSON
	sfile, err := os.OpenFile(ms.pathName, os.O_WRONLY|os.O_APPEND, 0777) //|os.O_APPEND
	if err != nil {
		log.Println("storage file opening error: ", err)
		return err
	}
	defer sfile.Close()

	js, err := json.Marshal(ms)
	if err != nil {
		log.Println("JSON marshalling from struct error: ", err)
		return err
	}
	js = append(js, '\n')
	sfile.Write(js)
	return nil
}

func NewFsStorage(s map[string]string, p string) *StorageFs {
	// загрузка  базы из JSON
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
		buf := bufio.NewScanner(sfile)
		for buf.Scan() {
			err = json.Unmarshal(buf.Bytes(), &d)
			if err != nil {
				log.Println("JSON unmarshalling to struct error:", err)
			}
			d.IDURL[d.Id] = d.URL
		}
	}
	return &StorageFs{
		IDURL:    s,
		pathName: p,
	}
}

func (ms *StorageFs) GetStorage(key string) (value string, err error) {
	value, ok := d.IDURL[key]
	if !ok {
		return "", fmt.Errorf("key %v not found", key)
	}
	return value, nil
}

func (ms *StorageFs) LenStorage() (lenn int) {
	lenn = len(ms.IDURL)
	return lenn
}
