package handler

import (
	"encoding/json"
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

		log.Tracef("docker event %v", event)

		switch event.Action {
		case "create":
			if container, err := cc.Client.FindContainer(event.ActorID); err == nil {
				response := map[string]interface{}{
					"action":    event.Action,
					"container": &container,
				}
				enc := json.NewEncoder(c.Response())

				if err := enc.Encode(response); err != nil {
					return err
				}

				log.Debugf("container %s was %s", event.ActorID, event.Action)
				fmt.Fprintf(c.Response().Writer, "\n")
			}

		case "start", "die", "destroy":
			response := map[string]interface{}{
				"action": event.Action,
				"id":     event.ActorID,
			}
			enc := json.NewEncoder(c.Response())

			if err := enc.Encode(response); err != nil {
				log.Debugf("container %s was %s", event.ActorID, event.Action)
				fmt.Fprintf(c.Response().Writer, "\n")
			}

		}
		f.Flush()
	}
}
