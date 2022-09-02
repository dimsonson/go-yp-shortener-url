package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/randomsuff"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// читаем Body
	b, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("b", string(b))
	//создаем ключ
	var key string //:= randSeq(5)
	// присваиваем значение ключа и проверяем уникальность ключа
	for {
		tmpKey := randomsuff.RandSeq(settings.KeyLeght)
		if _, ok := storage.Db[tmpKey]; !ok {
			key = tmpKey
			break
		}
	}
	fmt.Println("key", key)
	//создаем пару ключ-значение
	storage.Db[key] = string(b)
	fmt.Println("storage.Db[key]", storage.Db[key])
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	//устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)
	// пишем тело ответа
	fmt.Println("r.Host", r.Host)
	w.Write([]byte("http://" + r.Host + key))
}
