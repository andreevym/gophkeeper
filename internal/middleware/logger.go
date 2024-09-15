package middleware

import (
	"io"
	"net/http"
	"time"

	"github.com/andreevym/gophkeeper/pkg/logger"
	"go.uber.org/zap"
)

// WithRequestLoggerMiddleware returns a middleware handler that logs details about incoming HTTP requests and their responses.
// The middleware logs the HTTP method, request URI, response status, and duration of the request.
// It also logs the response body length if available.
//
// Parameters:
//   - h: The next http.Handler to execute in the middleware chain.
//
// Returns:
//   - http.Handler: A new HTTP handler that includes the request logging functionality.
func WithRequestLoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response recorder to capture the response details.
		rec := &responseRecorder{ResponseWriter: w}
		h.ServeHTTP(rec, r)

		end := time.Now()

		// Log request details.
		logger.Logger().Debug(
			"request",
			zap.String("method", r.Method),
			zap.String("URI", r.RequestURI),
			zap.Duration("duration", end.Sub(start)),
		)

		// Log response details.
		if r.Body != nil {
			defer func() {
				err := r.Body.Close()
				if err != nil {
					logger.Logger().Error("error closing request body", zap.Error(err))
				}
			}()
			bytes, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Logger().Error("error reading request body", zap.Error(err))
				return
			}
			logger.Logger().Debug(
				"response",
				zap.Int("status", rec.statusCode),
				zap.Int("body_length", len(bytes)),
			)
		}
	})
}

// responseRecorder is a custom implementation of http.ResponseWriter that captures the status code.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the HTTP status code.
func (rec *responseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// Write is used to capture the length of the response body.
func (rec *responseRecorder) Write(b []byte) (int, error) {
	return rec.ResponseWriter.Write(b)
}
