package repositories

import (
	"go-demo/database"
	"go-demo/models"
)

func GetUsers() ([]models.User, error) {
	rows, err := database.DB.Query(
		"SELECT id, uuid, name, role FROM users",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.UUID, &u.Name, &u.Role)
		users = append(users, u)
	}
	return users, nil
}

func CreateUser(u models.User) error {
	_, err := database.DB.Exec(
		"INSERT INTO users (uuid, name, role) VALUES ($1, $2, $3)",
		u.UUID,
		u.Name,
		u.Role,
	)
	return err
}

func UpdateUser(id int, u models.User) error {
	_, err := database.DB.Exec(
		"UPDATE users SET name=$1, role=$2 WHERE id=$3",
		u.Name,
		u.Role,
		id,
	)
	return err
}

func DeleteUser(id int) error {
	_, err := database.DB.Exec(
		"DELETE FROM users WHERE id=$1",
		id,
	)
	return err
}
