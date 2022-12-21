package servicemock

import (
	"context"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
)

type ServiceMock struct {
}

func (s *ServiceMock) ServiceCreateShortURL(ctx context.Context, url string, userid string) (key string, err error) {
	return
}

func (s *ServiceMock) ServiceGetShortURL(ctx context.Context, id string) (value string, del bool, err error) {
	return "", true, err
}

func (s *ServiceMock) ServiceGetUserShortURLs(ctx context.Context, userid string) (userURLsMap map[string]string, err error) {
	return
}

func (s *ServiceMock) ServiceStorageOkPing(ctx context.Context) (bool, error) {
	return true, nil
}

func (s *ServiceMock) ServiceCreateBatchShortURLs(ctx context.Context, dc models.BatchRequest, userid string) (ec []models.BatchResponse, err error) {
	return
}

func (s *ServiceMock) ServiceDeleteURL(shURLs []([2]string)) {

}
