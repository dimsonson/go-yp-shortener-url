package handlers

import (
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Error(w, "userId is empty", http.StatusBadRequest)
		return
	}
	// проверяем наличие ключа и получем длинную ссылку
	value, ok := storage.Db[r.URL.Path]
	if !ok {
		http.Error(w, "short URL not found", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, value, http.StatusTemporaryRedirect)
}
