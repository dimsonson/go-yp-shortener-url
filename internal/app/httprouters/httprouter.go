package httprouters

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(hn *handlers.Handler) chi.Router { // http.Handler {
	// chi роутер
	rout := chi.NewRouter()
	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	rout.Use(middleware.Logger)
	//rout.Use(middleware.Compress(1)) //, "/*"))
	rout.Use(middleware.Recoverer)
	rout.Use(gzipHandle)

	// маршрут GET "/{id}" id в URL
	rout.Get("/{id}", hn.HandlerGetShortURL)
	// маршрут POST "/api/shorten" c JSON в теле запроса
	rout.Post("/api/shorten", hn.HandlerCreateShortJSON)
	// маршрут POST "/" с текстовым URL в теле запроса
	rout.Post("/", hn.HandlerCreateShortURL)
	// возврат ошибки 400 для всех остальных запросов
	rout.HandleFunc("/*", hn.IncorrectRequests)

	return rout //gz
}

type gzipWriter struct {
	http.ResponseWriter
	gzipWriter io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.gzipWriter.Write(b)
}

type gzipBodyReader struct {
	http.Request
	gzipBody io.ReadCloser
}

func (r gzipBodyReader) Close() error {
	//
	return r.gzipBody.Close()
}

func (r gzipBodyReader) Read(b []byte) (int, error) {
	//
	return r.gzipBody.Read(b)
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается клиентом, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}
		// создаём gzip.Writer поверх текущего w для записи сжатого ответа
		gzW, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			log.Println("gzip encodimg error:", err) //io.WriteString(w, err.Error())
			return
		}
		defer gzW.Close()

		// проверяем, получены ли сжатые gzip данные
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// если не использован gzip в запросе, передаём управление дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// читаем и распаковываем тело запроса с gzip
		gzRb, err := gzip.NewReader(r.Body)
		if err != nil {
			log.Println("request body decoding error", err)
			return
		}
		defer gzRb.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter и r с расшиброванным body
		next.ServeHTTP(gzipWriter{ResponseWriter: w, gzipWriter: gzW}, r) //, gzipBody: gzRb})//r) //gzipReader{Request: r, gzipBody: gzRb})
	})
}
