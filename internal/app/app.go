package app

import (
	"context"
	"go-usdtrub/internal/config"
	"go-usdtrub/internal/controller"
	"go-usdtrub/internal/repository"
	"go-usdtrub/internal/router"
	"go-usdtrub/internal/service"
	"log"
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

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	repo, err := repository.NewRepository(cfg)
	if err != nil {
		return err
	}

	serv := service.NewService(repo)
	cont, err := controller.NewGrpcController(serv, cfg.GrpcAddress)
	if err != nil {
		return err
	}

	server := http.Server{
		Addr:    cfg.HttpAddress,
		Handler: router.NewRouter(controller.NewHttpController()),
	}
	go func() {
		// TODO: Zap logging
		log.Println("HTTP server started at " + cfg.HttpAddress)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	signal.Notify(app.StopSig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-app.StopSig
	log.Printf("Received signal: %s\n", sig)

	cont.Stop()
	server.Shutdown(context.Background())
	return repo.Close()
}
