package middlewares

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs request method, path and execution time
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		// call next handler
		next.ServeHTTP(w, r)

		duration := time.Since(start)

		log.Printf(
			"[%s] %s took %v",
			r.Method,
			r.URL.Path,
			duration,
		)
	})
}
