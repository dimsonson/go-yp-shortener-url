package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
)

type Services interface {
	ServiceCreateShortURL(url string) (key string)
	ServiceGetShortURL(id string) (value string, err error)
}

type Handler struct {
	handler Services
	Base    string
}

var hn = Handler{
	handler: nil,
	Base:    "",
}

// значение переменной BASE_URL по умолчанию
var defBase string = "http://localhost:8080"

func NewHandler(s Services, base string) *Handler {
	// проверка переменной окуржения и присвоение значения по умолчанию, если не установлено
	var ok bool
	if !govalidator.IsURL(base) {
		base, ok = os.LookupEnv("BASE_URL")
		if !ok || !govalidator.IsURL(base) {
			base = defBase
			log.Println("enviroment variable BASE_URL set to default value:", defBase)
		}
	}
	fmt.Println("base", base)
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
	// создаем ключ
	key := hn.handler.ServiceCreateShortURL(string(B))

	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	//устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)
	// пишем тело ответа
	fmt.Println(hn.Base)
	w.Write([]byte(hn.Base + "/" + key))
}

// обработка GET запроса c id и редирект по полному URL
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

// обработка всех остальных запросов и возврат кода 400
func (hn Handler) IncorrectRequests(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "request incorrect", http.StatusBadRequest)
}

// обработка POST запроса с JSON URL в теле и возврат короткого URL JSON в теле
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
	ec.Result = hn.Base + "/" + key
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
