package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers/servicemock"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/go-chi/chi/v5"
)

var s = &servicemock.ServiceMock{}

func ExamplePutHandler_Put() {
	hPut := handlers.NewPutHandler(s, "http://localhost:8080")

	wPut := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://localhost:8080/0xyz"))
	req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "ok"))
	// запускаем хендлер
	hPut.Put(wPut, req)
	fmt.Println(wPut.Code)

	wPutJSON := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"https://pkg.go.dev/io#Reader"}`))
	req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "ok"))
	// запускаем хендлер
	hPut.PutJSON(wPutJSON, req)
	fmt.Println(wPutJSON.Code)

	wPutBatch := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(`[{"correlation_id": "05id", "original_url": "https://ya.ru/"}]`))
	req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "ok"))
	// запускаем хендлер
	hPut.PutBatch(wPutBatch, req)
	fmt.Println(wPutBatch.Code)

	// Output:
	// 201
	// 201
	// 201

}

func ExampleGetHandler_Get() {
	hGet := handlers.NewGetHandler(s, "http://localhost:8080")

	wGet := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/{id}", nil)
	// Создание контекста id пользователя для передачи хендлером в сервис.
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strings.TrimPrefix("/xyz", "/"))
	req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "ok"))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	// запускаем хендлер
	hGet.Get(wGet, req)
	fmt.Println(wGet.Code)

	wGetBatch := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "ok"))
	// запускаем хендлер
	hGet.GetBatch(wGetBatch, req)
	fmt.Println(wGetBatch.Code)

	// Output:
	// 307
	// 200
}

func ExampleDeleteHandler_Delete() {
	hDelete := handlers.NewDeleteHandler(s, "http://localhost:8080")

	wDelete := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "ok"))
	// запускаем хендлер
	hDelete.Delete(wDelete, req)
	fmt.Print(wDelete.Code)

	// Output:
	// 202
}
