package service_test

import (
	"context"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/service"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/service/storagemock"
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
			InputUserid:   "123",
			expectedKey:   "8xyz9k",
			expectedError: nil,
		},
		/* 		{
		
		   		{
		   			name:                  "Negativae test POST - Put long URL for short - notUniq",
		   			inputMetod:            http.MethodPost,
		   			inputEndpoint:         "/",
		   			inputBody:             "https://pkg.go.dev/io#Reader",
		   			inputUserID:           "notUniq",
		   			expectedStatusCode:    http.StatusConflict,
		   			expectedResponseBody:  "http://localhost:8080/",
		   			expectedHeader:        "Content-Type",
		   			expectedHeaderContent: "text/plain; charset=utf-8",
		   		},*/
	}
	s := &storagemock.StorageMock{}
	rand := &storagemock.RandMock{}
	svc := service.NewPutService(s, "", rand)

	ctx := context.Background()

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			key, err := svc.Put(ctx, "", "")

			// проверяем код ответа
			assert.Equal(t, tt.expectedKey, key)
			// получаем и проверяем тело ответа
			assert.Equal(t, tt.expectedError, err)
			// проверка содержания строки в теле ответа
			/* 			assert.Containsf(t, string(key), tt.expectedResponseBody, "error message %s", "formatted")
			   			// заголовок ответа
			   			assert.Contains(t, (tt.expectedHeader), tt.expectedHeaderContent) */
		})
	}
}
