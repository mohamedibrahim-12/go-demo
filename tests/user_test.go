package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-demo/handlers"
)

func TestUserCRUD(t *testing.T) {

	// ---------- CREATE ----------
	createBody := []byte(`{"name":"Test User","role":"Tester"}`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handlers.UserHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("CREATE user failed, expected 201 got %d", rr.Code)
	}

	// ---------- READ ----------
	req = httptest.NewRequest(http.MethodGet, "/users", nil)
	rr = httptest.NewRecorder()

	handlers.UserHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("GET users failed, expected 200 got %d", rr.Code)
	}

	// ---------- UPDATE ----------
	updateBody := []byte(`{"name":"Updated User","role":"Manager"}`)
	req = httptest.NewRequest(http.MethodPut, "/users?id=1", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handlers.UserHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("UPDATE user failed, expected 200 got %d", rr.Code)
	}

	// ---------- DELETE ----------
	req = httptest.NewRequest(http.MethodDelete, "/users?id=1", nil)
	rr = httptest.NewRecorder()

	handlers.UserHandler(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("DELETE user failed, expected 204 got %d", rr.Code)
	}
}
