package repository

import (
	"context"
	"go-usdtrub/internal/config"
	"go-usdtrub/internal/metrics"
	"go-usdtrub/internal/models"
	"go-usdtrub/pkg/logger"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type testCase struct {
	time     time.Time
	ask, bid float64
}

var testCases []testCase

func TestRepository(t *testing.T) {
	start := time.Now()

	testCases = []testCase{
		{start, rand.Float64(), rand.Float64()},
		{start, rand.Float64(), rand.Float64()},
		{start, rand.Float64(), rand.Float64()},
		{start, rand.Float64(), rand.Float64()},
		{start, rand.Float64(), rand.Float64()},
	}

	repo := openRepo(t)
	defer repo.Close()

	ctx := context.WithValue(context.Background(), metrics.WrapperKey, metrics.ZeroHandler)

	for _, tcase := range testCases {
		err := repo.AddRates(ctx, models.CurrenceyRate{Ask: tcase.ask, Bid: tcase.bid})
		if err != nil {
			t.Fatal(err)
		}

		rate, err := repo.GetRates(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if rate.Ask != tcase.ask || rate.Bid != tcase.bid {
			t.Fatalf("Expected ask and bid to be %.2f and %.2f, got %.2f, %.2f", tcase.ask, tcase.bid, rate.Ask, rate.Bid)
		}

		rate.Ask, rate.Bid = tcase.ask+1, tcase.bid+1
		err = repo.SetRates(ctx, rate)
		if err != nil {
			t.Fatal(err)
		}

		rate, err = repo.GetRates(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if rate.Ask != tcase.ask+1 || rate.Bid != tcase.bid+1 {
			t.Fatalf("Expected ask and bid to be %.2f and %.2f, got %.2f, %.2f", tcase.ask+1, tcase.bid+1, rate.Ask, rate.Bid)
		}
	}
}

func openRepo(t *testing.T) *Repository {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	cfg.PostgresConfig.AutoMigrateUp = "true"
	cfg.PostgresConfig.AutoMigrateDown = "true"
	cfg.PostgresConfig.Conn = os.Getenv("INTEGRATION_TESTS_DB_CONN")

	repo, err := NewRepository(WithCfg(cfg))
	if err != nil {
		logger.Logger().Sugar().Named("repository_test").Debug("Could not connect DB on ", cfg.PostgresConfig.Conn, " falling back to mock DB, only unit tests will be performed")

		DB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		prepDBMock(mock)

		repo, err = NewRepository(WithCfg(cfg), WithDB(DB, &MockMigrator{}))
		if err != nil {
			t.Fatal(err)
		}
	}

	return repo
}

//// Mock DB

func prepDBMock(mock sqlmock.Sqlmock) {
	for _, tcase := range testCases {
		mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"created_at", "ask", "bid"}).AddRow(time.Now(), tcase.ask, tcase.bid))
		mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"created_at", "ask", "bid"}).AddRow(time.Now(), tcase.ask+1, tcase.bid+1))
	}
}

//// Mock Migrator

type MockMigrator struct{}

func (m *MockMigrator) Up() error {
	return nil
}

func (m *MockMigrator) Down() error {
	return nil
}
