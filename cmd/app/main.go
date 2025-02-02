package main

import (
	"context"
	"time"

	"github.com/mikaelchan/hamster/internal/application/library"
	"github.com/mikaelchan/hamster/pkg/app"
)

func main() {
	ctx := context.Background()
	config, err := library.LoadConfig("hamster.toml")
	if err != nil {
		panic(err)
	}
	service := library.NewService(ctx, config)
	defer func() {
		if err := service.Close(); err != nil {
			panic(err)
		}
	}()
	app := app.NewApp(time.Second * 10)
	app.Run(ctx)
}
