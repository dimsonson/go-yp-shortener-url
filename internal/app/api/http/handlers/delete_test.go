package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/api/http/handlers/servicemock"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/stretchr/testify/assert"
)

// TestDelete тест хендлера Delete.
func TestDelete(t *testing.T) {
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
			name:                  "Positive test - Delete - OK",
			inputMetod:            http.MethodDelete,
			inputEndpoint:         "/api/user/urls",
			inputBody:             "",
			inputUserID:           "ok",
			expectedStatusCode:    http.StatusAccepted,
			expectedResponseBody:  "",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "application/json; charset=utf-8",
		},
		{
			name:                  "Negativae test - Delete - no content",
			inputMetod:            http.MethodDelete,
			inputEndpoint:         "/api/user/urls",
			inputBody:             "{}",
			inputUserID:           "bad",
			expectedStatusCode:    http.StatusBadRequest,
			expectedResponseBody:  "",
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "application/json; charset=utf-8",
		},
	}
	s := &servicemock.ServiceMock{}
	h := NewDeleteHandler(s, "http://localhost:8080")
	for _, tt := range tests {
		// Запускаем каждый тест.
		t.Run(tt.name, func(t *testing.T) {
			// Cоздание тестирующего запроса и рекордер.
			req := httptest.NewRequest(tt.inputMetod, tt.inputEndpoint, strings.NewReader(tt.inputBody))
			w := httptest.NewRecorder()
			// Создание контекста id пользователя для передачи хендлером в сервис.
			req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, tt.inputUserID))
			// Запускаем хендлер.
			h.Delete(w, req)
			// Проверяем код ответа.
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			// Заголовок ответа.
			assert.Contains(t, w.Header().Get(tt.expectedHeader), tt.expectedHeaderContent)
		})
	}
}
