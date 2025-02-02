package app

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/mikaelchan/hamster/pkg/logger"
)

type Adapter interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}

type App interface {
	Run(ctx context.Context) error
}

type app struct {
	shutdownTimeout time.Duration
	adapters        []Adapter
}

// Close implements App.
func (a *app) stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, a.shutdownTimeout)
	defer cancel()
	logger.Infof("app has been shutdown")

	errCh := make(chan error, len(a.adapters))
	for _, adapter := range a.adapters {
		go func(adapter Adapter) {
			errCh <- adapter.Stop(ctx)
		}(adapter)
	}

	for range a.adapters {
		if err := <-errCh; err != nil {
			logger.Errorf("adapter failed to stop: %v", err)
		}
	}
	logger.Sync()
	return nil
}

// Run implements App.
func (a *app) Run(ctx context.Context) error {
	for _, adapter := range a.adapters {
		go func(adapter Adapter) {
			if err := adapter.Run(ctx); err != nil {
				logger.Fatalf("adapter %s failed to start: %v", adapter, err)
			}
		}(adapter)
	}

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

func NewApp(shutdownTimeout time.Duration, adapters ...Adapter) App {
	return &app{shutdownTimeout: shutdownTimeout, adapters: adapters}
}
