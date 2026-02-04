package repositories

import (
	"go-demo/database"
	"go-demo/models"
)

func GetProducts() ([]models.Product, error) {
	var products []models.Product
	if err := database.GormDB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func CreateProduct(p models.Product) error {
	return database.GormDB.Create(&p).Error
}

func UpdateProduct(id int, p models.Product) error {
	return database.GormDB.Model(&models.Product{}).Where("id = ?", id).Updates(map[string]interface{}{"name": p.Name, "price": p.Price}).Error
}

func DeleteProduct(id int) error {
	return database.GormDB.Delete(&models.Product{}, id).Error
}
