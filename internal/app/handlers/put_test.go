package handlers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers/servicemock"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/stretchr/testify/assert"
)

func TestPut(t *testing.T) {
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name                  string
		inputMetod            string
		inputEndpoint         string
		inputBody             string
		inputUserID           string
		expectedStatusCode    int
		expectedResponseBody  string
		expectedHeader        string
		expectedHeaderContent string
	}{
		{
			name:                  "Positive test POST - Put long URL for short - OK",
			inputMetod:            http.MethodPost,
			inputEndpoint:         "/",
			inputBody:             "https://pkg.go.dev/io#Reader",
			inputUserID:           "ok",
			expectedStatusCode:    http.StatusCreated,
			expectedResponseBody:  "http://localhost:8080/0xyz",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "text/plain; charset=utf-8",
		},
		{
			name:                  "Negativae test POST - Put long URL for short - server error",
			inputMetod:            http.MethodPost,
			inputEndpoint:         "/",
			inputBody:             "https://pkg.go.dev/io#Reader",
			inputUserID:           "srv",
			expectedStatusCode:    http.StatusInternalServerError,
			expectedResponseBody:  "http://localhost:8080/",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "text/plain; charset=utf-8",
		},
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
		},
	}
	s := &servicemock.ServiceMock{}
	h := NewPutHandler(s, "http://localhost:8080")
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			// Cоздание тестирующего запроса и рекордерю
			req := httptest.NewRequest(tt.inputMetod, tt.inputEndpoint, strings.NewReader(tt.inputBody))
			w := httptest.NewRecorder()
			// Создание контекста id пользователя для передачи хендлером в сервис.
			req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, tt.inputUserID))
			// запускаем сервер
			h.Put(w, req)
			// проверяем код ответа
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			// получаем и проверяем тело ответа
			resBody, err := io.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			// проверка содержания строки в теле ответа
			assert.Containsf(t, string(resBody), tt.expectedResponseBody, "error message %s", "formatted")
			// заголовок ответа
			assert.Contains(t, w.Header().Get(tt.expectedHeader), tt.expectedHeaderContent)
		})
	}
}

func TestPutJSON(t *testing.T) {
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name                  string
		inputMetod            string
		inputEndpoint         string
		inputBody             string
		inputUserID           string
		expectedStatusCode    int
		expectedResponseBody  string
		expectedHeader        string
		expectedHeaderContent string
	}{
		{
			name:                  "Positive test - PutJSON - OK",
			inputMetod:            http.MethodPost,
			inputEndpoint:         "/api/shorten",
			inputBody:             `{"url":"https://pkg.go.dev/io#Reader"}`,
			inputUserID:           "ok",
			expectedStatusCode:    http.StatusCreated,
			expectedResponseBody:  "http://localhost:8080/0xyz",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "application/json; charset=utf-8",
		},
		{
			name:                  "Negativae test - PutJSON server error",
			inputMetod:            http.MethodPost,
			inputEndpoint:         "/api/shorten",
			inputBody:             `{"url":"https://pkg.go.dev/io#Reader"}`,
			inputUserID:           "srv",
			expectedStatusCode:    http.StatusInternalServerError,
			expectedResponseBody:  "",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "application/json; charset=utf-8",
		},
		{
			name:                  "Negativae test - PutJSON - notUniq",
			inputMetod:            http.MethodPost,
			inputEndpoint:         "/api/shorten",
			inputBody:             `{"url":"https://pkg.go.dev/io#Reader"}`,
			inputUserID:           "notUniq",
			expectedStatusCode:    http.StatusConflict,
			expectedResponseBody:  "",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "application/json; charset=utf-8",
		},
	}
	s := &servicemock.ServiceMock{}
	h := NewPutHandler(s, "http://localhost:8080")
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			// Cоздание тестирующего запроса и рекордерю
			req := httptest.NewRequest(tt.inputMetod, tt.inputEndpoint, strings.NewReader(tt.inputBody))
			w := httptest.NewRecorder()
			// Создание контекста id пользователя для передачи хендлером в сервис.
			req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, tt.inputUserID))
			// запускаем сервер
			h.PutJSON(w, req)
			// проверяем код ответа
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			// получаем и проверяем тело ответа
			resBody, err := io.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			// проверка содержания строки в теле ответа
			assert.Containsf(t, string(resBody), tt.expectedResponseBody, "error message %s", "formatted")
			// заголовок ответа
			assert.Contains(t, w.Header().Get(tt.expectedHeader), tt.expectedHeaderContent)
		})
	}
}

func TestPutBatch(t *testing.T) {
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name                  string
		inputMetod            string
		inputEndpoint         string
		inputBody             string
		inputUserID           string
		expectedStatusCode    int
		expectedResponseBody  string
		expectedHeader        string
		expectedHeaderContent string
	}{
		{
			name:                  "Positive test - PutBatch - OK",
			inputMetod:            http.MethodPost,
			inputEndpoint:         "/api/shorten/batch",
			inputBody:             `[{"correlation_id": "05id", "original_url": "https://ya.ru/"}]`,
			inputUserID:           "ok",
			expectedStatusCode:    http.StatusCreated,
			expectedResponseBody:  "http://localhost:8080/0xyz",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "application/json; charset=utf-8",
		},
		{
			name:                  "Negativae test POST - PutBatch - server error",
			inputMetod:            http.MethodPost,
			inputEndpoint:         "/api/shorten/batch",
			inputBody:             `[{"correlation_id": "05id", "original_url": "https://ya.ru/"}]`,
			inputUserID:           "srv",
			expectedStatusCode:    http.StatusInternalServerError,
			expectedResponseBody:  "",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "application/json; charset=utf-8",
		},
		{
			name:                  "Negativae test POST - PutBatch - notUniq",
			inputMetod:            http.MethodPost,
			inputEndpoint:         "/api/shorten/batch",
			inputBody:             `[{"correlation_id": "05id", "original_url": "https://ya.ru/"}]`,
			inputUserID:           "notUniq",
			expectedStatusCode:    http.StatusConflict,
			expectedResponseBody:  "",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "application/json; charset=utf-8",
		},
	}
	s := &servicemock.ServiceMock{}
	h := NewPutHandler(s, "http://localhost:8080")
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			// Cоздание тестирующего запроса и рекордерю
			req := httptest.NewRequest(tt.inputMetod, tt.inputEndpoint, bytes.NewBufferString(tt.inputBody))
			w := httptest.NewRecorder()
			// Создание контекста id пользователя для передачи хендлером в сервис.
			req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, tt.inputUserID))
			// запускаем сервер
			h.PutBatch(w, req)
			// проверяем код ответа
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			// получаем и проверяем тело ответа
			resBody, err := io.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			// проверка содержания строки в теле ответа
			assert.Containsf(t, string(resBody), tt.expectedResponseBody, "error message %s", "formatted")
			// заголовок ответа
			assert.Contains(t, w.Header().Get(tt.expectedHeader), tt.expectedHeaderContent)
		})
	}
}
