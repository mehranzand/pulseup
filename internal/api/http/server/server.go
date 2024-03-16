package server

import (
	"net/http"

	"github.com/mehranzand/pulseup/internal/docker"

	"github.com/labstack/echo/v4"

	"github.com/mehranzand/pulseup/web"
	log "github.com/sirupsen/logrus"
)

type AuthProvider string

const (
	NONE  AuthProvider = "none"
	BASIC AuthProvider = "basic"
)

type Config struct {
	Base         string
	Addr         string
	Version      string
	Hostname     string
	AuthProvider AuthProvider
}

func CreateServer(client docker.Client, config Config) *echo.Echo {
	list, _ := client.ListContainers()

	for container := range list {
		c := list[container]
		log.Info(c.Name)
	}

	server := echo.New()

	web.RegisterHandlers(server)

	server.GET("/api", func(c echo.Context) error {
		return c.String(http.StatusOK, "pulseUp")
	})

	server.Logger.Fatal(server.Start(config.Addr))

	return server
}
