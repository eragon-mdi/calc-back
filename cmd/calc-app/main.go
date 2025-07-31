package main

import (
	"log"

	"github.com/eragon-mdi/calc-back/internal/common/api"
	"github.com/eragon-mdi/calc-back/internal/common/configs"
	"github.com/eragon-mdi/calc-back/internal/common/logger"
	"github.com/eragon-mdi/calc-back/internal/common/middlewares"
	"github.com/eragon-mdi/calc-back/internal/common/storage"
	"github.com/eragon-mdi/calc-back/internal/repository"
	"github.com/eragon-mdi/calc-back/internal/service"
	"github.com/eragon-mdi/calc-back/internal/transport"
	"github.com/eragon-mdi/calc-back/pkg/server"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := configs.Get()

	l, err := logger.New(cfg.Logger)
	if err != nil {
		log.Fatal(err)
	}

	store, err := storage.Conn(&cfg.Storage, storage.ConnTimeoutDefault)
	if err != nil {
		log.Fatal(err)
	}

	r := repository.New(store)
	s := service.New(r)
	t := transport.New(s, l)

	e := echo.New()
	m := middlewares.New(&cfg.Middlerware)
	api.RegisterRoutes(e, t, m)

	srv := server.New(e, &cfg.Server)
	go func() {
		l.Infow("starting server", "addr", cfg.Server.Address)
		if err := srv.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	srv.WaitingShutdownSignal()

	if err := srv.GracefulShutdown(); err != nil {
		l.Errorw("error during server shutdown", "error", err)
	}
	if err := store.GracefulShutdown(); err != nil {
		l.Errorw("error disconnect store", "error", err)
	}
}
