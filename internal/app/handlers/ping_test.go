package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers/servicemock"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
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
			name:                  "Positive test - Ping - OK",
			inputMetod:            http.MethodGet,
			inputEndpoint:         "/ping",
			inputBody:             "",
			inputUserID:           "ok",
			expectedStatusCode:    http.StatusOK,
			expectedResponseBody:  "https://pkg.go.dev/io#Reader",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "",
		},
		{
			name:                  "Negativae test - Ping - server error",
			inputMetod:            http.MethodGet,
			inputEndpoint:         "/ping",
			inputBody:             "",
			inputUserID:           "bad",
			expectedStatusCode:    http.StatusInternalServerError,
			expectedResponseBody:  "",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "",
		},
	}
	s := &servicemock.ServiceMock{}
	h := NewPingHandler(s, "http://localhost:8080")
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			// Cоздание тестирующего запроса и рекордерю
			req := httptest.NewRequest(tt.inputMetod, tt.inputEndpoint, strings.NewReader(tt.inputBody))
			w := httptest.NewRecorder()
			// Создание контекста id пользователя для передачи хендлером в сервис.
			req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, tt.inputUserID))
			// Запускаем хендлер.
			h.Ping(w, req)
			// проверяем код ответа
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			// заголовок ответа
			assert.Contains(t, w.Header().Get(tt.expectedHeader), tt.expectedHeaderContent)
		})
	}
}
