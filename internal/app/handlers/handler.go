package handlers

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/randomsuff"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/services"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-delve/delve/service"
)

func HandlerCreateShortURL(w http.ResponseWriter, r *http.Request) {
	// читаем Body
	B, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = url.ParseRequestURI(string(B))
	if err != nil {
		http.Error(w, "invalid URL received to make short one", http.StatusBadRequest)
		return
	}
	//создаем ключ
    key := services.ServiseCreateShortURL(string(B))
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	//устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)
	// пишем тело ответа
	w.Write([]byte("http://" + r.Host + "/" + key))
}

func HandlerGetShortURL(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Error(w, "userId is empty", http.StatusBadRequest)
		return
	}
	// проверяем наличие ключа и получем длинную ссылку
	value, ok := storage.DB[chi.URLParam(r, "id")]
	if !ok {
		http.Error(w, "short URL not found", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, value, http.StatusTemporaryRedirect)
}

func IncorrectRequests(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "request incorect", http.StatusBadRequest)
}
