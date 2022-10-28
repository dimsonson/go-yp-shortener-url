package handlers_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/services"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

func TestHandlerGetShortURL(t *testing.T) {
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
				contentType: "text/plain; charset=utf-8",
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
		
			req := httptest.NewRequest(tt.req.metod, tt.req.endpoint, nil) //strings.NewReader("http://localhost:8080/"))

			// создаём новый Recorder

			w := httptest.NewRecorder()

			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, settings.StorageTimeout)
			defer cancel()

			// определяем хендлер
			s := storage.NewMapStorage(make(map[string]string), make(map[string]string))
			srvs := services.NewService(s)
			h := handlers.NewHandler(srvs, "")
			r := httprouters.NewRouter(h)
			//	h := http.HandlerFunc(handlers.NewHandler())
			s.PutToStorage(ctx, 1, "xyz", "https://pkg.go.dev/github.com/stretchr/testify@v1.8.0/assert#Containsf")

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strings.TrimPrefix(tt.req.endpoint, "/"))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// запускаем сервер
			r.ServeHTTP(w, req)
			resp := w.Result()

			defer resp.Body.Close()

			// проверяем код ответа
			if resp.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// заголовок ответа
			if resp.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, resp.Header.Get("Content-Type"))
			}
		})
	}
}
