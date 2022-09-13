package handlers_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/services"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
	"github.com/stretchr/testify/assert"
)

/* type HandlerTest interface {
	ServiceCreateShortURL(url string) (key string)
	ServiceGetShortURL(id string) (value string, err error)
}

type Handler struct {
	handler Services
}

func NewHandler(s Services) *Handler {
	return &Handler{
		s,
	}
} */

func TestHandlerCreateShortURL(t *testing.T) {
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
			name: "POST #1",
			req: req{
				metod:    "POST",
				endpoint: "/",
				body:     "https://pkg.go.dev/io#Reader",
			},
			want: want{
				code:        201,
				response:    "http://example.com/",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {

			//создаем тестирующий запрос
			req := httptest.NewRequest(tt.req.metod, "/", strings.NewReader("http://localhost:8080/"))

			// создаём новый Recorder

			w := httptest.NewRecorder()

			// определяем хендлер
			s := storage.NewMapStorage("map")
			srvs := services.NewService(s)
			h := handlers.NewHandler(srvs)
			r := httprouters.NewRouter(h)
			//h := http.HandlerFunc(handlers.)
			//rctx := chi.NewRouteContext()
			//rctx.URLParams.Add("id", strings.TrimPrefix(tt.req.endpoint, "/"))
			//req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			//	h := http.HandlerFunc(handlers.NewHandler())

			// запускаем сервер
			r.ServeHTTP(w, req)
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
