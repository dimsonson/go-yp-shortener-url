package handlers

import (
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Error(w, "userId is empty", http.StatusBadRequest)
		return
	}
	ur := (chi.URLParam(r, "id"))
	log.Println("chi.URLParam: " + ur)
	// проверяем наличие ключа и получем длинную ссылку
	value, ok := storage.DB["/" + chi.URLParam(r, "id")]
	if !ok {
		http.Error(w, "short URL not found", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, value, http.StatusTemporaryRedirect)
}
