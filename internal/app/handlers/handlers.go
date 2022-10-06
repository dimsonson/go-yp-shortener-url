package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
)

// интерфейс методов бизнес логики
type Services interface {
	ServiceCreateShortURL(url string) (key string)
	ServiceGetShortURL(id string) (value string, err error)
}

// структура для конструктура обработчика
type Handler struct {
	handler Services
	base    string
}

// конструктор обработчика
func NewHandler(s Services, base string) *Handler {
	return &Handler{
		s,
		base,
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

// обработка POST запроса с text URL в теле и возврат короткого URL в теле
func (hn Handler) HandlerCreateShortURL(w http.ResponseWriter, r *http.Request) {
	// читаем Body
	bs, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b := string(bs)
	// валидация URL
	if !govalidator.IsURL(b) {
		http.Error(w, "invalid URL received to make short one", http.StatusBadRequest)
		return
	}
	// создаем ключ
	key := hn.handler.ServiceCreateShortURL(b)
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	//устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)
	// пишем тело ответа
	w.Write([]byte(hn.base + "/" + key))
}

// обработка GET запроса c id и редирект по полному URL
func (hn Handler) HandlerGetShortURL(w http.ResponseWriter, r *http.Request) {
	// пролучаем id из URL через chi, проверяем наличие
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "userId is empty", http.StatusBadRequest)
		return
	}
	// получаем ссылку по id
	value, err := hn.handler.ServiceGetShortURL(id)
	if err != nil {
		http.Error(w, "short URL not found", http.StatusBadRequest)
	}
	// устанавливаем заголовок content-type
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	// перенаправление по ссылке
	http.Redirect(w, r, value, http.StatusTemporaryRedirect)
	// пишем тело ответа
	w.Write([]byte(value))
}

// обработка всех остальных запросов и возврат кода 400
func (hn Handler) IncorrectRequests(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "request incorrect", http.StatusBadRequest)
}

// обработка POST запроса с JSON URL в теле и возврат короткого URL JSON в теле
func (hn Handler) HandlerCreateShortJSON(w http.ResponseWriter, r *http.Request) {
	// десериализация тела запроса
	dc := DecodeJSON{}
	err := json.NewDecoder(r.Body).Decode(&dc)
	if err != nil {
		log.Printf("Unmarshal error: %s", err)
	}
	// валидация URL
	if !govalidator.IsURL(dc.URL) {
		http.Error(w, "invalid URL received to make short one", http.StatusBadRequest)
		return
	}
	//создаем ключ
	key := hn.handler.ServiceCreateShortURL(dc.URL)
	// сериализация тела запроса
	ec := EncodeJSON{}
	ec.Result = hn.base + "/" + key
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "application/json; charset=utf-8")
	//устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)
	// пишем тело ответа
	json.NewEncoder(w).Encode(ec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
