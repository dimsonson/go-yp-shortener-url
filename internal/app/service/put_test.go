package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/service"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/service/storagemock"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"
)

func TestPut(t *testing.T) {
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name          string
		inputURL      string
		InputUserid   string
		expectedKey   string
		expectedError error
	}{
		{
			name:          "Positive test Put service level - OK",
			inputURL:      "https://pkg.go.dev/io#Reader",
			InputUserid:   "ok",
			expectedKey:   "8xyz9k",
			expectedError: nil,
		},
		{
			name:          "Negaive test Put service level - error UniqueViolation",
			inputURL:      "https://pkg.go.dev/io#Reader",
			InputUserid:   "notUniq",
			expectedKey:   "8xyz1000k",
			expectedError: errors.New(pgerrcode.UniqueViolation),
		},
		{
			name:          "Negative test Put service level - Server error",
			inputURL:      "https://pkg.go.dev/io#Reader",
			InputUserid:   "srv",
			expectedKey:   "",
			expectedError: errors.New("server error"),
		},
	}
	s := &storagemock.StorageMock{}
	rand := &storagemock.RandMock{}
	svc := service.NewPutService(s, "", rand)

	ctx := context.Background()

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			key, err := svc.Put(ctx, "", tt.InputUserid)
			// проверяем
			assert.Equal(t, tt.expectedKey, key)
			// получаем
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestPutBatch(t *testing.T) {
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name          string
		inputURL      string
		InputUserid   string
		expectedKey   []models.BatchResponse
		expectedError error
	}{
		{
			name:          "Positive test Put service level - OK",
			inputURL:      "https://pkg.go.dev/io#Reader",
			InputUserid:   "ok",
			expectedKey:   []models.BatchResponse{{CorrelationID: "05d", ShortURL: "http://localhost:8080/0xyz"}},
			expectedError: nil,
		},
		{
			name:          "Negaive test Put service level - error UniqueViolation",
			inputURL:      "https://pkg.go.dev/io#Reader",
			InputUserid:   "notUniq",
			expectedKey:   nil,
			expectedError: errors.New(pgerrcode.UniqueViolation),
		},

		{
			name:          "Negative test Put service level - Server error",
			inputURL:      "https://pkg.go.dev/io#Reader",
			InputUserid:   "srv",
			expectedKey:   nil,
			expectedError: errors.New("server error"),
		},
	}
	s := &storagemock.StorageMock{}
	rand := &storagemock.RandMock{}
	svc := service.NewPutService(s, "http://localhost:8080", rand)

	ctx := context.Background()

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			key, err := svc.PutBatch(ctx, nil, tt.InputUserid)
			// проверяем
			assert.Equal(t, tt.expectedKey, key)
			// получаем
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
