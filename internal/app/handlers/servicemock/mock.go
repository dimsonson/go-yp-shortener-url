package servicemock

import (
	"context"
	"errors"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
	"github.com/jackc/pgerrcode"
)

type ServiceMock struct {
}

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

func (s *ServiceMock) GetBatch(ctx context.Context, userid string) (userURLsMap map[string]string, err error) {
	kv := map[string]string{"xyz":"https://pkg.go.dev/io#Reader" }
	switch userid {
	case "ok":
		return kv, nil
	case "bad":
		return nil, errors.New("noContent")
	}
	return nil, err
}

func (s *ServiceMock) Ping(ctx context.Context) (bool, error) {
	return true, nil
}

func (s *ServiceMock) Delete(shURLs []([2]string)) {

}
