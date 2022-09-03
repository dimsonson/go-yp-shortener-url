package handlers_test

import (
	"net/http"
	"testing"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
)

func TestPostHandler(t *testing.T) {
	// определяем структуру теста
	type want struct {
		code        int
		response    string
		contentType string
	}

	type req struct {
		metod string
		endpoint string
		body string
	}
	
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name string
		req string
		want want
	}{
		// определяем все тесты
		{
			name: "positive test #1",
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},
	}
	
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PostHandler(tt.args.w, tt.args.r)
		})
	}
}
