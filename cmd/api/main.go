package main

import (
	"net/http"

	"go-demo/config"
	"go-demo/database"
	apphandlers "go-demo/handlers"
	"go-demo/middlewares"
	"go-demo/pkg/logger"
	"go-demo/pkg/validator"

	ghandlers "github.com/gorilla/handlers"
)

func main() {
	config.LoadEnv()
	logger.Init()
	// initialize validator before connecting/handling requests
	validator.Init()

	database.Connect()

	mux := http.NewServeMux()

	mux.HandleFunc("/users", apphandlers.UserHandler)
	mux.HandleFunc("/products", apphandlers.ProductHandler)

	// Build handler chain:
	// 1) base mux
	// 2) logging middleware
	// 3) rate limiting
	// 4) compression
	// 5) CORS
	// 6) recovery (outermost)
	handler := middlewares.LoggingMiddleware(mux)
	handler = middlewares.RateLimitMiddleware(handler)
	handler = ghandlers.CompressHandler(handler)
	handler = ghandlers.CORS(ghandlers.AllowedOrigins([]string{"*"}))(handler)
	handler = ghandlers.RecoveryHandler(ghandlers.PrintRecoveryStack(true))(handler)

	logger.Log.Info().Msg("🚀 Server running on :8080")
	http.ListenAndServe(":8080", handler)
}
