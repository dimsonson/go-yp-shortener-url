package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

type Services interface {
	ServiceCreateShortURL(url string) (key string)
	ServiceGetShortURL(id string) (value string, err error)
}

type Handler struct {
	handler Services
}

func NewHandler(s Services) *Handler {
	return &Handler{
		s,
	}
}

type PostJSON struct {
	URL    string `json:"url,omitempty"`
	Result string `json:"result,omitempty"`
}

func (hn Handler) HandlerCreateShortURL(w http.ResponseWriter, r *http.Request) {
	// читаем Body
	B, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// валидация URL
	_, err = url.ParseRequestURI(string(B))
	if err != nil {
		http.Error(w, "invalid URL received to make short one", http.StatusBadRequest)
		return
	}
	//создаем ключ
	key := hn.handler.ServiceCreateShortURL(string(B))
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	//устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)
	// пишем тело ответа
	w.Write([]byte("http://" + r.Host + "/" + key))
}

func (hn Handler) HandlerGetShortURL(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Error(w, "userId is empty", http.StatusBadRequest)
		return
	}
	// проверяем наличие ключа и получем длинную ссылку
	id := chi.URLParam(r, "id")

	value, err := hn.handler.ServiceGetShortURL(id)
	if err != nil {
		http.Error(w, "short URL not found", http.StatusBadRequest)
	}

	http.Redirect(w, r, value, http.StatusTemporaryRedirect)
}

func (hn Handler) IncorrectRequests(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "request incorrect", http.StatusBadRequest)
}

func (hn Handler) HandlerCreateShortJSON(w http.ResponseWriter, r *http.Request) {
	// читаем Body
	B, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// десериализация тела запроса
	b := PostJSON{}
	if err := json.Unmarshal(B, &b); err != nil {
		http.Error(w, "invalid URL received to make short one", http.StatusBadRequest)
	}

	// валидация URL
	_, err = url.ParseRequestURI(b.URL)
	if err != nil {
		http.Error(w, "invalid URL received to make short one", http.StatusBadRequest)
		return
	}
	//создаем ключ
	key := hn.handler.ServiceCreateShortURL(b.URL)
	b.Result = "http://" + r.Host + "/" + key
	b.URL = ""
	jsn, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "application/json; charset=utf-8")
	//устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)
	// пишем тело ответа
	w.Write(jsn)
}
