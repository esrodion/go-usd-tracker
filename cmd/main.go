package main

import (
	"go-usdtrub/internal/app"
	"log"

	"github.com/joho/godotenv"
)

func main() {
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
