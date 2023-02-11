package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/asaskevich/govalidator"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/jackc/pgerrcode"
)

// PutServiceProvider интерфейс методов бизнес логики слоя Put.
type PutServiceProvider interface {
	Put(ctx context.Context, url string, userid string) (key string, err error)
	PutBatch(ctx context.Context, dc models.BatchRequest, userid string) (ec []models.BatchResponse, err error)
}

// PutHandler структура для конструктура обработчика.
type PutHandler struct {
	service PutServiceProvider
	base    string
}

// NewPutHandler конструктор обработчика.
func NewPutHandler(s PutServiceProvider, base string) *PutHandler {
	return &PutHandler{
		s,
		base,
	}
}

// Put метод обработки POST запроса с text URL в теле и возврат короткого URL в теле.
func (hn PutHandler) Put(w http.ResponseWriter, r *http.Request) {
	// получаем значение userid из контекста запроса
	userid := r.Context().Value(settings.CtxKeyUserID).(string)

	// читаем Body
	var bf bytes.Buffer
	_, err := io.Copy(&bf, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b := bf.String()

	// не эффективные по памяти варианты чтения Body

	//now := bufio.NewScanner(r.Body)
	//now.Scan()
	//b := now.Text()

	// не эффективные варианты чтения Body 2

	//bs, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//b := string(bs)

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
	key, err := hn.service.Put(ctx, b, userid)
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
	_, err = w.Write([]byte(hn.base + "/" + key))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

// PutJSON метод обработки POST запроса с одним JSON URL в теле и возврат короткого URL JSON в теле
func (hn PutHandler) PutJSON(w http.ResponseWriter, r *http.Request) {
	// получаем значение iserid из контекста запроса
	userid := r.Context().Value(settings.CtxKeyUserID).(string)
	// десериализация тела запроса
	var dc models.DecodeJSON
	err := json.NewDecoder(r.Body).Decode(&dc)
	if err != nil && err != io.EOF {
		log.Printf("unmarshal PutJSON error: %s", err)
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
	// струтктура сериализация тела запроса
	var ec models.EncodeJSON
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "application/json; charset=utf-8")
	// создаем ключ, userid token, ошибку создания в случае налияи URL в базе
	key, err := hn.service.Put(ctx, dc.URL, userid)
	// устанавливаем статус-код 201, 500 или 409
	switch {
	case err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation):
		ec.Result = hn.base + "/" + key
		w.WriteHeader(http.StatusConflict)
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		ec.Result = hn.base + "/" + key
		w.WriteHeader(http.StatusCreated)
	}
	// пишем тело ответа
	err = json.NewEncoder(w).Encode(ec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

// PutBatch метод обработки POST запроса с JSON batch в теле и возврат Batch JSON c короткими URL
func (hn PutHandler) PutBatch(w http.ResponseWriter, r *http.Request) {
	// получаем значение iserid из контекста запроса
	userid := r.Context().Value(settings.CtxKeyUserID).(string)
	// наследуем контекcт запроса r *http.Request, оснащая его Timeout
	ctx, cancel := context.WithTimeout(r.Context(), settings.StorageTimeout)
	// не забываем освободить ресурс
	defer cancel()
	// десериализация тела запроса
	var dc models.BatchRequest
	err := json.NewDecoder(r.Body).Decode(&dc)
	if err != nil && err != io.EOF {
		log.Printf("unmarshal error: %s", err)
		http.Error(w, "invalid JSON structure received", http.StatusBadRequest)
	}
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "application/json; charset=utf-8")
	//устанавливаем статус-код 201, 500 или 409

	// запрос на получение correlation_id  - original_url
	ec, err := hn.service.PutBatch(ctx, dc, userid)
	switch {
	case err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation):
		w.WriteHeader(http.StatusConflict)
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusCreated)
	}
	// пишем тело ответа
	err = json.NewEncoder(w).Encode(ec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}
