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
	IDURL map[string]string `json:"idurl,omitempty"`
	pathName string
}

var d = StorageFs{
	IDURL: make(map[string]string),
	pathName: "",
}

//var pathName = os.Getenv("FILE_STORAGE_PATH")

func (ms *StorageFs) PutStorage(key string, value string) (err error) {
	if _, ok := d.IDURL[key]; ok {
		return fmt.Errorf("key is already in database")
	}
	d.IDURL[key] = string(value)

	// запись в JSON
	sfile, err := os.OpenFile(d.pathName, os.O_WRONLY|os.O_CREATE, 0777) //|os.O_APPEND
	if err != nil {
		log.Println("storage file opening/creating error: ", err)
		return err
	}
	defer sfile.Close()

	js, err := json.Marshal(&d.IDURL)
	if err != nil {
		log.Println("JSON marshalling from struct error: ", err)
		return err
	}

	js = append(js, '\n')
	sfile.Write(js)
	return nil
}

func NewFsStorage(s map[string]string, p string) *StorageFs {
	// загрузка базы из JSON
	_, pathOk := os.Stat(filepath.Dir(d.pathName))

	if os.IsNotExist(pathOk) {
		os.MkdirAll(filepath.Dir(p), 0777)
		log.Printf("folder %s created\n", filepath.Dir(d.pathName))
	}

	sfile, err := os.OpenFile(p, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal("file creating error: ", err)
	}
	fileInfo, _ := os.Stat(p)
	if fileInfo.Size() != 0 {
		buf := bufio.NewReader(sfile)
		js, err := buf.ReadBytes('\n')
		if err != nil {
			log.Println("file storage reading error:", err)
		}

		err = json.Unmarshal(js, &d.IDURL)
		if err != nil {
			log.Println("JSON unmarshalling to struct error:", err)
		}
	}
	
	return &StorageFs{
		IDURL: s,
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
	lenn = len(d.IDURL)
	return lenn
}
