package httprouters

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

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

func middlewareGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// читаем и распаковываем тело запроса с gzip
			var err error
			r.Body, err = gzip.NewReader(r.Body)
			if err != nil {
				log.Println("request body decoding error", err)
				next.ServeHTTP(w, r)
				return
			}
			defer r.Body.Close()
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
