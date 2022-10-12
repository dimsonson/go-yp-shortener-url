package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	//"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/services"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGetUserURLs(t *testing.T) {
	// определяем структуру теста
	type want struct {
		code        int
		response    string
		contentType string
	}

	type req struct {
		metod    string
		endpoint string
		cookie   http.Cookie
		id int
		body     string
	}

	//var cook http.Cookie
	cook := http.Cookie{
		Name:   "token",
		Value:  "00000000b38aaf6c89467a765a15a5d40098d050c80503562bebef1c64ded15cc4fbdaeb",
		MaxAge: 300,
	}
	cookBad := http.Cookie{}

	
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name string
		req  req
		want want
	}{
		// определяем все тесты
		{
			name: "GET - set of URLs Successfully created",
			req: req{
				metod:    "GET",
				endpoint: "/api/user/urls",
				cookie:   cook,
				id : 0,
				body:     `{"url":"https://yandex.ru/search/?text=AToi+go&lr=213"}`,
			},
			want: want{
				code:        200,
				response:    "[",
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name: "GET wrong cookie",
			req: req{
				metod:    "GET",
				endpoint: "/api/user/urls",
				cookie:   cookBad,
				id : 5,
				body:     "",
			},
			want: want{
				code:     204,
				response: "",

				contentType: "text/plain; charset=utf-8",
			},
		},
		
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {

			//создаем тестирующий запрос
			req := httptest.NewRequest(tt.req.metod, tt.req.endpoint, strings.NewReader(tt.req.body))
			// установим куку в ответ
			req.AddCookie(&cook)

			defer req.Body.Close()

			// создаём новый Recorder

			w := httptest.NewRecorder()

			// определяем хендлер
			s := storage.NewMapStorage(make(map[string]int), make(map[string]string))
			s.PutToStorage(tt.req.id, "xyz", "https://pkg.go.dev/github.com/stretchr/testify@v1.8.0/assert#Containsf")
			srvs := services.NewService(s)
			h := handlers.NewHandler(srvs, "")
			//r := httprouters.NewRouter(h)
			h.HandlerGetUserURLs(w, req)
			fmt.Println(s.UserID)
			// запускаем сервер
			//r.ServeHTTP(w, req)
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
