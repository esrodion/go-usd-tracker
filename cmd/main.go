package main

import (
	"go-usdtrub/internal/app"
	"go-usdtrub/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	log := logger.Logger().Sugar().Named("main")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	app := app.NewApp()
	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
