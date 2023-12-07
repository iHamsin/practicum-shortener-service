package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func WithCompressionResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(rw, r)
			return
		}

		logrus.Debug("Will zip response")

		gz, gzipError := gzip.NewWriterLevel(rw, gzip.BestCompression)
		if gzipError != nil {
			http.Error(rw, gzipError.Error(), http.StatusBadRequest)
			return
		}
		defer gz.Close()

		rw.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: rw, Writer: gz}, r)

	})
}
