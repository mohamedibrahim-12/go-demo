package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"go-demo/handlers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Product API", func() {

	It("should create a product", func() {
		body := []byte(`{"name":"Ginkgo Product","price":199.99}`)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handlers.ProductHandler(rr, req)

		Expect(rr.Code).To(Equal(http.StatusCreated))
	})

	It("should get products", func() {
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		rr := httptest.NewRecorder()

		handlers.ProductHandler(rr, req)

		Expect(rr.Code).To(Equal(http.StatusOK))
	})

	It("should update a product", func() {
		body := []byte(`{"name":"Updated Product","price":299.99}`)
		req := httptest.NewRequest(http.MethodPut, "/products?id=1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handlers.ProductHandler(rr, req)

		Expect(rr.Code).To(Equal(http.StatusOK))
	})

	It("should delete a product", func() {
		req := httptest.NewRequest(http.MethodDelete, "/products?id=1", nil)
		rr := httptest.NewRecorder()

		handlers.ProductHandler(rr, req)

		Expect(rr.Code).To(Equal(http.StatusNoContent))
	})
})
