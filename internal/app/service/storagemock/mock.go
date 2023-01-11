// Package servicemock пакет заглушек обращения хендлеров к сервисам.
package storagemock

import (
	"context"
	"errors"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

// ServiceMock структура заглушки обращения хендлеров к сервисам.
type StorageMock struct {
}

// Put метод реализованный для запросов через заглушку.
func (s *StorageMock) Put(ctx context.Context, key string, value string, userid string) (existKey string, err error) {
	/* 	switch userid {
	   	case "ok":
	   		return "0xyz", nil
	   	case "srv":
	   		return "", errors.New("server error")
	   	case "notUniq":
	   		return "", errors.New(pgerrcode.UniqueViolation)
	   	}
	*/
	return "8xyz9k", nil
}

// PutBatch метод реализованный для запросов через заглушку.
func (s *StorageMock) PutBatch(ctx context.Context, dc models.BatchRequest, userid string) (dcCorr models.BatchRequest, err error) {
	/* switch userid {
	case "ok":
		ec = []models.BatchResponse{{CorrelationID: "05id", ShortURL: "http://localhost:8080/0xyz"}}
		return ec, nil
	case "srv":
		return nil, errors.New("server error")
	case "notUniq":
		return ec, errors.New(pgerrcode.UniqueViolation)
	} */

	return nil, nil
}

// Get метод реализованный для запросов через заглушку.
func (s *StorageMock) Get(ctx context.Context, key string) (value string, del bool, err error) {
	/* 	switch id {
	   	case "xyz":
	   		return "https://pkg.go.dev/io#Reader", false, nil
	   	case "bad":
	   		return "", false, errors.New("bad")
	   	case "del":
	   		return "", true, nil
	   	} */
	return "", true, err
}

// GetBatch метод реализованный для запросов через заглушку.
func (s *StorageMock) GetBatch(ctx context.Context, userid string) (userURLs map[string]string, err error) {
	/* 	kv := map[string]string{"xyz": "https://pkg.go.dev/io#Reader"}
	   	switch userid {
	   	case "ok":
	   		return kv, nil
	   	case "bad":
	   		return nil, errors.New("noContent")
	   	} */
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

func (s *StorageMock) Len(ctx context.Context) (lenn int) {

	return 8

}

func (s *StorageMock) Load() {
}

func (s *StorageMock) Close() {

}

type RandMock struct {
	StorageMock
}

func (s *RandMock) RandSeq(n int) (random string, ok error) {
	return "xyz9k", ok
}
