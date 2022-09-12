package services

import (
	"log"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/randomsuff"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
)

func ShortService(url string, id string) (value string, key string) {
	//var key string //:= randSeq(5)
	// присваиваем значение ключа и проверяем уникальность ключа
	for {
		tmpKey, err := randomsuff.RandSeq(settings.KeyLeght)
		if err != nil {
			log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
		}
		
		if _, ok := storage.DB[tmpKey]; !ok {
			key = tmpKey
			break
		}
	}
	//создаем пару ключ-значение
	storage.DB[key] = string(handlers.B)

	return 

}
