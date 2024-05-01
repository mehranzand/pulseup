package api

import (
	"io/fs"

	"github.com/mehranzand/pulseup/internal/api/handler"
	"github.com/mehranzand/pulseup/internal/api/middleware"
	"github.com/mehranzand/pulseup/internal/api/router"
	"github.com/mehranzand/pulseup/internal/docker"
)

func CreateServer(clients map[string]docker.Client, config *handler.Config, assets fs.FS) {
	h := handler.NewHandler(config, clients, assets)
	r := router.New(assets, h)
	r.Use(middleware.DockerMiddleware(r, clients))
	base := r.Group(config.Base + "/api")
	h.Register(base)
	r.Logger.Fatal(r.Start(config.Addr))
}
