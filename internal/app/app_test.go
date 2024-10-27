package app

import (
	"os"
	"testing"
)

func TestApp(t *testing.T) {
	var err error

	err = os.Setenv("POSTGRES_CONN", "postgres://postgres:yourpassword@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("MIGRATIONS_URL", "file://D:/Kata/Repo/goUSDtracker/go-usd-tracker/internal/repository/db/migrations")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("AUTO_MIGRATE_UP", "false")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("AUTO_MIGRATE_DOWN", "false")
	if err != nil {
		t.Fatal(err)
	}

	app := NewApp()

	go func() {
		err = app.Run()
	}()

	app.StopSig <- os.Interrupt
	<-app.Done

	if err != nil {
		t.Fatal(err)
	}
}
