// Package handlers пакет обработчиков http запросов.
package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

// DeleteServiceProvider интерфейс методов бизнес логики для слоя Delete.
type DeleteServiceProvider interface {
	Delete(shURLs []([2]string))
}

// DeleteHandler структура для конструктура обработчика.
type DeleteHandler struct {
	service DeleteServiceProvider
	base    string
}

// NewDeleteHandler конструктор обработчика.
func NewDeleteHandler(s DeleteServiceProvider, base string) *DeleteHandler {
	return &DeleteHandler{
		s,
		base,
	}
}

// Delete метод обработки DELETE запроса с слайсом short_url в теле.
func (hn DeleteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// получаем значение iserid из контекста запроса
	userid := r.Context().Value(settings.CtxKeyUserID).(string)
	// десериализация тела запроса
	var d []string
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil && err != io.EOF {
		log.Printf("unmarshal error: %s", err)
		http.Error(w, "invalid slice of short_urls received", http.StatusBadRequest)
	}
	// создание и наполнение слайса массивов для передачи в fanout-fanin
	var shURLs []([2]string)
	for _, v := range d {
		shURLs = append(shURLs, [2]string{v, userid})
	}
	// запуск сервиса внесения записей о удалении
	go hn.service.Delete(shURLs)
	// устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "application/json; charset=utf-8")
	// записываем статус-код 202
	w.WriteHeader(http.StatusAccepted)
}
