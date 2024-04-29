package router

import (
	"html/template"
	"io"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	log "github.com/sirupsen/logrus"
)

var (
	indexTmpl *template.Template
)

func New(assets fs.FS) *echo.Echo {
	initIndexTemplate(assets)

	e := echo.New()

	e.GET("/", indexHandler)
	e.GET("/*", indexHandler, middleware.StaticWithConfig(middleware.StaticConfig{
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

func initIndexTemplate(assets fs.FS) {
	file, err := assets.Open("index.html")
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Panic(err)
	}

	indexTmpl, err = template.New("index.html").Parse(string(bytes))
	if err != nil {
		log.Panic(err)
	}
}

func indexHandler(c echo.Context) error {
	data := map[string]interface{}{
		"Config": "config json",
	}

	err := indexTmpl.Execute(c.Response().Writer, data)
	if err != nil {
		log.Panic(err)
		http.Error(c.Response().Writer, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
