package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/go-chi/chi/v5"
)

// интерфейс методов бизнес логики
type Services interface {
	ServiceCreateShortURL(ctx context.Context, url string) (key string, err error)
	ServiceGetShortURL(ctx context.Context, id string) (value string, err error)
	ServiceGetUserShortURLs(ctx context.Context) (UserURLsMap map[string]string, err error)
	ServiceStorageOkPing(ctx context.Context) (bool, error)
	ServiceCreateBatchShortURLs(ctx context.Context, dc settings.DecodeBatchJSON) (ec []settings.EncodeBatch, err error)
}

// структура для конструктура обработчика
type Handler struct {
	service Services
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
	// наследуем контекcт запроса r *http.Request, оснащая его Timeout
	ctx, cancel := context.WithTimeout(r.Context(), settings.StorageTimeout)
	// не забываем освободить ресурс
	defer cancel()
	// создаем ключ и userid token
	key, err := hn.service.ServiceCreateShortURL(ctx, b)
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	//устанавливаем статус-код 201, 500 или 409
	switch {
	case err != nil && strings.Contains(err.Error(), "23505"):
		w.WriteHeader(http.StatusConflict)
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusCreated)
	}
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
	// наследуем контекcт запроса r *http.Request, оснащая его Timeout
	ctx, cancel := context.WithTimeout(r.Context(), settings.StorageTimeout)
	// не забываем освободить ресурс
	defer cancel()
	// получаем ссылку по id
	value, err := hn.service.ServiceGetShortURL(ctx, id)
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
	// наследуем контекcт запроса r *http.Request, оснащая его Timeout
	ctx, cancel := context.WithTimeout(r.Context(), settings.StorageTimeout)
	// не забываем освободить ресурс
	defer cancel()
	// создаем ключ, userid token, ошибку создания в случае налияи URL в базе
	key, err := hn.service.ServiceCreateShortURL(ctx, dc.URL)
	// сериализация тела запроса
	ec := EncodeJSON{}
	ec.Result = hn.base + "/" + key
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "application/json; charset=utf-8")
	//устанавливаем статус-код 201, 500 или 409
	switch {
	case err != nil && strings.Contains(err.Error(), "23505"):
		w.WriteHeader(http.StatusConflict)
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusCreated)
	}
	// пишем тело ответа
	json.NewEncoder(w).Encode(ec)
}

// обработка GET запроса /api/user/urls c возвратом пользователю всех когда-либо сокращённых им URL
func (hn Handler) HandlerGetUserURLs(w http.ResponseWriter, r *http.Request) {
	// наследуем контекcт запроса r *http.Request, оснащая его Timeout
	ctx, cancel := context.WithTimeout(r.Context(), settings.StorageTimeout)
	// не забываем освободить ресурс
	defer cancel()
	// получаем map всех URLs по usertoken
	userURLsMap, err := hn.service.ServiceGetUserShortURLs(ctx)
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
	// сериализация тела запроса
	w.Header().Set("content-type", "application/json; charset=utf-8")
	//устанавливаем статус-код 201
	w.WriteHeader(http.StatusOK)
	// сериализуем и пишем тело ответа
	json.NewEncoder(w).Encode(UserURLs)
}

// структура для создания среза surl:url и дельнейшего encode
type UserURL struct {
	ShortURL    string `json:"short_url,omitempty"`
	OriginalURL string `json:"original_url,omitempty"`
}

// проверка доступности базы SQL
func (hn Handler) HandlerSQLping(w http.ResponseWriter, r *http.Request) {
	// наследуем контекcт запроса r *http.Request, оснащая его Timeout
	ctx, cancel := context.WithTimeout(r.Context(), settings.StorageTimeout)
	// не забываем освободить ресурс
	defer cancel()
	// создаем переменную для визуального возврата пользователю в теле отвта
	var result []byte
	ok, err := hn.service.ServiceStorageOkPing(ctx)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		result = []byte("DB ping NOT OK")
		log.Println(err)
	} else {
		w.WriteHeader(http.StatusOK)
		result = []byte("DB ping OK")
	}
	w.Write(result)
}

// обработка POST запроса с JSON batch в теле и возврат Batch JSON c короткими URL
// посмотреть в будущем вариант записи через отдельный метод хранилища с стейтментами
func (hn Handler) HandlerCreateBatchJSON(w http.ResponseWriter, r *http.Request) {
	// наследуем контекcт запроса r *http.Request, оснащая его Timeout
	ctx, cancel := context.WithTimeout(r.Context(), settings.StorageTimeout)
	// не забываем освободить ресурс
	defer cancel()
	// десериализация тела запроса
	dc := settings.DecodeBatchJSON{} //DecodeBatchJSON{}
	err := json.NewDecoder(r.Body).Decode(&dc)
	if err != nil {
		log.Printf("Unmarshal error: %s", err)
		http.Error(w, "invalid JSON structure received", http.StatusBadRequest)
	}

	/* 	reqMap := make(map[string]string)
	   	for _, v := range dc {
	   		if !govalidator.IsURL(v.OriginalURL) {
	   			http.Error(w, "invalid URL received to make short one", http.StatusBadRequest)
	   			return
	   		}
	   		reqMap[v.CorrelationID] = v.OriginalURL
	   	} */
	// запрос на получение пар кароткий - длинный URL
	ec, err := hn.service.ServiceCreateBatchShortURLs(ctx, dc)
	if err != nil {
		log.Println(err) // подумать над обработкой
	}
	// сериализация тела ответа
	//	ec := []EncodeBatchJSON{}

	// итерируем по полученнму map, пишем в исходящий слайс стркутур
	/* 	for k, v := range respMap {
		// добавляем структуру в слайс
		ec = append(ec, EncodeBatchJSON{
			CorrelationID: k,
			ShortURL:      v,
		})
	} */
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "application/json; charset=utf-8")
	//устанавливаем статус-код 201, 500 или 409
	switch {
	case err != nil && strings.Contains(err.Error(), "23505"):
		w.WriteHeader(http.StatusConflict)
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusCreated)
	}
	// пишем тело ответа
	json.NewEncoder(w).Encode(ec)
}

// слайс структур декодирования JSON из POST запроса
/* type DecodeBatchJSON []struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
}

// структура кодирования JSON для POST Batch ответа
type EncodeBatchJSON struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
} */

type DecodeBatchMap map[string]string
