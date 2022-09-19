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
	LenStorage() (lenn int)
}

type Services struct {
	storage Storages
}

func NewService(s Storages) *Services {
	return &Services{
		s,
	}
}
// создание пары id : URL
func (sr *Services) ServiceCreateShortURL(url string) (key string) {
	// присваиваем значение ключа
	key, err := RandSeq(settings.KeyLeght)
	if err != nil {
		log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
	}
	// добавляем уникальный префикс к ключу
	key = fmt.Sprintf("%d%s", sr.storage.LenStorage(), key)
	// создаем пару ключ-значение в базе
	sr.storage.PutStorage(key, url)

	return key
}
// возврат URL по id
func (sr *Services) ServiceGetShortURL(id string) (value string, err error) {
	value, err = sr.storage.GetStorage(id)
	if err != nil {
		err = fmt.Errorf("id not found")
	}
	return
}
// генерация случайной последовательности знаков
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
