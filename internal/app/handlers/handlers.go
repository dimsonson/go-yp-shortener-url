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
	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
)

// интерфейс методов бизнес логики
type Services interface {
	ServiceCreateShortURL(ctx context.Context, url string, userid string) (key string, err error)
	ServiceGetShortURL(ctx context.Context, id string) (value string, del bool, err error)
	ServiceGetUserShortURLs(ctx context.Context, userid string) (userURLsMap map[string]string, err error)
	ServiceStorageOkPing(ctx context.Context) (bool, error)
	ServiceCreateBatchShortURLs(ctx context.Context, dc models.BatchRequest, userid string) (ec []models.BatchResponse, err error)
	ServiceDeleteURL(shURLs []([2]string))
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
	// получаем значение iserid из контекста запроса
	userid := r.Context().Value(settings.CtxKeyUserID).(string)
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
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	// создаем ключ и userid token
	key, err := hn.service.ServiceCreateShortURL(ctx, b, userid)
	//устанавливаем статус-код 201, 500 или 409
	switch {
	case err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation):
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
	key := chi.URLParam(r, "id")
	if key == "" {
		http.Error(w, "userId is empty", http.StatusBadRequest)
		return
	}
	// наследуем контекcт запроса r *http.Request, оснащая его Timeout
	ctx, cancel := context.WithTimeout(r.Context(), settings.StorageTimeout)
	// освобождаем ресурс
	defer cancel()
	// устанавливаем заголовок content-type
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	// получаем ссылку по id
	value, del, err := hn.service.ServiceGetShortURL(ctx, key)
	if err != nil {
		http.Error(w, "short URL not found", http.StatusBadRequest)
	}
	if del {
		// сообщаем что ссылка удалена
		http.Error(w, "short URL is deleted", http.StatusGone)
	} else {
		// перенаправление по ссылке
		http.Redirect(w, r, value, http.StatusTemporaryRedirect)
		// пишем тело ответа
		w.Write([]byte(value))
	}
}

// обработка всех остальных запросов и возврат кода 400
func (hn Handler) IncorrectRequests(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "request incorrect", http.StatusBadRequest)
}

// обработка POST запроса с JSON URL в теле и возврат короткого URL JSON в теле
func (hn Handler) HandlerCreateShortJSON(w http.ResponseWriter, r *http.Request) {
	// получаем значение iserid из контекста запроса
	userid := r.Context().Value(settings.CtxKeyUserID).(string)
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
	key, err := hn.service.ServiceCreateShortURL(ctx, dc.URL, userid)
	// сериализация тела запроса
	ec := EncodeJSON{}
	ec.Result = hn.base + "/" + key
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "application/json; charset=utf-8")
	//устанавливаем статус-код 201, 500 или 409
	switch {
	case err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation):
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
	// получаем значение iserid из контекста запроса
	userid := r.Context().Value(settings.CtxKeyUserID).(string)
	// наследуем контекcт запроса r *http.Request, оснащая его Timeout
	ctx, cancel := context.WithTimeout(r.Context(), settings.StorageTimeout)
	// не забываем освободить ресурс
	defer cancel()
	// получаем map всех URLs по usertoken
	userURLsMap, err := hn.service.ServiceGetUserShortURLs(ctx, userid)
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
func (hn Handler) HandlerCreateBatchJSON(w http.ResponseWriter, r *http.Request) {
	// получаем значение iserid из контекста запроса
	userid := r.Context().Value(settings.CtxKeyUserID).(string)
	// наследуем контекcт запроса r *http.Request, оснащая его Timeout
	ctx, cancel := context.WithTimeout(r.Context(), settings.StorageTimeout)
	// не забываем освободить ресурс
	defer cancel()
	// десериализация тела запроса
	dc := models.BatchRequest{} //DecodeBatchJSON{}
	err := json.NewDecoder(r.Body).Decode(&dc)
	if err != nil {
		log.Printf("Unmarshal error: %s", err)
		http.Error(w, "invalid JSON structure received", http.StatusBadRequest)
	}
	// запрос на получение correlation_id  - original_url
	ec, err := hn.service.ServiceCreateBatchShortURLs(ctx, dc, userid)
	if err != nil {
		log.Println(err) // подумать над обработкой
	}
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "application/json; charset=utf-8")
	//устанавливаем статус-код 201, 500 или 409
	switch {
	case err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation):
		w.WriteHeader(http.StatusConflict)
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusCreated)
	}
	// пишем тело ответа
	json.NewEncoder(w).Encode(ec)
}

// обработка DELETE запроса с слайсом short_url в теле
func (hn Handler) HandlerDeleteBatch(w http.ResponseWriter, r *http.Request) {
	// получаем значение iserid из контекста запроса
	userid := r.Context().Value(settings.CtxKeyUserID).(string)
	// десериализация тела запроса
	d := []string{}
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		log.Printf("Unmarshal error: %s", err)
		http.Error(w, "invalid slice of short_urls received", http.StatusBadRequest)
	}
	// создание и наполнение слайса массивов для передачи в fanout-fanin
	var shURLs []([2]string)
	for _, v := range d {
		shURLs = append(shURLs, [2]string{v, userid})
	}
	// запуск сервиса внесения записей о удалении
	go hn.service.ServiceDeleteURL(shURLs)
	// устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "application/json; charset=utf-8")
	// записываем статус-код 202
	w.WriteHeader(http.StatusAccepted)
}
