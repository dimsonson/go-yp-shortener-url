package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/randomsuff"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// проверяем пришел ли POST запрос
	if r.Method != "POST" {
		http.Error(w, "Error: method is not POST", http.StatusBadRequest)
		return
	}
	// читаем Body
	b, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//создаем ключ
	var key string //:= randSeq(5)
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
	storage.DB[key] = string(b)
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	//устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)
	// пишем тело ответа
	w.Write([]byte("http://" + r.Host + key))
}
