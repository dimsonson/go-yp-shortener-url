package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/services"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
	"github.com/stretchr/testify/assert"
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
				metod:    http.MethodGet,
				endpoint: "/",
				body:     "https://pkg.go.dev/io#Reader",
			},
			want: want{
				code:        201,
				response:    "http://",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {

			//создаем тестирующий запрос
			request := httptest.NewRequest(tt.req.metod, tt.req.endpoint, nil) //strings.NewReader("http://localhost:8080/"))

			// создаём новый Recorder

			w := httptest.NewRecorder()

			// определяем хендлер
			s := storage.NewMapStorage("map")
			srvs := services.NewService(s)
			h := handlers.NewHandler(srvs)
			r := httprouters.NewRouter(h)
			//	h := http.HandlerFunc(handlers.NewHandler())

			// запускаем сервер
			r.ServeHTTP(w, request)
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
