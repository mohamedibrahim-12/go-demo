package models

import "time"

type Product struct {
	ID        int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UUID      string    `json:"uuid,omitempty" gorm:"type:uuid;default:gen_random_uuid();column:uuid"`
	Name      string    `json:"name" validate:"required" gorm:"column:name;not null"`
	Price     float64   `json:"price" validate:"required,gt=0" gorm:"column:price;not null"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"column:created_at;autoCreateTime"`
}

func (Product) TableName() string {
	return "products"
}
