package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connnect() *sql.DB {
	dsn := "postgres://crypto_wallet:wallet@localhost:5432/wallet_db" // Replace with env or config
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	// Optional: ping to ensure DB is reachable
	if err := db.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	return db
}
