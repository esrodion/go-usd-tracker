package repository

import (
	"context"
	"database/sql"
	"fmt"

	"go-usdtrub/internal/config"
	"go-usdtrub/internal/models"

	postgres "go-usdtrub/internal/repository/db/postgres"
)

type Repository struct {
	db  *sql.DB
	cfg *config.Config
}

func NewRepository(db *sql.DB, cfg *config.Config) (*Repository, error) {
	var err error

	repo := &Repository{
		db:  db,
		cfg: cfg,
	}

	if repo.cfg == nil {
		repo.cfg, err = config.NewConfig()
		if err != nil {
			return nil, fmt.Errorf("repository.NewRepository: could not load postgres config: %w", err)
		}
	}

	if repo.db == nil {
		repo.db, err = postgres.NewPostgresDB(repo.cfg.PostgresConfig)
		if err != nil {
			return nil, fmt.Errorf("repository.NewRepository: could not open postgres db: %w", err)
		}
	}

	if repo.cfg.AutoMigrateUp == "true" {
		err = postgres.MigrateUp(repo.db, repo.cfg.MigrationsURL)
		if err != nil {
			return nil, fmt.Errorf("repository.NewRepository: auto migration failed: %w", err)
		}
	}

	return repo, nil
}

func (repo *Repository) Close() error {
	return repo.db.Close()
}

func (repo *Repository) GetRates(ctx context.Context) (models.CurrenceyRate, error) {
	row := repo.db.QueryRowContext(ctx, `
		SELECT created_at, ask, bid 
		FROM usdtrub 
		ORDER BY created_at DESC
		LIMIT 1`)

	result := models.CurrenceyRate{}
	err := row.Scan(&result.Timestamp, &result.Ask, &result.Bid)
	if err != nil {
		return result, fmt.Errorf("repository.Repository: could not fetch rates: %w", err)
	}

	return result, nil
}

func (repo *Repository) AddRates(ctx context.Context, rate models.CurrenceyRate) error {
	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO usdtrub (ask, bid) 
		VALUES ($1, $2) 
		ON CONFLICT (created_at) DO 
			UPDATE SET (ask, bid) = ($1, $2)`, rate.Ask, rate.Bid)

	if err != nil {
		return fmt.Errorf("repository.Repository: could not add rates: %w", err)
	}

	return nil
}

func (repo *Repository) SetRates(ctx context.Context, rate models.CurrenceyRate) error {
	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO usdtrub (created_at, ask, bid) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (created_at) DO 
			UPDATE SET (ask, bid) = ($2, $3)`, rate.Timestamp, rate.Ask, rate.Bid)

	if err != nil {
		return fmt.Errorf("repository.Repository: could not set rates: %w", err)
	}

	return nil
}
