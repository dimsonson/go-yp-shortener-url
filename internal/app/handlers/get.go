package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/go-chi/chi/v5"
)

// GetServiceProvider интерфейс методов бизнес логики для слоя Get.
type GetServiceProvider interface {
	Get(ctx context.Context, id string) (value string, del bool, err error)
	GetBatch(ctx context.Context, userid string) (userURLsMap map[string]string, err error)
}

// GetHandler структура для конструктура обработчика.
type GetHandler struct {
	service GetServiceProvider
	base    string
}

// NewGetHandler конструктор обработчика.
func NewGetHandler(s GetServiceProvider, base string) *GetHandler {
	return &GetHandler{
		s,
		base,
	}
}

// GetHandler метод обработки GET запроса c id и редирект по полному URL.
func (hn GetHandler) Get(w http.ResponseWriter, r *http.Request) {
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
	value, del, err := hn.service.Get(ctx, key)
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
		_, err := w.Write([]byte(value))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
}

// GetBatch метод обработки GET запроса /api/user/urls c возвратом пользователю всех когда-либо сокращённых им URL.
func (hn GetHandler) GetBatch(w http.ResponseWriter, r *http.Request) {
	// получаем значение iserid из контекста запроса
	userid := r.Context().Value(settings.CtxKeyUserID).(string)
	// наследуем контекcт запроса r *http.Request, оснащая его Timeout
	ctx, cancel := context.WithTimeout(r.Context(), settings.StorageTimeout)
	// не забываем освободить ресурс
	defer cancel()
	// получаем map всех URLs по usertoken
	userURLsMap, err := hn.service.GetBatch(ctx, userid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	// создаем и заполняем слайс структур
	var UserURLs []models.UserURL
	for k, v := range userURLsMap {
		k = hn.base + "/" + k
		UserURLs = append(UserURLs, models.UserURL{ShortURL: k, OriginalURL: v})
	}
	// сериализация тела запроса
	w.Header().Set("content-type", "application/json; charset=utf-8")
	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusOK)
	// сериализуем и пишем тело ответа
	err = json.NewEncoder(w).Encode(UserURLs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
