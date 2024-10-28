package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go-usdtrub/internal/config"
	"go-usdtrub/internal/metrics"
	"go-usdtrub/internal/models"

	database "go-usdtrub/internal/repository/db"

	"github.com/golang-migrate/migrate/v4"
)

type Repository struct {
	db       *sql.DB
	migrator *migrate.Migrate
	cfg      *config.Config
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	var err error

	repo := &Repository{
		cfg: cfg,
	}

	if repo.cfg == nil {
		repo.cfg, err = config.NewConfig()
		if err != nil {
			return nil, fmt.Errorf("repository.NewRepository: could not load postgres config: %w", err)
		}
	}

	db, m, err := database.NewPostgresDB(repo.cfg.PostgresConfig)
	if err != nil {
		return nil, fmt.Errorf("repository.NewRepository: could not open postgres db: %w", err)
	}
	repo.db, repo.migrator = db, m

	if repo.cfg.AutoMigrateUp == "true" {
		err = repo.migrator.Up()
		if err != nil {
			repo.db.Close()
			return nil, fmt.Errorf("repository.NewRepository: auto migration failed: %w", err)
		}
	}

	return repo, nil
}

func (repo *Repository) Close() error {
	if repo.cfg.AutoMigrateDown == "true" {
		err := repo.migrator.Down()
		if err != nil {
			return errors.Join(fmt.Errorf("repository.NewRepository: auto migration failed: %w", err), repo.db.Close())
		}
	}

	return repo.db.Close()
}

func (repo *Repository) GetRates(ctx context.Context) (models.CurrenceyRate, error) {
	defer metrics.DBAccessDuration(ctx, time.Now())

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
	defer metrics.DBAccessDuration(ctx, time.Now())

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
	defer metrics.DBAccessDuration(ctx, time.Now())

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
