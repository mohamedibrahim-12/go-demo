package main

import (
	"log"
	"net/http"

	"go-demo/config"
	"go-demo/database"
	"go-demo/handlers"
	"go-demo/middlewares"
)

func main() {
	config.LoadEnv()
	database.Connect()

	mux := http.NewServeMux()

	mux.HandleFunc("/users", handlers.UserHandler)
	mux.HandleFunc("/products", handlers.ProductHandler)

	log.Println("🚀 Server running on :8080")
	http.ListenAndServe(":8080", middlewares.LoggingMiddleware(mux))
}
