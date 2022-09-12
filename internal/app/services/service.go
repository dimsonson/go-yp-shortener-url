package services

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
)

func ServiseCreateShortURL(url string) (key string) {
	//var key string //:= randSeq(5)
	// присваиваем значение ключа и проверяем уникальность ключа
	for {
		tmpKey, err := RandSeq(settings.KeyLeght)
		if err != nil {
			log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
		}

		if _, ok := storage.DB[tmpKey]; !ok {
			key = tmpKey
			break
		}
	}
	//создаем пару ключ-значение
	//storage.DB[key] = string(handlers.B)
	storage.PutMapStorage(key, url)

	return key

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
