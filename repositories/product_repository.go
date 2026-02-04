package repositories

import (
	"go-demo/database"
	"go-demo/models"
)

func GetProducts() ([]models.Product, error) {
	rows, err := database.DB.Query(
		"SELECT id, uuid, name, price FROM products",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		rows.Scan(&p.ID, &p.UUID, &p.Name, &p.Price)
		products = append(products, p)
	}
	return products, nil
}

func CreateProduct(p models.Product) error {
	_, err := database.DB.Exec(
		"INSERT INTO products (uuid, name, price) VALUES ($1, $2, $3)",
		p.UUID,
		p.Name,
		p.Price,
	)
	return err
}

func UpdateProduct(id int, p models.Product) error {
	_, err := database.DB.Exec(
		"UPDATE products SET name=$1, price=$2 WHERE id=$3",
		p.Name,
		p.Price,
		id,
	)
	return err
}

func DeleteProduct(id int) error {
	_, err := database.DB.Exec(
		"DELETE FROM products WHERE id=$1",
		id,
	)
	return err
}
