package models

type User struct {
	ID   int    `json:"id"`
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name" validate:"required"`
	Role string `json:"role" validate:"required"`
}
