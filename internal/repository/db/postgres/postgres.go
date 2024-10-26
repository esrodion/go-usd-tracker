package postgres

import (
	"database/sql"
	"go-usdtrub/internal/config"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg config.PostgresConfig) (*sql.DB, error) {
	// "TODO: Zap logging"
	log.Println("Connecting db on: ", cfg.Conn)
	db, err := sql.Open("postgres", cfg.Conn)

	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
