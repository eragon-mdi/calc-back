package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/eragon-mdi/calc-back/internal/common/configs"
	"github.com/go-faster/errors"
)

const (
	ErrListenAndServe = "server failed to start"
	ErrShutdown       = "server shutdown failed"
)

type Server interface {
	Start() error
	GracefulShutdown() error
	WaitingShutdownSignal()
}

type server struct {
	*http.Server
}

func New(rHandler http.Handler, cfg *configs.Server) Server {
	return server{
		Server: &http.Server{
			Addr:              fmt.Sprintf("%s:%s", cfg.Address, cfg.Port),
			Handler:           rHandler,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			IdleTimeout:       cfg.IdleTimeout,
		},
	}
}

func (s server) Start() error {
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return errors.Wrap(err, ErrListenAndServe)
	}

	return nil
}

func (s server) GracefulShutdown() error {
	if err := s.Shutdown(context.Background()); err != nil {
		return errors.Wrap(err, ErrShutdown)
	}

	return nil
}

func (s server) WaitingShutdownSignal() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	<-done
}
