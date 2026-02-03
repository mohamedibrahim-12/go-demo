package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	root, _ := os.Getwd()

	// Move to project root if running from tests/
	if filepath.Base(root) == "tests" {
		root = filepath.Dir(root)
	}

	env := os.Getenv("APP_ENV")

	var envFile string
	if env == "test" {
		envFile = filepath.Join(root, ".env.test")
		log.Println("🧪 Loading test environment")
	} else {
		envFile = filepath.Join(root, ".env")
		log.Println("🚀 Loading production environment")
	}

	log.Println("Loading env file:", envFile)

	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("❌ Failed to load env file: %s", envFile)
	}
}
