package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-demo/models"
	uuidpkg "go-demo/pkg/uuid"
	"go-demo/pkg/validator"
	"go-demo/repositories"
)

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		products, err := repositories.GetProducts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(products)

	case http.MethodPost:
		var product models.Product
		json.NewDecoder(r.Body).Decode(&product)

		// assign UUID for the new product
		product.UUID = uuidpkg.New()

		if err := validator.Validate.Struct(product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := repositories.CreateProduct(product); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodPut:
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid product id", http.StatusBadRequest)
			return
		}

		var product models.Product
		json.NewDecoder(r.Body).Decode(&product)

		if err := validator.Validate.Struct(product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := repositories.UpdateProduct(id, product); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid product id", http.StatusBadRequest)
			return
		}

		if err := repositories.DeleteProduct(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
