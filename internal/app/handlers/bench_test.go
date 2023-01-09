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

// BenchmarkHandlers(b *testing.B) Бенчмарк хендлеров.
func BenchmarkHandlers(b *testing.B) {
	s := &servicemock.ServiceMock{}

	h := NewPutHandler(s, "http://localhost:8080")

	b.Run("Put", func(b *testing.B) {
		b.ReportAllocs()
		// конфигурируем аналоги http.Request и http.ResponseWriter для бенчмарка
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://ya.ru"))
		w := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff"))
		// обнуляемтаймер бенчмарка и запукаем его
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.Put(w, req)
		}
	})

	b.Run("PutJSON", func(b *testing.B) {
		b.ReportAllocs()
		// конфигурируем аналоги http.Request и http.ResponseWriter для бенчмарка
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"https://yandex.ru"}`))
		w := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff"))
		// обнуляемтаймер бенчмарка и запукаем его
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.PutJSON(w, req)
		}
	})

	b.Run("PutBatch", func(b *testing.B) {
		b.ReportAllocs()
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
			h.PutBatch(w, req)
		}
	})

	hG := NewGetHandler(s, "http://localhost:8080")

	b.Run("GetBatch", func(b *testing.B) {
		b.ReportAllocs()
		// конфигурируем аналоги http.Request и http.ResponseWriter для бенчмарка
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", strings.NewReader(`{"url":"https://yandex.ru"}`))
		w := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff"))
		// обнуляемтаймер бенчмарка и запукаем его
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			hG.GetBatch(w, req)
		}
	})

	b.Run("Get", func(b *testing.B) {
		b.ReportAllocs()
		// конфигурируем аналоги http.Request и http.ResponseWriter для бенчмарка
		req := httptest.NewRequest(http.MethodGet, "/0124", nil)
		w := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff"))
		// обнуляемтаймер бенчмарка и запукаем его
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			hG.Get(w, req)
		}
	})

	hD := NewDeleteHandler(s, "http://localhost:8080")

	b.Run("Delete", func(b *testing.B) {
		b.ReportAllocs()
		// конфигурируем аналоги http.Request и http.ResponseWriter для бенчмарка
		req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(`["0SNGmH", "1GSuBf", "4pLhqd", "4", "1GSuB33f"]`))
		w := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff"))
		// обнуляемтаймер бенчмарка и запукаем его
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			hD.Delete(w, req)
		}
	})

	hP := NewPingHandler(s, "http://localhost:8080")

	b.Run("SQLping", func(b *testing.B) {
		b.ReportAllocs()
		// конфигурируем аналоги http.Request и http.ResponseWriter для бенчмарка
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		w := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), settings.CtxKeyUserID, "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff"))
		// обнуляемтаймер бенчмарка и запукаем его
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			hP.Ping(w, req)
		}
	})
}
