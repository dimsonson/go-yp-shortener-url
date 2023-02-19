package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/api/http/handlers/servicemock"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// TestGet тест хендлера Get
func TestGet(t *testing.T) {
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
			name:                  "Positive test - Get - OK",
			inputMetod:            http.MethodGet,
			inputEndpoint:         "/xyz",
			inputBody:             "",
			inputUserID:           "ok",
			expectedStatusCode:    http.StatusTemporaryRedirect,
			expectedResponseBody:  "https://pkg.go.dev/io#Reader",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "text/plain; charset=utf-8",
		},
		{
			name:                  "Negativae test - Get - server error",
			inputMetod:            http.MethodGet,
			inputEndpoint:         "/bad",
			inputBody:             "",
			inputUserID:           "bad",
			expectedStatusCode:    http.StatusBadRequest,
			expectedResponseBody:  "",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "text/plain; charset=utf-8",
		},
		{
			name:                  "Negativae test - Get - notUniq",
			inputMetod:            http.MethodGet,
			inputEndpoint:         "/del",
			inputBody:             "",
			inputUserID:           "del",
			expectedStatusCode:    http.StatusGone,
			expectedResponseBody:  "",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "text/plain; charset=utf-8",
		},
	}
	s := &servicemock.ServiceMock{}
	h := NewGetHandler(s, "http://localhost:8080")
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			// Cоздание тестирующего запроса и рекордерю
			req := httptest.NewRequest(tt.inputMetod, tt.inputEndpoint, strings.NewReader(tt.inputBody))
			w := httptest.NewRecorder()
			// Создание контекста id пользователя для передачи хендлером в сервис.
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strings.TrimPrefix(tt.inputEndpoint, "/"))
			// req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, tt.inputUserID))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			// Запускаем хендлер.
			h.Get(w, req)
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

// TestGetBatch тест хендлера GetBatch
func TestGetBatch(t *testing.T) {
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
			name:                  "Positive test - GetBatch - OK",
			inputMetod:            http.MethodGet,
			inputEndpoint:         "/api/user/urls",
			inputBody:             "",
			inputUserID:           "ok",
			expectedStatusCode:    http.StatusOK,
			expectedResponseBody:  "https://pkg.go.dev/io#Reader",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "application/json; charset=utf-8",
		},
		{
			name:                  "Negativae test - Get - no content",
			inputMetod:            http.MethodGet,
			inputEndpoint:         "/api/user/urls",
			inputBody:             "",
			inputUserID:           "bad",
			expectedStatusCode:    http.StatusNoContent,
			expectedResponseBody:  "",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "text/plain; charset=utf-8",
		},
	}
	s := &servicemock.ServiceMock{}
	h := NewGetHandler(s, "http://localhost:8080")
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			// Cоздание тестирующего запроса и рекордерю
			req := httptest.NewRequest(tt.inputMetod, tt.inputEndpoint, strings.NewReader(tt.inputBody))
			w := httptest.NewRecorder()
			// Создание контекста id пользователя для передачи хендлером в сервис.
			req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, tt.inputUserID))
			// Запускаем хендлер.
			h.GetBatch(w, req)
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
