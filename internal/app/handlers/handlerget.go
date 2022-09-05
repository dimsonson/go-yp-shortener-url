package handlers

import (
	"fmt"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("метод GET")
	if r.URL.Path == "/" {
		http.Error(w, "userId is empty", http.StatusBadRequest)
		return
	}
	// проверяем наличие ключа и получем длинную ссылку
	value, ok := storage.DB[r.URL.Path]
	if !ok {
		http.Error(w, "short URL not found", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, value, http.StatusTemporaryRedirect)
}
