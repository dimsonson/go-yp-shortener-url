package httprouters

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// gzipWriter структура для записи сжатого ответа.
type gzipWriter struct {
	http.ResponseWriter
	gzWriter io.Writer
}

// gzipWriter метод для записи сжатого ответа.
func (w gzipWriter) Write(b []byte) (int, error) {
	// w.gzWriter будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.gzWriter.Write(b)
}

// gzipReader структура для чтения сжатого тела запроса.
type gzipReader struct {
	gzipReader *gzip.Reader
	gzipBody   io.ReadCloser
}

// gzipReader метод для закрытия тела запроса.
func (r gzipReader) Close() error {
	//закрываем gzipReader
	r.gzipReader.Close()
	//закрываем тело
	return r.gzipBody.Close()
}

// Read метод для чтения зашифрованного тела запроса.
func (r gzipReader) Read(b []byte) (int, error) {
	return r.gzipReader.Read(b)
}

// middlewareGzip функция распаковки-сжатия http алгоритмом gzip для использования в роутере.
func middlewareGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что запрос содежит сжатые данные
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// читаем и распаковываем тело запроса с gzip
			gzR, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Print("gzip error: ", err)
				return
			}
			defer gzR.Close()
			r.Body = gzipReader{gzipReader: gzR, gzipBody: r.Body}
			defer r.Body.Close()
		}
		// проверяем, что клиент поддерживает gzip-сжатие
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// создаём gzip.Writer поверх текущего w для записи сжатого ответа
			gzW, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				log.Print("gzip encodimg error:", err)
				return
			}
			defer gzW.Close()
			// устанавливаем заголовок сжатия содержимого ответа
			w.Header().Set("Content-Encoding", "gzip")
			// отправляем ответ с сжатым содержанием
			next.ServeHTTP(gzipWriter{ResponseWriter: w, gzWriter: gzW}, r)
			return
		}
		// если gzip не поддерживается клиентом, передаём управление дальше без изменений
		next.ServeHTTP(w, r)
	})
}
