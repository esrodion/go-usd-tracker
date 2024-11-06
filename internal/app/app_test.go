package app

import (
	"go-usdtrub/internal/repository"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAppStartup(t *testing.T) {
	var err error

	err = os.Setenv("AUTO_MIGRATE_UP", "false")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("AUTO_MIGRATE_DOWN", "false")
	if err != nil {
		t.Fatal(err)
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectClose()

	repo, err := repository.NewRepository(repository.WithDB(db, &MockMigrator{}))
	if err != nil {
		t.Fatal(err)
	}

	app := NewApp(WithRepo(repo))

	go func() {
		err := app.Run()
		if err != nil {
			t.Log(err) // handled at app.Stop()
		}
	}()

	err = app.Stop()
	if err != nil {
		t.Fatal(err)
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
