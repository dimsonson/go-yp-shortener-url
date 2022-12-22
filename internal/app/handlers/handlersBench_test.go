package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers/servicemock"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

func BenchmarkHandlers(b *testing.B) {
	s := &servicemock.ServiceMock{}
	h := NewHandler(s, "http://localhost:8080")

	b.Run("CreateShortURL", func(b *testing.B) {
		// конфигурируем аналоги http.Request и http.ResponseWriter для бенчмарка
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://ya.ru"))
		w := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff"))
		// обнуляемтаймер бенчмарка и запукаем его
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.HandlerCreateShortURL(w, req)
		}
	})

	b.Run("CreateShortJSON", func(b *testing.B) {
		// конфигурируем аналоги http.Request и http.ResponseWriter для бенчмарка
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"https://yandex.ru"}`))
		w := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff"))
		// обнуляемтаймер бенчмарка и запукаем его
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.HandlerCreateShortJSON(w, req)
		}
	})

	b.Run("GetUserURLs", func(b *testing.B) {
		// конфигурируем аналоги http.Request и http.ResponseWriter для бенчмарка
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", strings.NewReader(`{"url":"https://yandex.ru"}`))
		w := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff"))
		// обнуляемтаймер бенчмарка и запукаем его
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.HandlerGetUserURLs(w, req)
		}
	})

	b.Run("IncorrectReques1ts", func(b *testing.B) {
		// конфигурируем аналоги http.Request и http.ResponseWriter для бенчмарка
		req := httptest.NewRequest(http.MethodPatch, "/api/user/urls", strings.NewReader(`{"url":"https://yandex.ru"}`))
		w := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff"))
		// обнуляемтаймер бенчмарка и запукаем его
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.IncorrectRequests(w, req)
		}
	})

	b.Run("GetShortURL", func(b *testing.B) {
		// конфигурируем аналоги http.Request и http.ResponseWriter для бенчмарка
		req := httptest.NewRequest(http.MethodGet, "/0124", nil)
		w := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff"))
		// обнуляемтаймер бенчмарка и запукаем его
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.HandlerGetShortURL(w, req)
		}
	})

	b.Run("CreateBatchJSON", func(b *testing.B) {
		// конфигурируем аналоги http.Request и http.ResponseWriter для бенчмарка
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`[
			{
				"correlation_id": "05id",
				"original_url": "https://y1131223a.ru/"
			},
			{
				"correlation_id": "02id",
				"original_url": "http://ru1.wikipedia.org/wiki/Go"
			},
			{
				"correlation_id": "03id",
				"original_url": "http://tproger.ru/translations/golang-basics/"
			}
		]`))
		w := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff"))
		// обнуляемтаймер бенчмарка и запукаем его
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.HandlerCreateBatchJSON(w, req)
		}
	})

	b.Run("DeleteBatch", func(b *testing.B) {
		// конфигурируем аналоги http.Request и http.ResponseWriter для бенчмарка
		req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(`["0SNGmH", "1GSuBf", "4pLhqd", "4", "1GSuB33f"]`))
		w := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff"))
		// обнуляемтаймер бенчмарка и запукаем его
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.HandlerDeleteBatch(w, req)
		}
	})

	b.Run("SQLping", func(b *testing.B) {
		// конфигурируем аналоги http.Request и http.ResponseWriter для бенчмарка
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		w := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff"))
		// обнуляемтаймер бенчмарка и запукаем его
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.HandlerSQLping(w, req)
		}
	})
}
