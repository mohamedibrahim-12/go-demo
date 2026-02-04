package models

type User struct {
	ID   int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UUID string `json:"uuid,omitempty" gorm:"type:uuid;default:gen_random_uuid();column:uuid"`
	Name string `json:"name" validate:"required" gorm:"column:name;not null"`
	Role string `json:"role" validate:"required" gorm:"column:role;not null"`
}

func (User) TableName() string {
	return "users"
}
