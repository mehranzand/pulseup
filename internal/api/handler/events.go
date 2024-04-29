package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/mehranzand/pulseup/internal/api/middleware"
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
		case "create", "start", "die":
			if container, err := cc.Client.FindContainer(event.ActorID); err == nil {
				s := map[string]interface{}{
					"action":    event.Action,
					"container": &container,
				}

				if buffer, err := json.Marshal(s); err != nil {
					log.Errorf("json encoding error while streaming %v", err.Error())
				} else {
					fmt.Fprintf(c.Response().Writer, "event: container\ndata: %s\n\n", buffer)
				}

				log.Debugf("container %s was %s", event.ActorID, event.Action)
			}

		case "destroy":
			s := map[string]interface{}{
				"action":    event.Action,
				"container": map[string]interface{}{"id": event.ActorID},
			}

			if buffer, err := json.Marshal(s); err != nil {
				log.Errorf("json encoding error while streaming %v", err.Error())
			} else {
				fmt.Fprintf(c.Response().Writer, "event: container\ndata: %s\n\n", buffer)
			}
		}

		f.Flush()
	}
}
