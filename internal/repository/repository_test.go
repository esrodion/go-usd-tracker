package repository

import (
	"context"
	"go-usdtrub/internal/config"
	"go-usdtrub/internal/models"
	"math/rand"
	"testing"
)

// Test DB connection string for integration tests
const testDBConn = "postgres://postgres:yourpassword@localhost:5432/postgres?sslmode=disable"

// DB migrations URL to access from test executable
const testMigrations = "file://D:/Kata/Repo/goUSDtracker/go-usd-tracker/internal/repository/db/migrations"

func TestRepository(t *testing.T) {
	ask, bid := rand.Float64(), rand.Float64()

	repo := openRepo(t)
	defer repo.Close()

	ctx := context.Background()

	err := repo.AddRates(ctx, models.CurrenceyRate{Ask: ask, Bid: bid})
	if err != nil {
		t.Fatal(err)
	}

	rate, err := repo.GetRates(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if rate.Ask != ask || rate.Bid != bid {
		t.Fatalf("Expected ask and bid to be %.2f and %.2f, got %.2f, %.2f", ask, bid, rate.Ask, rate.Bid)
	}

	ask, bid = rand.Float64(), rand.Float64()
	rate.Ask, rate.Bid = ask, bid
	err = repo.SetRates(ctx, rate)
	if err != nil {
		t.Fatal(err)
	}

	rate, err = repo.GetRates(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if rate.Ask != ask || rate.Bid != bid {
		t.Fatalf("Expected ask and bid to be %.2f and %.2f, got %.2f, %.2f", ask, bid, rate.Ask, rate.Bid)
	}
}

func openRepo(t *testing.T) *Repository {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	cfg.PostgresConfig.AutoMigrateUp = "true"
	cfg.PostgresConfig.AutoMigrateDown = "true"
	cfg.PostgresConfig.MigrationsURL = testMigrations
	cfg.PostgresConfig.Conn = testDBConn

	repo, err := NewRepository(cfg)
	if err != nil {
		t.Fatal(err)
	}

	return repo
}
