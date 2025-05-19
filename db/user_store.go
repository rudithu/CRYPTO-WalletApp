package db

import (
	"database/sql"
	"errors"

	"github.com/rudithu/CRYPTO-WalletApp/models"
)

func GetUserById(db *sql.DB, id int64) (*models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, name, created_at FROM users where id=$1", id).Scan(&user.ID, &user.Name, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
