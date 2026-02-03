package tests

import (
	"os"
	"testing"

	"go-demo/config"
	"go-demo/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAPIs(t *testing.T) {
	os.Setenv("APP_ENV", "test")

	config.LoadEnv()
	database.Connect()

	RegisterFailHandler(Fail)
	RunSpecs(t, "API Integration Test Suite")
}
