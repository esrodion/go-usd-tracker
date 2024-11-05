package app

import (
	"context"
	"errors"
	"go-usdtrub/internal/config"
	"go-usdtrub/internal/controller"
	"go-usdtrub/internal/repository"
	"go-usdtrub/internal/router"
	"go-usdtrub/internal/service"
	"go-usdtrub/internal/traces"
	"go-usdtrub/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	StopSig chan os.Signal
	Done    chan os.Signal
	err     error
	repo    *repository.Repository
}

func WithRepo(repo *repository.Repository) func(app *App) {
	return func(app *App) {
		app.repo = repo
	}
}

func NewApp(opts ...func(app *App)) *App {
	app := &App{
		Done:    make(chan os.Signal),
		StopSig: make(chan os.Signal, 2),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (app *App) Run() error {
	defer close(app.Done)

	ctx := context.Background()

	// init config and logger

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	config.ReadFlags(cfg)
	if !config.ReadFlags(cfg) {
		config.PrintCLIUsage()
		return nil
	}

	log := logger.Logger().Sugar().Named("App")

	// init tracer

	err = traces.Init(ctx, "go-usd-tracker")
	if err != nil {
		return err
	}

	// init repository and service layers

	if app.repo == nil {
		repo, err := repository.NewRepository(repository.WithCfg(cfg))
		if err != nil {
			return err
		}
		app.repo = repo
	}

	serv := service.NewService(app.repo)

	// gRPC controller

	cont, err := controller.NewGrpcController(serv, cfg.GrpcAddress)
	if err != nil {
		return err
	}

	// HTTP controller for healthcheck

	server := http.Server{
		Addr:    cfg.HttpAddress,
		Handler: router.NewRouter(controller.NewHttpController(serv)),
	}
	go func() {
		log.Info("HTTP server started at ", cfg.HttpAddress)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.DPanic(err.Error())
		}
	}()

	// Graceful shutdown

	signal.Notify(app.StopSig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-app.StopSig
	log.Info("Received signal: ", sig)

	log.Info("Closing grpc server...")
	cont.Stop()

	log.Info("Closing http server...")
	errHttp := server.Shutdown(ctx)

	log.Info("Closing database connection...")
	errRepo := app.repo.Close()

	app.err = errors.Join(errHttp, errRepo)
	return app.err
}

func (app *App) Stop() error {
	app.StopSig <- os.Interrupt
	<-app.Done

	return app.err
}

func (app *App) Err() error {
	return app.err
}
