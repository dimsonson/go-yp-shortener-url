package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

// PingServiceProvider интерфейс методов бизнес логики слоя Ping.
type PingServiceProvider interface {
	Ping(ctx context.Context) (bool, error)
}

// PingHandler структура для конструктура обработчика.
type PingHandler struct {
	service PingServiceProvider
	base    string
}

// NewPingHandler конструктор обработчика.
func NewPingHandler(s PingServiceProvider, base string) *PingHandler {
	return &PingHandler{
		s,
		base,
	}
}

// Ping метод проверки доступности базы SQL.
func (hn PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	// наследуем контекcт запроса r *http.Request, оснащая его Timeout
	ctx, cancel := context.WithTimeout(r.Context(), settings.StorageTimeout)
	// не забываем освободить ресурс
	defer cancel()
	// создаем переменную для визуального возврата пользователю в теле отвта
	var result []byte
	ok, err := hn.service.Ping(ctx)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		result = []byte("DB ping NOT OK")
		log.Println(err)
	} else {
		w.WriteHeader(http.StatusOK)
		result = []byte("DB ping OK")
	}
	_, err = w.Write(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
