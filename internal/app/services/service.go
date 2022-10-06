package services

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

// интерфейс методов хранилища
type StorageProvider interface {
	PutToStorage(key string, value string) (err error)
	GetFromStorage(key string) (value string, err error)
	LenStorage() (lenn int)
}

// структура конструктора бизнес логики
type Services struct {
	storage StorageProvider
}

// конструктор бизнес логики
func NewService(s StorageProvider) *Services {
	return &Services{
		s,
	}
}

// метод создание пары id : URL
func (sr *Services) ServiceCreateShortURL(url string) (key string) {
	// присваиваем значение ключа
	key, err := RandSeq(settings.KeyLeght)
	if err != nil {
		log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
	}
	// добавляем уникальный префикс к ключу
	key = fmt.Sprintf("%d%s", sr.storage.LenStorage(), key)
	// создаем пару ключ-значение в базе
	sr.storage.PutToStorage(key, url)
	return key
}

// метод возврат URL по id
func (sr *Services) ServiceGetShortURL(id string) (value string, err error) {
	// используем метод хранилища
	value, err = sr.storage.GetFromStorage(id)
	if err != nil {
		log.Println("id not found:", err)
	}
	return value, err
}

// функция генерации случайной последовательности знаков
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
