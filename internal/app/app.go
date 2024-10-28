package app

import (
	"context"
	"go-usdtrub/internal/config"
	"go-usdtrub/internal/controller"
	"go-usdtrub/internal/repository"
	"go-usdtrub/internal/router"
	"go-usdtrub/internal/service"
	"go-usdtrub/pkg/logger"
	"net/http"
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

	// prepare logger and config

	log := logger.Logger().Sugar().Named("App")

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	// init repository and service layers

	repo, err := repository.NewRepository(cfg)
	if err != nil {
		return err
	}

	serv := service.NewService(repo)

	// gRPC controller

	cont, err := controller.NewGrpcController(serv, cfg.GrpcAddress)
	if err != nil {
		return err
	}

	// HTTP controller for healthcheck

	server := http.Server{
		Addr:    cfg.HttpAddress,
		Handler: router.NewRouter(controller.NewHttpController()),
	}
	go func() {
		log.Info("HTTP server started at ", cfg.HttpAddress)
		err := server.ListenAndServe()
		if err != nil {
			log.DPanic(err.Error())
		}
	}()

	// Graceful shutdown

	signal.Notify(app.StopSig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-app.StopSig
	log.Info("Received signal: ", sig)

	cont.Stop()
	server.Shutdown(context.Background())
	return repo.Close()
}
