package router

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mehranzand/pulseup/internal/api/handler"

	log "github.com/sirupsen/logrus"
)

func New(assets fs.FS, h *handler.Handler) *echo.Echo {
	e := echo.New()

	e.GET("/", h.IndexHandler)
	e.GET("/*", h.IndexHandler, middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "/",
		HTML5:      true,
		Filesystem: http.FS(assets),
	}))

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(log.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	}))

	return e
}
