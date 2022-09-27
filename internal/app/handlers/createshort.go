package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/randomsuff"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
)

func CreateShortURL(w http.ResponseWriter, r *http.Request) {
	// читаем Body
	b, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = url.ParseRequestURI(string(b))
	if err != nil {
		http.Error(w, "invalid URL received to make short one", http.StatusBadRequest)
		return
	}
	// присваиваем значение не уникального ключа
	tmpKey, err := randomsuff.RandSeq(settings.KeyLeght)
	if err != nil {
		log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
	}
	unique := len(storage.DB) // уникальная часть
	// присваиваем значение уникального ключа
	key := fmt.Sprintf("%v%v", unique, tmpKey)
	//создаем пару ключ-значение
	storage.DB[key] = string(b)
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	//устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)
	// пишем тело ответа
	w.Write([]byte("http://" + r.Host + "/" + key))
}
