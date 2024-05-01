package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mehranzand/pulseup/internal/docker"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) Register(v1 *echo.Group) {
	v1.GET("/:host/containers", h.GetContainers)
	v1.GET("/logs/stream/:host/:id", h.StreamLogs)
	v1.GET("/events/stream/:host", h.StreamContainerEvents)
}

func (h *Handler) IndexHandler(c echo.Context) error {
	hosts := make([]*docker.Host, 0, len(h.clients))
	for _, v := range h.clients {
		hosts = append(hosts, v.Host())
	}

	config := map[string]interface{}{
		"authProvider": h.config.AuthProvider,
		"version":      h.config.Version,
		"hostname":     h.config.Hostname,
		"hosts":        hosts,
	}

	data := map[string]interface{}{
		"Config": config,
	}

	err := h.indexTmpl.Execute(c.Response().Writer, data)
	if err != nil {
		log.Panic(err)
		http.Error(c.Response().Writer, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
