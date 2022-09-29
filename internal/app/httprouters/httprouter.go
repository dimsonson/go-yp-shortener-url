package httprouters

import (
	"compress/gzip"
	"fmt"
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
	gzWriter io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.gzWriter.Write(b)
}

type gzipReader struct {
	gzipReader io.Reader
	gzipBody   io.ReadCloser
}

func (r gzipReader) Close() error {
	//
	return r.gzipBody.Close()
}

func (r gzipReader) Read(b []byte) (int, error) {
	//
	return r.gzipBody.Read(b)
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// читаем и распаковываем тело запроса с gzip
			gzRb, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Println("request body decoding error", err)
				next.ServeHTTP(w, r)
				return
			}
			defer gzRb.Close()
			//
			data, err := io.ReadAll(gzRb)
			if err != nil {
				log.Println(err)
			}
			//
			r.Body.Read(data)
			fmt.Println(r.Body)
		}
		// проверяем, что клиент поддерживает gzip-сжатие
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// создаём gzip.Writer поверх текущего w для записи сжатого ответа
			gzW, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				log.Println("gzip encodimg error:", err)
				return
			}
			defer gzW.Close()
			w.Header().Set("Content-Encoding", "gzip")
			//
			next.ServeHTTP(gzipWriter{ResponseWriter: w, gzWriter: gzW}, r)
			return
		}
		// если gzip не поддерживается клиентом, передаём управление дальше без изменений
		next.ServeHTTP(w, r)
	})
}
