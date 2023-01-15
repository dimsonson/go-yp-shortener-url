// Package servicemock пакет заглушек обращения хендлеров к сервисам.
package storagemock

import (
	"context"
	"errors"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/jackc/pgerrcode"
)

// ServiceMock структура заглушки обращения хендлеров к сервисам.
type StorageMock struct {
}

// Put метод реализованный для запросов через заглушку.
func (s *StorageMock) Put(ctx context.Context, key string, value string, userid string) (existKey string, err error) {
	switch userid {
	case "ok":
		return "8xyz9k", nil
	case "srv":
		return "", errors.New("server error")
	case "notUniq":
		return "8xyz1000k", errors.New(pgerrcode.UniqueViolation)
	}

	return "notExist", nil
}

// PutBatch метод реализованный для запросов через заглушку.
func (s *StorageMock) PutBatch(ctx context.Context, dc models.BatchRequest, userid string) (dcCorr models.BatchRequest, err error) {
	switch userid {
	case "ok":
		dcCorr = models.BatchRequest{{CorrelationID: "05d", OriginalURL: "05id", ShortURL: "0xyz"}}
		return dcCorr, nil
	case "srv":
		return nil, errors.New("server error")
	case "notUniq":
		return dcCorr, errors.New(pgerrcode.UniqueViolation)
	}

	return nil, nil
}

// Get метод реализованный для запросов через заглушку.
func (s *StorageMock) Get(ctx context.Context, key string) (value string, del bool, err error) {
	switch key {
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
func (s *StorageMock) GetBatch(ctx context.Context, userid string) (userURLs map[string]string, err error) {
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
func (s *StorageMock) Ping(ctx context.Context) (ok bool, err error) {
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
func (s *StorageMock) Delete(key string, userid string) (err error) {
	return err
}

// Len метод реализованный для запросов через заглушку.
func (s *StorageMock) Len(ctx context.Context) (lenn int) {
	return 8
}

// RandSeq структура реализованная для запросов через заглушку.
type RandMock struct {
	StorageMock
}

// RandSeq метод реализованный для запросов через заглушку.
func (s *RandMock) RandSeq(n int) (random string, ok error) {
	return "xyz9k", ok
}
