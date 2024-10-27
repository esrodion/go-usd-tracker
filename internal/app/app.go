package app

import (
	"go-usdtrub/internal/config"
	"go-usdtrub/internal/controller"
	"go-usdtrub/internal/repository"
	"go-usdtrub/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	StopSig chan os.Signal
	Done    chan os.Signal
}

func NewApp() *App {
	return &App{
		Done:    make(chan os.Signal),
		StopSig: make(chan os.Signal, 2),
	}
}

func (app *App) Run() error {
	defer close(app.Done)

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	repo, err := repository.NewRepository(cfg)
	if err != nil {
		return err
	}

	serv := service.NewService(repo)
	cont, err := controller.NewGrpcController(serv, cfg.ServerAddress)
	if err != nil {
		return err
	}

	signal.Notify(app.StopSig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-app.StopSig
	log.Printf("Received signal: %s\n", sig)

	cont.Stop()
	return repo.Close()
}
