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
		log.Print("ERROR: error reading config")
		return nil, err
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		conf[config.DB_USER], conf[config.DB_PASS], conf[config.DB_HOST],
		conf[config.DB_PORT], conf[config.DB_NAME])
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Print("ERROR: Failed to open DB:", err)
		return nil, err
	}

	// ping to ensure DB is reachable
	if err := db.Ping(); err != nil {
		log.Print("ERROR: DB not reachable:", err)
		return nil, err
	}
	log.Printf("'%s' database connected", conf[config.DB_NAME])
	return db, nil
}
