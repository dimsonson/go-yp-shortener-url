package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
)

func TestDefHandler(t *testing.T) {
	// определяем структуру теста
	type want struct {
		code int
	}

	type req struct {
		metod    string
		endpoint string
		body     string
	}

	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name string
		req  req
		want want
	}{
		// определяем все тесты
		{
			name: "GET #1",
			req: req{
				metod:    "GET",
				endpoint: "/xyz",
				body:     "",
			},
			want: want{
				code: 400,
			},
		},
	}

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {

			//создаем тестирующий запрос
			request := httptest.NewRequest(http.MethodGet, tt.req.endpoint, nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()

			// определяем хендлер
			h := http.HandlerFunc(handlers.DefHandler)

			// запускаем сервер
			h.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			// проверяем код ответа
			if resp.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

		})
	}
}
