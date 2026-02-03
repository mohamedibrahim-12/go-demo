package tests

import (
	"log"
	"os"
	"testing"

	"go-demo/config"
	"go-demo/database"
)

func TestMain(m *testing.M) {
	// Ensure test env is loaded before connecting to DB
	os.Setenv("APP_ENV", "test")
	config.LoadEnv()
	database.Connect()

	// Run tests
	code := m.Run()

	log.Println("Tests finished")

	os.Exit(code)
}
