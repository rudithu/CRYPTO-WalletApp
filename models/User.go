package models

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func GetUserById(db *sql.DB, id int64) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, name, created_at FROM users where id=$1", id).Scan(&user.ID, &user.Name, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
