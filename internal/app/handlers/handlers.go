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

// структура декодирования JSON
type DecodeJSON struct {
	URL string `json:"url,omitempty"`
}

// структура кодирования JSON
type EncodeJSON struct {
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
	// проверяем наличие id
	if r.URL.Path == "/" {
		http.Error(w, "userId is empty", http.StatusBadRequest)
		return
	}
	// пролучаем id из URL через chi
	id := chi.URLParam(r, "id")
	// получаем ссылку по id
	value, err := hn.handler.ServiceGetShortURL(id)
	if err != nil {
		http.Error(w, "short URL not found", http.StatusBadRequest)
	}
	// перенаправление по ссылке
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
	dc := DecodeJSON{}
	if err := json.Unmarshal(B, &dc); err != nil {
		http.Error(w, "invalid URL received to make short one", http.StatusBadRequest)
	}
	// валидация URL
	_, err = url.ParseRequestURI(dc.URL)
	if err != nil {
		http.Error(w, "invalid URL received to make short one", http.StatusBadRequest)
		return
	}
	//создаем ключ
	key := hn.handler.ServiceCreateShortURL(dc.URL)
	// сериализация тела запроса
	ec := EncodeJSON{}
	ec.Result = "http://" + r.Host + "/" + key
	jsn, err := json.Marshal(ec)
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
