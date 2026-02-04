package main

import (
	"net/http"

	"go-demo/config"
	"go-demo/database"
	"go-demo/handlers"
	"go-demo/middlewares"
	"go-demo/pkg/logger"
	"go-demo/pkg/validator"
)

func main() {
	config.LoadEnv()
	logger.Init()
	// initialize validator before connecting/handling requests
	validator.Init()

	database.Connect()

	mux := http.NewServeMux()

	mux.HandleFunc("/users", handlers.UserHandler)
	mux.HandleFunc("/products", handlers.ProductHandler)

	logger.Log.Info().Msg("🚀 Server running on :8080")
	// Rate limit outer, then logging, then mux
	handler := middlewares.RateLimitMiddleware(middlewares.LoggingMiddleware(mux))
	http.ListenAndServe(":8080", handler)
}
