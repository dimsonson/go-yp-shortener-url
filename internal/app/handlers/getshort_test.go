package handlers_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetShortURL(t *testing.T) {
	// определяем структуру теста
	type want struct {
		code        int
		response    string
		contentType string
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
				code:        307,
				response:    "https://pkg.go.dev/github.com/stretchr/testify@v1.8.0/assert#Containsf",
				contentType: "text/html; charset=utf-8",
			},
		},
		{
			name: "GET #2",
			req: req{
				metod:    "GET",
				endpoint: "/",
				body:     "",
			},
			want: want{
				code:        400,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {

			//создаем запись в базе url
			storage.DB["xyz"] = "https://pkg.go.dev/github.com/stretchr/testify@v1.8.0/assert#Containsf"

			r := httptest.NewRequest(http.MethodGet, tt.req.endpoint, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strings.TrimPrefix(tt.req.endpoint, "/"))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			//создаем тестирующий запрос
			//request := httptest.NewRequest(http.MethodGet, tt.req.endpoint, nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()

			// определяем хендлер
			handlers.GetShortURL(w, r)
			// записываем ответ
			resp := w.Result()

			// проверяем код ответа
			if resp.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело ответа
			defer resp.Body.Close()
			resBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			// проверка содержания строки в теле ответа
			assert.Containsf(t, string(resBody), tt.want.response, "error message %s", "formatted")

			// заголовок ответа
			if resp.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, resp.Header.Get("Content-Type"))
			}
		})
	}
}