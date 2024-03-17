package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mehranzand/pulseup/api/middleware"
	"github.com/mehranzand/pulseup/utils"
	log "github.com/sirupsen/logrus"
)

// GetContainers
// @Summary Get list of containers
func (h *Handler) GetContainers(c echo.Context) error {
	cc := c.(*middleware.DockerContext)

	if cc.Client != nil {
		list, _ := cc.Client.ListContainers()

		for container := range list {
			c := list[container]
			log.Info(c.Name)
		}
	}

	if cc.Client == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	return c.String(http.StatusOK, "/api/:host/containers")
}
