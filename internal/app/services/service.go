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
	PutToStorage(userid string, key string, value string) (err error)
	GetFromStorage(key string) (value string, err error)
	LenStorage() (lenn int)
	URLsByUserID(userid string) (userURLs map[string]string, err error)
	LoadFromFileToStorage() 
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
func (sr *Services) ServiceCreateShortURL(url string, userCookie string) (key string, userToken string) {
	// присваиваем значение ключа
	key, err := RandSeq(settings.KeyLeght)
	if err != nil {
		log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
	}
	userid := "testuser"
	// добавляем уникальный префикс к ключу
	key = fmt.Sprintf("%d%s", sr.storage.LenStorage(), key)
	// создаем пару ключ-значение в базе
	sr.storage.PutToStorage(userid, key, url)
	return key, userToken
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

// метод возврат всех URLs по userid
func (sr *Services) ServiceGetUserShortURLs(userToken string) (UserURLsMap map[string]string, err error) {

	userid := userToken
	// используем метод хранилища для получения map URLs по userid
	userURLsMap, err := sr.storage.URLsByUserID(userid)
	fmt.Println("userURLsMap:", userURLsMap)
	if err != nil {
		log.Println(err)
		return map[string]string{"": ""}, err //userURLsMap, err
	}

	return userURLsMap, err
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
