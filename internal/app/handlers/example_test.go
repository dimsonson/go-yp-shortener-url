package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers/servicemock"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

func Example() {
	s := &servicemock.ServiceMock{}
	h := NewPutHandler(s, "http://localhost:8080")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://localhost:8080/0xyz"))
	w := httptest.NewRecorder()
	// Создание контекста id пользователя для передачи хендлером в сервис.
	req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "ok"))
	// запускаем сервер
	h.Put(w, req)

	fmt.Println(w.Code)

	// Output:
	// 201
}
