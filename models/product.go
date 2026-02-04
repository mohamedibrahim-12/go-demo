package models

type Product struct {
	ID    int     `json:"id"`
	UUID  string  `json:"uuid,omitempty"`
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price" validate:"required,gt=0"`
}
