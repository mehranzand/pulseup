package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mehranzand/pulseup/api/middleware"
	"github.com/mehranzand/pulseup/internal/docker"
	log "github.com/sirupsen/logrus"
)

// GetContainers
// @Summary Get list of containers
func (h *Handler) GetContainers(c echo.Context) error {
	cc := c.(*middleware.DockerContext)

	if cc.Client == nil {
		http.Error(c.Response().Writer, "Docker host not found!", http.StatusInternalServerError)
	}

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

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-transform")
	c.Response().Header().Add(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no")

	cc := c.(*middleware.DockerContext)
	id := c.Param("id")

	f, ok := cc.Context.Response().Writer.(http.Flusher)
	if !ok {

		http.Error(c.Response().Writer, "Streaming unsupported!", http.StatusInternalServerError)
	}

	container, err := cc.Client.FindContainer(id)
	if err != nil {
		log.Error(err)
		http.Error(c.Response().Writer, "Continer not found!", http.StatusInternalServerError)
	}

	since := time.Now().AddDate(0, 0, -100)
	reader, err := cc.Client.ContainerLogs(cc.Context.Request().Context(), container.ID, strconv.FormatInt(since.Unix(), 10), stdTypes)
	if err != nil {
		http.Error(cc.Context.Response().Writer, "Continer not found!", http.StatusInternalServerError)
	}

	lr := docker.NewLogReader(reader, container.Tty)

loop:
	for {
		event, ok := <-lr.Events
		if !ok {
			break loop
		}

		fmt.Fprintf(c.Response().Writer, "data: %s\n", event.Message)
		f.Flush()
	}

	return nil
}
