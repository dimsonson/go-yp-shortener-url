package handlers

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers/servicemock"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

func Example_Put() {
/* 	var s Student
	s.SetName("dima")
	fmt.Println(s.Name)

	n := s.GetName()
	fmt.Println(n) */

	// Output:
	// Dima
	// Dima
}

func ExampleHandlers_Put() {
	s := &servicemock.ServiceMock{}
	h := NewPutHandler(s, "http://localhost:8080")

	req := httptest.NewRequest(tt.inputMetod, tt.inputEndpoint, strings.NewReader(tt.inputBody))
	w := httptest.NewRecorder()
	// Создание контекста id пользователя для передачи хендлером в сервис.
	req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, tt.inputUserID))
	// запускаем сервер
	h.Put(w, req)

	
	fmt.Println(w)

	
	fmt.Println(w)

	// Output:
	// Dima
	// Dima
}
