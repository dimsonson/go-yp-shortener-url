package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type StorageFs struct {
	IdURL map[string]string `json:"idurl,omitempty"`
}

var d = StorageFs{
	IdURL: make(map[string]string),
}

var path = os.Getenv("FILE_STORAGE_PATH")
var fileName = path + "/" + "keyvalue.json"

func (ms *StorageFs) PutStorage(key string, value string) (err error) {
	if _, ok := d.IdURL[key]; ok {
		return fmt.Errorf("key is already in database")
	}
	d.IdURL[key] = string(value)

	// запись в JSON
	sfile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777) //|os.O_APPEND
	if err != nil {
		log.Println("storage file opening/creating error: ", err)
		return err
	}
	defer sfile.Close()

	js, err := json.Marshal(&d.IdURL)
	if err != nil {
		log.Println("JSON marshalling from struct error: ", err)
		return err
	}

	js = append(js, '\n')
	sfile.Write(js)
	return nil
}

func NewFsStorage(s map[string]string) *StorageFs {
	// загрузка базы из JSON
	_, pathOk := os.Stat(path)
	fileInfo, _ := os.Stat(fileName)

	if os.IsNotExist(pathOk) {
		os.Mkdir(path, 0777)
		log.Printf("folder %s created\n", path)
	}

	sfile, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal("file creating error: ", err)
	}

	if fileInfo.Size() != 0 {
		buf := bufio.NewReader(sfile)
		js, err := buf.ReadBytes('\n')
		if err != nil {
			log.Println("file storage reading error:", err)
		}

		err = json.Unmarshal(js, &d.IdURL)
		if err != nil {
			log.Println("JSON unmarshalling to struct error:", err)
		}
	}
	return &StorageFs{
		IdURL: s,
	}
}

func (ms *StorageFs) GetStorage(key string) (value string, err error) {
	value, ok := d.IdURL[key]
	if !ok {
		return "", fmt.Errorf("key %v not found", key)
	}
	return value, nil
}

func (ms *StorageFs) LenStorage() (lenn int) {
	lenn = len(d.IdURL)
	return lenn
}
