package app

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/mikaelchan/hamster/pkg/logger"
)

type App interface {
	Run(ctx context.Context) error
}

type app struct {
	shutdownTimeout time.Duration
}

// Close implements App.
func (a *app) stop(ctx context.Context) error {
	// ctx, cancel := context.WithTimeout(ctx, a.shutdownTimeout)
	// defer cancel()
	logger.Infof("app has been shutdown")

	logger.Sync()
	return nil
}

// Start implements App.
func (a *app) Run(ctx context.Context) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig
	logger.Infof("received interrupt signal, shutting down")
	go func() {
		<-sig
		logger.Fatalf("received interrupt signal again, force shutting down")
	}()
	return a.stop(ctx)
}

func NewApp(shutdownTimeout time.Duration) App {
	return &app{shutdownTimeout: shutdownTimeout}
}
