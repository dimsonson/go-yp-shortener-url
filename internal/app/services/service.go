package services

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

type Storages interface {
	PutStorage(key string, value string) (err error)
	GetStorage(key string) (value string, err error)
}

type Services struct {
	Storage Storages
}

func NewService(s Storages) *Services {
	return &Services{
		s,
	}
}

func (sr *Services) ServiceCreateShortURL(url string) (key string) {

	// создать метод в storage для len

	// присваиваем значение ключа и проверяем уникальность ключа
	key, err := RandSeq(settings.KeyLeght)
	if err != nil {
		log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
	}

	//создаем пару ключ-значение
	sr.Storage.PutStorage(key, url)
	return key
}

func (sr *Services) ServiceGetShortURL(id string) (value string, err error) {
	value, err = sr.Storage.GetStorage(id)
	if err != nil {
		err = fmt.Errorf("id not found")

	}
	return
}

func RandSeq(n int) (string, error) {
	if n < 1 {
		err := fmt.Errorf("wromg argument: number %v less than 1\n ", n)
		return "", err
	}
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b), nil
}
