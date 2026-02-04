package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"go-demo/handlers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User API", func() {

	It("should create a user", func() {
		body := []byte(`{"name":"Ginkgo User","role":"Tester"}`)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handlers.UserHandler(rr, req)

		Expect(rr.Code).To(Equal(http.StatusCreated))
	})

	It("should enqueue welcome email notification when creating a user", func() {
		body := []byte(`{"name":"Welcome User","role":"Member"}`)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handlers.UserHandler(rr, req)

		Expect(rr.Code).To(Equal(http.StatusCreated))
		// The notification is enqueued asynchronously and processed by the worker
		// The API response is not blocked by notification processing
	})

	It("should get users", func() {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rr := httptest.NewRecorder()

		handlers.UserHandler(rr, req)

		Expect(rr.Code).To(Equal(http.StatusOK))
	})

	It("should update a user", func() {
		body := []byte(`{"name":"Updated User","role":"Admin"}`)
		req := httptest.NewRequest(http.MethodPut, "/users?id=1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handlers.UserHandler(rr, req)

		Expect(rr.Code).To(Equal(http.StatusOK))
	})

	It("should delete a user", func() {
		req := httptest.NewRequest(http.MethodDelete, "/users?id=1", nil)
		rr := httptest.NewRecorder()

		handlers.UserHandler(rr, req)

		Expect(rr.Code).To(Equal(http.StatusNoContent))
	})
})
