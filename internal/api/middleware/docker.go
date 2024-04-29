package middleware

import (
	"strings"

	"github.com/mehranzand/pulseup/internal/docker"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type DockerContext struct {
	echo.Context
	Client docker.Client
}

func DockerMiddleware(e *echo.Echo, clients map[string]docker.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			if strings.Contains(strings.ToLower(c.Request().RequestURI), "api") { //attach docker client to only api path
				cc := &DockerContext{}
				cc.Context = c

				if client, ok := clients[c.Param("host")]; ok {
					cc.Client = client
				} else {
					log.Debugf("Middleware could not be infer the Docker client from the URL host:%s", c.Param("host"))
				}

				return next(cc)
			}

			return next(c)
		})
	}
}
