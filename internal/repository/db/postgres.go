package postgres

import (
	"database/sql"
	"fmt"
	"go-usdtrub/internal/config"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg config.PostgresConfig) (*sql.DB, *migrate.Migrate, error) {
	// "TODO: Zap logging"
	log.Println("Connecting db on: ", cfg.Conn)
	db, err := sql.Open("postgres", cfg.Conn)

	if err != nil {
		return nil, nil, fmt.Errorf("postgres.NewPostgresDB: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("postgres.NewPostgresDB: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(cfg.MigrationsURL, "postgres", driver)
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("postgres.NewPostgresDB: %w", err)
	}

	return db, m, nil
}
