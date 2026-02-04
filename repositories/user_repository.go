package repositories

import (
	"go-demo/database"
	"go-demo/models"
)

func GetUsers() ([]models.User, error) {
	var users []models.User
	if err := database.GormDB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func CreateUser(u models.User) error {
	return database.GormDB.Create(&u).Error
}

func UpdateUser(id int, u models.User) error {
	return database.GormDB.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{"name": u.Name, "role": u.Role}).Error
}

func DeleteUser(id int) error {
	return database.GormDB.Delete(&models.User{}, id).Error
}
