package handlers

import (
	"fmt"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r", r)
	if r.URL.Path == "/" {
		http.Error(w, "userId is empty", http.StatusBadRequest)
		return
	}
	fmt.Println("r.URL.Path", r.URL.Path)
	// проверяем наличие ключа и получем длинную ссылку
	value, ok := storage.Db[r.URL.Path]
	if !ok {
		http.Error(w, "short URL not found", http.StatusBadRequest)
		return
	}
	fmt.Println("value", value)
	http.Redirect(w, r, value, http.StatusTemporaryRedirect)
}
