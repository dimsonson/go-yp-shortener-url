package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/jackc/pgerrcode"
)

// PutStorageProvider интерфейс методов хранилища.
type PutStorageProvider interface {
	Put(ctx context.Context, key string, value string, userid string) (existKey string, err error)
	PutBatch(ctx context.Context, dc models.BatchRequest, userid string) (dcCorr models.BatchRequest, err error)
	Len(ctx context.Context) (lenn int)
}

// PutServices структура конструктора бизнес логики.
type PutServices struct {
	storage PutStorageProvider
	base    string
	RandProvider
}

// NewPutService конструктор бизнес логики.
func NewPutService(s PutStorageProvider, base string, rand RandProvider) *PutServices {
	return &PutServices{
		s,
		base,
		rand,
	}
}

// Put метод создание пары id : URL.
func (sr *PutServices) Put(ctx context.Context, url string, userid string) (key string, err error) {
	// создаем и присваиваем значение короткой ссылки
	key, err = sr.RandSeq(settings.KeyLeght)
	if err != nil {
		log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
	}
	// добавляем уникальный префикс к ключу
	key = fmt.Sprintf("%d%s", sr.storage.Len(ctx), key)
	// создаем запись userid-ключ-значение в базе
	existKey, err := sr.storage.Put(ctx, key, url, userid)
	switch {
	case err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation):
		key = existKey
	case err != nil:
		return "", err
	}
	return key, err
}

// PutBatch метод создание пакета пар id : URL.
func (sr *PutServices) PutBatch(ctx context.Context, dc models.BatchRequest, userid string) (ec []models.BatchResponse, err error) {
	// добавление shorturl
	for i := range dc {
		key, err := sr.RandSeq(settings.KeyLeght)
		if err != nil {
			log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
		}
		key = fmt.Sprintf("%d%s", sr.storage.Len(ctx), key)
		dc[i].ShortURL = key
	}
	// пишем в базу и получаем слайс с обновленными shorturl в случае конфликта
	dc, err = sr.storage.PutBatch(ctx, dc, userid)
	switch {
	case err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation):
		break
	case err != nil:
		return nil, err
	}
	// заполняем слайс ответа
	for _, v := range dc {
		elem := models.BatchResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      sr.base + "/" + v.ShortURL,
		}
		ec = append(ec, elem)
	}
	return ec, err
}

// RandProvider интерфейс для вызова метода генератора псевдо случайной последовательности знаков.
type RandProvider interface {
	RandSeq(n int) (random string, ok error)
}

// Rand структура для вызова метода генератора псевдо случайной последовательности знаков.
type Rand struct{}

// RandSeq функция генерации псевдо случайной последовательности знаков.
func (sr *Rand) RandSeq(n int) (random string, ok error) {
	if n < 1 {
		err := fmt.Errorf("wromg argument: number %v less than 1\n ", n)
		return "", err
	}
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	random = string(b)
	return random, nil
}
