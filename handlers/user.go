package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-demo/models"
	uuidpkg "go-demo/pkg/uuid"
	"go-demo/pkg/validator"
	"go-demo/repositories"
	"go-demo/worker"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		users, err := repositories.GetUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)

		// fire‑and‑forget audit log; API response is not blocked
		worker.Publish(worker.NewEvent(
			"READ",
			"user",
			0,
			"listed users",
		))

	case http.MethodPost:
		var user models.User
		json.NewDecoder(r.Body).Decode(&user)

		// assign a UUID for this new user
		user.UUID = uuidpkg.New()

		if err := validator.Validate.Struct(user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := repositories.CreateUser(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

		worker.Publish(worker.NewEvent(
			"CREATE",
			"user",
			0,
			"created user",
		))

		// enqueue welcome email notification asynchronously
		// API response is not blocked by notification processing
		// enqueue welcome email notification asynchronously using Outbox Pattern
		// API response is not blocked by notification processing
		// Payload: JSON
		recipient := user.Name + "@example.com"
		payloadMap := map[string]string{
			"recipient": recipient,
			"message":   "Welcome to our platform, " + user.Name + "!",
		}
		payloadBytes, _ := json.Marshal(payloadMap)

		repositories.CreateNotificationOutbox("WELCOME_EMAIL", string(payloadBytes))

	case http.MethodPut:
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		var user models.User
		json.NewDecoder(r.Body).Decode(&user)

		if err := validator.Validate.Struct(user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := repositories.UpdateUser(id, user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		worker.Publish(worker.NewEvent(
			"UPDATE",
			"user",
			id,
			"updated user",
		))

	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		if err := repositories.DeleteUser(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

		worker.Publish(worker.NewEvent(
			"DELETE",
			"user",
			id,
			"deleted user",
		))
	}
}
