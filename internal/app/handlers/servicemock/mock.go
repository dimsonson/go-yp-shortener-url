// Package servicemock пакет заглушек обращения хендлеров к сервисам.
package servicemock

import (
	"context"
	"errors"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/jackc/pgerrcode"
)

// ServiceMock структура заглушки обращения хендлеров к сервисам.
type ServiceMock struct {
}

// Put метод реализованный для запросов через заглушку.
func (s *ServiceMock) Put(ctx context.Context, url string, userid string) (key string, err error) {
	switch userid {
	case "ok":
		return "0xyz", nil
	case "srv":
		return "", errors.New("server error")
	case "notUniq":
		return "", errors.New(pgerrcode.UniqueViolation)
	}

	return "", nil
}

// PutBatch метод реализованный для запросов через заглушку.
func (s *ServiceMock) PutBatch(ctx context.Context, dc models.BatchRequest, userid string) (ec []models.BatchResponse, err error) {
	switch userid {
	case "ok":
		ec = []models.BatchResponse{{CorrelationID: "05id", ShortURL: "http://localhost:8080/0xyz"}}
		return ec, nil
	case "srv":
		return nil, errors.New("server error")
	case "notUniq":
		return ec, errors.New(pgerrcode.UniqueViolation)
	}

	return nil, nil
}

// Get метод реализованный для запросов через заглушку.
func (s *ServiceMock) Get(ctx context.Context, id string) (value string, del bool, err error) {
	switch id {
	case "xyz":
		return "https://pkg.go.dev/io#Reader", false, nil
	case "bad":
		return "", false, errors.New("bad")
	case "del":
		return "", true, nil
	}
	return "", true, err
}

// GetBatch метод реализованный для запросов через заглушку.
func (s *ServiceMock) GetBatch(ctx context.Context, userid string) (userURLsMap map[string]string, err error) {
	kv := map[string]string{"xyz": "https://pkg.go.dev/io#Reader"}
	switch userid {
	case "ok":
		return kv, nil
	case "bad":
		return nil, errors.New("noContent")
	}
	return nil, err
}

// Ping метод реализованный для запросов через заглушку.
func (s *ServiceMock) Ping(ctx context.Context) (bool, error) {
	userid := ctx.Value(settings.CtxKeyUserID).(string)
	switch userid {
	case "ok":
		return true, nil
	case "bad":
		return false, errors.New("DB error")
	}
	return true, nil
}

// Delete метод реализованный для запросов через заглушку.
func (s *ServiceMock) Delete(shURLs []([2]string)) {
}

func  (sr *ServiceMock)RandSeq(n int) (random string, ok error) {
	return "", ok
}