package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
)

// интерфейс методов бизнес логики
type Services interface {
	ServiceCreateShortURL(url string, userTokenIn string) (key string, userTokenOut string)
	ServiceGetShortURL(id string) (value string, err error)
	ServiceGetUserShortURLs(userToken string) (UserURLsMap map[string]string, err error)
	ServiceStorageOkPing() bool
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

// структура декодирования JSON для POST запроса
type DecodeJSON struct {
	URL string `json:"url,omitempty"`
}

// структура кодирования JSON для POST запроса
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

	// читаем куку с userid
	var userToken string
	userCookie, err := r.Cookie("token")
	if err != nil {
		log.Println("Request does not consist token cookie - err:", err)
	} else {
		userToken = userCookie.Value
	}
	fmt.Println("userCookie.Value:", userToken)

	// создаем ключ и userid token
	key, userTokenNew := hn.handler.ServiceCreateShortURL(b, userToken)

	// создаем и записываем куку в ответ если ее нет в запросе или она создана сервисом
	if err != nil || userTokenNew != userToken {
		cookie := &http.Cookie{
			Name:   "token",
			Value:  userTokenNew,
			MaxAge: 300,
		}
		fmt.Println("cookie:  ", cookie)
		// установим куку в ответ
		http.SetCookie(w, cookie)
	}

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
		http.Error(w, "invalid JSON structure received", http.StatusBadRequest)
	}
	// валидация URL
	if !govalidator.IsURL(dc.URL) {
		http.Error(w, "invalid URL received to make short one", http.StatusBadRequest)
		return
	}
	// читаем куку с userid
	var userToken string
	userCookie, err := r.Cookie("token")
	if err != nil {
		log.Println("Request does not consist token cookie - err:", err)
	} else {
		userToken = userCookie.Value
	}
	fmt.Println("userCookie.Value:", userToken)

	// создаем ключ и userid token
	key, userTokenNew := hn.handler.ServiceCreateShortURL(dc.URL, userToken)

	// создаем и записываем куку в ответ если ее нет в запросе или она создана сервисом
	if err != nil || userTokenNew != userToken {
		cookie := &http.Cookie{
			Name:   "token",
			Value:  userTokenNew,
			MaxAge: 300,
		}
		fmt.Println("cookie:  ", cookie)
		// установим куку в ответ
		http.SetCookie(w, cookie)
	}
	// сериализация тела запроса
	ec := EncodeJSON{}
	ec.Result = hn.base + "/" + key
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "application/json; charset=utf-8")
	//устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)
	// пишем тело ответа
	json.NewEncoder(w).Encode(ec)
}

// обработка GET запроса /api/user/urls c возвратом пользователю всех когда-либо сокращённых им URL
func (hn Handler) HandlerGetUserURLs(w http.ResponseWriter, r *http.Request) {

	// читаем куку с userid
	userCookie, err := r.Cookie("token")
	if err != nil {
		log.Println("request does not consist token cookie - err:", err)
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	fmt.Println("userCookie.Value:", userCookie.Value)
	// получаем map всех URLs по usertoken
	userURLsMap, err := hn.handler.ServiceGetUserShortURLs(userCookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	// создаем и заполняем слайс структур
	UserURLs := []UserURL{}
	for k, v := range userURLsMap {
		k = hn.base + "/" + k
		UserURLs = append(UserURLs, UserURL{k, v})
	}
	fmt.Println("UserURLs:", UserURLs)
	// сериализация тела запроса
	w.Header().Set("content-type", "application/json; charset=utf-8")
	//устанавливаем статус-код 201
	w.WriteHeader(http.StatusOK)
	// сериализуем и пишем тело ответа
	json.NewEncoder(w).Encode(UserURLs)
	/*
		 	if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	*/
}

// структура для создания среза surl:url и дельнейшего ecode
type UserURL struct {
	ShortURL    string `json:"short_url,omitempty"`
	OriginalURL string `json:"original_url,omitempty"`
}

func (hn Handler) HandlerSQLping(w http.ResponseWriter, r *http.Request) {
	var result []byte
	if !hn.handler.ServiceStorageOkPing() {
		w.WriteHeader(http.StatusInternalServerError)
		result = []byte("DB ping NOT OK")

	} else {
		w.WriteHeader(http.StatusOK)
		result = []byte("DB ping OK")
	}
	w.Write(result)
}
