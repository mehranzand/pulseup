package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mehranzand/pulseup/web"
)

func main() {
	e := echo.New()
	web.RegisterHandlers(e)
	e.GET("/api", func(c echo.Context) error {
		return c.String(http.StatusOK, "pulseUp")
	})

	e.Logger.Fatal(e.Start(":7070"))
}
