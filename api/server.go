package api

import (
	"github.com/mehranzand/pulseup/api/handler"
	"github.com/mehranzand/pulseup/api/middleware"
	"github.com/mehranzand/pulseup/api/router"
	"github.com/mehranzand/pulseup/internal/docker"
	"github.com/mehranzand/pulseup/web"
)

func CreateServer(clients map[string]docker.Client, config *handler.Config) {
	h := handler.NewHandler(config)
	r := router.New()
	//Inject docker client to request context
	r.Use(middleware.DockerMiddleware(r, clients))
	//Rebase by config args
	base := r.Group(config.Base + "/api")
	// Register routes and handlers
	h.Register(base)
	//handle embeded react
	web.RegisterHandlers(r)
	// run web server
	r.Logger.Fatal(r.Start(config.Addr))
}
