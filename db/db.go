package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rudithu/CRYPTO-WalletApp/config"
)

func Connnect() (*sql.DB, error) {

	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal("error reading config")
		return nil, err
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		conf[config.DB_USER], conf[config.DB_PASS], conf[config.DB_HOST],
		conf[config.DB_PORT], conf[config.DB_NAME])
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	// Optional: ping to ensure DB is reachable
	if err := db.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}
	fmt.Println("db connected")
	return db, nil
}
