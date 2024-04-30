package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) Register(v1 *echo.Group) {
	v1.GET("/:host/containers", h.GetContainers)
	v1.GET("/logs/stream/:host/:id", h.StreamLogs)
	v1.GET("/events/stream/:host", h.StreamContainerEvents)
}

func (h *Handler) IndexHandler(c echo.Context) error {
	data := map[string]interface{}{
		"Config": h.config,
	}

	err := h.indexTmpl.Execute(c.Response().Writer, data)
	if err != nil {
		log.Panic(err)
		http.Error(c.Response().Writer, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
