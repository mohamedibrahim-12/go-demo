package middlewares

import (
	"net/http"
	"time"

	"go-demo/pkg/logger"
)

// LoggingMiddleware logs request method, path and execution time
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		// call next handler
		next.ServeHTTP(w, r)

		duration := time.Since(start)

		logger.Log.Info().Str("method", r.Method).Str("path", r.URL.Path).Dur("latency", duration).Msg("request")
	})
}
