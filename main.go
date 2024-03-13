package main

import (
	"net/http"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/mehranzand/pulseup/web"
	log "github.com/sirupsen/logrus"
)

var (
	version = "head"
)

type args struct {
	Addr          string              `arg:"env:PULSEUP_ADDR" default:":7070" help:"sets host:port to bind for server. This is rarely needed inside a docker container."`
	Base          string              `arg:"env:PULSEUP_BASE" default:"/" help:"sets the base for http router."`
	Hostname      string              `arg:"env:PULSEUP_HOSTNAME" help:"sets the hostname for display. This is useful with multiple pulseUp instances."`
	Level         string              `arg:"env:PULSEUP_LEVEL" default:"info" help:"set pulseUp log level. Use debug for more logging."`
	Username      string              `arg:"env:PULSEUP_USERNAME" help:"sets the username for auth."`
	Password      string              `arg:"env:PULSEUP_PASSWORD" help:"sets password for auth"`
	FilterStrings []string            `arg:"env:PULSEUP_FILTER,--filter,separate" help:"filters docker containers using Docker syntax."`
	Filter        map[string][]string `arg:"-"`
}

func (args) Version() string {
	return version
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	args := parseArgs()
	if len(args.Addr) == 0 {
		log.Fatal("PULSEUP_ADDR can't be null or empty")
	}

	log.Infof("pulseUp version %s", version)

	e := echo.New()
	web.RegisterHandlers(e)
	e.GET("/api", func(c echo.Context) error {
		return c.String(http.StatusOK, "pulseUp")
	})

	e.Logger.Fatal(e.Start(args.Addr))
}

func parseArgs() args {
	var args args
	parser := arg.MustParse(&args)

	configureLogger(args.Level)

	args.Filter = make(map[string][]string)

	for _, filter := range args.FilterStrings {
		pos := strings.Index(filter, "=")
		if pos == -1 {
			parser.Fail("each filter should be of the form key=value")
		}
		key := filter[:pos]
		val := filter[pos+1:]
		args.Filter[key] = append(args.Filter[key], val)
	}

	return args
}

func configureLogger(level string) {
	if l, err := log.ParseLevel(level); err == nil {
		log.SetLevel(l)
	} else {
		panic(err)
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
}
