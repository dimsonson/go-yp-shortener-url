package httprouters_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
)

func TestHTTPRouter(t *testing.T) {

	// определяем структуру теста
	type want struct {
		handlerOutStatus int
	}

	type req struct {
		methodIn string
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
				methodIn: "POST",
			},
			want: want{
				handlerOutStatus: 201,
			},
		},
		{
			name: "DEF #1",
			req: req{
				methodIn: "PATCH",
			},
			want: want{
				handlerOutStatus: 400,
			},
		},
		{
			name: "GET #1",
			req: req{
				methodIn: "GET",
			},
			want: want{
				handlerOutStatus: 400, //307,
			},
		},
	}

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {

			//создаем тестирующий запрос
			request := httptest.NewRequest(tt.req.methodIn, "/", nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()

			// определяем хендлер
			h := http.HandlerFunc(httprouters.HTTPRouter)

			// запускаем сервер
			h.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			// проверяем код ответа вызываемой функции
			if tt.want.handlerOutStatus != resp.StatusCode {
				t.Errorf("Expected status code %d, got %d", tt.want.handlerOutStatus, resp.StatusCode)
			}
		})
	}
}
