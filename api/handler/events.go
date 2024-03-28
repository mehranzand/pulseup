package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/mehranzand/pulseup/api/middleware"
	"github.com/mehranzand/pulseup/internal/docker"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) StreamContainerEvents(c echo.Context) error {
	cc := c.(*middleware.DockerContext)

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-transform")
	c.Response().Header().Add(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no")

	f, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		http.Error(c.Response().Writer, "Streaming unsupported!", http.StatusInternalServerError)
		return nil
	}

	events := make(chan docker.ContainerEvent)
	cc.Client.Events(c.Request().Context(), events)

	for {
		event := <-events
		log.Tracef("received event: %+v", event)
		switch event.Name {
		case "create":
			log.Debugf("container %s created", event.ActorID)
			fmt.Fprintf(c.Response().Writer, "event: %s created\n", event.ActorID)
		case "start":
			log.Debugf("container %s started", event.ActorID)
			fmt.Fprintf(c.Response().Writer, "event: %s started\n", event.ActorID)
		case "destroy":
			log.Debugf("container %s destroyed", event.ActorID)
			fmt.Fprintf(c.Response().Writer, "event: %s destroyed\n", event.ActorID)
		case "die":
			log.Debugf("container %s died", event.ActorID)
			fmt.Fprintf(c.Response().Writer, "event: %s died\n", event.ActorID)
		}

		f.Flush()
	}

}
