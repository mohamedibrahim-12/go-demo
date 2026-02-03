package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-demo/middlewares"
)

// dummy handler
func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestLoggingMiddleware(t *testing.T) {
	handler := middlewares.LoggingMiddleware(http.HandlerFunc(testHandler))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}
