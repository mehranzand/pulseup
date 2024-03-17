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
	//Default route group
	v1 := r.Group("/api")
	// Register routes and handlers
	h.Register(v1)
	//handle embeded react
	web.RegisterHandlers(r)
	// run web server
	r.Logger.Fatal(r.Start(config.Addr))
}
