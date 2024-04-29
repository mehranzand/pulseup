package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {
	v1.GET("/:host/containers", h.GetContainers)
	v1.GET("/logs/stream/:host/:id", h.StreamLogs)
	v1.GET("/events/stream/:host", h.StreamContainerEvents)
}
