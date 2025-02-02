package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mikaelchan/hamster/pkg/app"
)

type adapter struct {
	srv *http.Server
}

func (a *adapter) Run(ctx context.Context) error {
	if err := a.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (a *adapter) Stop(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}

func NewAdapter(addr string, router *gin.Engine) app.Adapter {
	return &adapter{
		srv: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}
