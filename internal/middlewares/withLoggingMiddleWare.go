package middlewares

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func WithLoggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: rw, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		next.ServeHTTP(&lw, r)
		duration := time.Since(start)

		logrus.WithFields(logrus.Fields{
			"duration": duration,
			"uri":      r.RequestURI,
			"method":   r.Method,
		}).Info("REQUEST")

		logrus.WithFields(logrus.Fields{
			"size":   responseData.size,
			"status": responseData.status,
		}).Info("RESPONSE")
	})
}
