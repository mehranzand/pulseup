package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mehranzand/pulseup/api/middleware"
	"github.com/mehranzand/pulseup/internal/docker"
)

// GetContainers
// @Summary Get list of containers
func (h *Handler) GetContainers(c echo.Context) error {
	cc := c.(*middleware.DockerContext)

	containers, err := cc.Client.ListContainers()
	if err != nil {
		http.Error(c.Response().Writer, "Continer not found!", http.StatusInternalServerError)
	}

	enc := json.NewEncoder(c.Response())
	for _, l := range containers {
		if err := enc.Encode(l.Name); err != nil {
			return err
		}
		fmt.Fprintf(c.Response().Writer, "| ")
		c.Response().Flush()
	}

	return nil
}

func (h *Handler) StreamLogs(c echo.Context) error {
	var stdTypes docker.StdType
	stdTypes |= docker.STDOUT

	cc := c.(*middleware.DockerContext)

	_, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		http.Error(c.Response().Writer, "Streaming unsupported!", http.StatusInternalServerError)
	}

	_, err := cc.Client.ContainerLogs(c.Request().Context(), "329c358a5e81", "", stdTypes)
	if err != nil {
		http.Error(c.Response().Writer, "Continer not found!", http.StatusInternalServerError)
	}

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-transform")
	c.Response().Header().Add("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no")

	return nil
}
