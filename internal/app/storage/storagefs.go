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

type read struct {
	file    *os.File
	decoder *json.Decoder
}

var fileName = "keyvalue.json"

func (ms *StorageFs) PutStorage(key string, value string) (err error) {

	if _, ok := d.IdURL[key]; ok {
		return fmt.Errorf("key is already in database")
	}
	d.IdURL[key] = string(value)

	// запись в JSON
	sfile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777) //|os.O_APPEND
	if err != nil {
		log.Fatal("storage file opening/creating error")
	}
	defer sfile.Close()

	js, err := json.Marshal(&d.IdURL)
	if err != nil {
		log.Println("JSON from struct error: ", err)
	}

	js = append(js, '\n')
	sfile.Write(js)

	fmt.Println(d)
	return nil
}

func NewFsStorage(s map[string]string) *StorageFs {
	// загрузка базы из JSON

	sfile, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Println("file creating error", err)
	}
	fmt.Println("sfile: ", sfile)
    buf :=bufio.NewReader(sfile)
	js, _ := buf.ReadBytes('\n')

	//var js []byte
	//_, err = sfile.Read(js)
	if err != nil {
		log.Fatal("storage file opening/creating error")
	}
	fmt.Println("JSON для деодинга: ", js)
	json.Unmarshal(js, &d.IdURL)

	fmt.Println("база после чтения из файла: ", d.IdURL)

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

func (p *read) ReadFile(event *StorageFs) error {
	return p.decoder.Decode(&event)
}
