package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/service"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/service/storagemock"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name          string
		inputURL      string
		InputUserID   string
		expectedKey   string
		expectedOk    bool
		expectedError error
	}{
		{
			name:          "Positive test Put service level - OK",
			inputURL:      "https://pkg.go.dev/io#Reader",
			InputUserID:   "ok",
			expectedKey:   "",
			expectedOk:    true,
			expectedError: nil,
		},
		{
			name:          "Negative test Put service level - Server error",
			inputURL:      "https://pkg.go.dev/io#Reader",
			InputUserID:   "bad",
			expectedKey:   "",
			expectedOk:    false,
			expectedError: errors.New("DB error"),
		},
	}
	s := &storagemock.StorageMock{}
	svc := service.NewPingService(s)

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			// Создание контекста id пользователя для передачи хендлером в сервис.
			ctx := (context.WithValue(context.Background(), settings.CtxKeyUserID, tt.InputUserID))
			// Запрос.
			ok, err := svc.Ping(ctx)
			// получаем
			assert.Equal(t, tt.expectedError, err)
			// проверяем
			assert.Equal(t, tt.expectedOk, ok)
		})
	}
}
