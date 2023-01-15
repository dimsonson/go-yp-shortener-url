package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/service"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/service/storagemock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name          string
		inputURL      string
		InputUserid   string
		expectedKey   string
		expectedDel   bool
		expectedError error
	}{
		{
			name:          "Positive test Put service level - OK",
			inputURL:      "https://pkg.go.dev/io#Reader",
			InputUserid:   "del",
			expectedKey:   "",
			expectedDel:   true,
			expectedError: nil,
		},
		{
			name:          "Negative test Put service level - Server error",
			inputURL:      "https://pkg.go.dev/io#Reader",
			InputUserid:   "bad",
			expectedKey:   "",
			expectedDel:   false,
			expectedError: errors.New("bad"),
		},
	}
	s := &storagemock.StorageMock{}
	svc := service.NewGetService(s, "http://localhost:8080")
	ctx := context.Background()

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			key, del, err := svc.Get(ctx, tt.InputUserid)
			// проверяем
			assert.Equal(t, tt.expectedKey, key)
			// получаем
			assert.Equal(t, tt.expectedError, err)
			// проверяем
			assert.Equal(t, tt.expectedDel, del)
		})
	}
}

func TestGetBatch(t *testing.T) {
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name          string
		inputURL      string
		InputUserid   string
		expectedKey   map[string]string
		expectedDel   bool
		expectedError error
	}{
		{
			name:          "Positive test Put service level - OK",
			inputURL:      "https://pkg.go.dev/io#Reader",
			InputUserid:   "ok",
			expectedKey:   map[string]string{"xyz": "https://pkg.go.dev/io#Reader"},
			expectedError: nil,
		},
		{
			name:          "Negative test Put service level - Server error",
			inputURL:      "https://pkg.go.dev/io#Reader",
			InputUserid:   "bad",
			expectedKey:   nil,
			expectedError: errors.New("noContent"),
		},
	}
	s := &storagemock.StorageMock{}
	svc := service.NewGetService(s, "http://localhost:8080")
	ctx := context.Background()

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			key, err := svc.GetBatch(ctx, tt.InputUserid)
			// проверяем
			assert.Equal(t, tt.expectedKey, key)
			// получаем
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
