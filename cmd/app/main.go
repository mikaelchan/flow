package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mikaelchan/hamster/internal/application/library"
	"github.com/mikaelchan/hamster/internal/interfaces/http"
	"github.com/mikaelchan/hamster/pkg/app"
	"github.com/mikaelchan/hamster/pkg/env"
	"github.com/mikaelchan/hamster/pkg/logger"
)

func main() {
	isRelease := env.IsRelease()
	ctx := context.Background()
	config, err := LoadConfig("hamster.toml")
	if err != nil {
		panic(err)
	}
	service := library.NewService(ctx, config.Service)
	defer func() {
		if err := service.Close(); err != nil {
			panic(err)
		}
	}()

	if isRelease {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(logger.GinLogger())
	router.Use(logger.GinRecovery(true))
	// TODO: maybe we need a auth middleware
	http.WrapCommands(router, service.CommandBus)

	app := app.NewApp(time.Second*10, http.NewAdapter(fmt.Sprintf("%s:%d", config.App.Host, config.App.Port), router))
	app.Run(ctx)
}
