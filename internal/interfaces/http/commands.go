package http

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mikaelchan/hamster/internal/infrastructure/serializer"
	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/logger"
	"github.com/mikaelchan/hamster/pkg/messaging"
	"github.com/mikaelchan/hamster/pkg/serializer"
)

func WrapCommands(router *gin.Engine, cb messaging.CommandBus) {
	router.POST("/commands/:contract", wrap(cb))
}

func wrap(cb messaging.CommandBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		contract := c.Param("contract")
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		command, err := serializer.DeserializeNew(domain.Type(contract), body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		logger.Infof("dispatching command: %s", contract)
		cb.Dispatch(c.Request.Context(), command.(domain.Command))
		c.Status(http.StatusAccepted)
	}
}
