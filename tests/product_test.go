package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-demo/handlers"
)

func TestProductCRUD(t *testing.T) {

	// ---------- CREATE ----------
	createBody := []byte(`{"name":"Test Product","price":100.50}`)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handlers.ProductHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("CREATE product failed, expected 201 got %d", rr.Code)
	}

	// ---------- READ ----------
	req = httptest.NewRequest(http.MethodGet, "/products", nil)
	rr = httptest.NewRecorder()

	handlers.ProductHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("GET products failed, expected 200 got %d", rr.Code)
	}

	// ---------- UPDATE ----------
	updateBody := []byte(`{"name":"Updated Product","price":250.00}`)
	req = httptest.NewRequest(http.MethodPut, "/products?id=1", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handlers.ProductHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("UPDATE product failed, expected 200 got %d", rr.Code)
	}

	// ---------- DELETE ----------
	req = httptest.NewRequest(http.MethodDelete, "/products?id=1", nil)
	rr = httptest.NewRecorder()

	handlers.ProductHandler(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("DELETE product failed, expected 204 got %d", rr.Code)
	}
}
