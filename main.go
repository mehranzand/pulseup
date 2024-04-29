package main

import (
	"embed"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/mehranzand/pulseup/internal/api"
	"github.com/mehranzand/pulseup/internal/api/handler"
	"github.com/mehranzand/pulseup/internal/docker"

	"github.com/alexflint/go-arg"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	version = "development"
)

//go:embed all:web/dist
var content embed.FS

type args struct {
	Addr                 string              `arg:"env:PULSEUP_ADDR" default:":7070" help:"sets host:port to bind for server. This is rarely needed inside a docker container."`
	Base                 string              `arg:"env:PULSEUP_BASE" default:"/" help:"sets the base for http router."`
	Hostname             string              `arg:"env:PULSEUP_HOSTNAME" help:"sets the hostname for display. This is useful with multiple pulseUp instances."`
	LogLevel             string              `arg:"env:PULSEUP_LOGLEVEL" default:"info" help:"set pulseUp log level. Use debug for more logging."`
	Username             string              `arg:"env:PULSEUP_USERNAME" help:"sets the username for auth."`
	Password             string              `arg:"env:PULSEUP_PASSWORD" help:"sets password for auth"`
	FilterStrings        []string            `arg:"env:PULSEUP_FILTER,--filter,separate" help:"filters docker containers using Docker syntax."`
	Filter               map[string][]string `arg:"-"`
	WaitForDockerSeconds int                 `arg:"--wait-for-docker-seconds,env:PULEUP_WAIT_FOR_DOCKER_SECONDS" help:"wait for docker to be available for at most this many seconds before starting the server."`
	RemoteHost           []string            `arg:"env:PULSEUP_REMOTE_HOST,--remote-host" help:"list of remote address of re dockerd to connect remotely"`
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

	log.Infof("💡 pulseUp version %s", version)

	clients := createDockerClients(args.Hostname, args)

	if len(clients) == 0 {
		log.Fatal("Could not connect to any Dockerd")
	} else {
		log.Infof("Totaly connected to %d Dockerd(s)", len(clients))
	}

	createServer(args, clients)
}

func createDockerClients(hostname string, args args) map[string]docker.Client {
	clients := make(map[string]docker.Client)

	for i := 1; ; i++ {
		dockerClient, err := docker.NewLocalClientWithOpts(args.Filter)

		if err == nil {
			if hostname != "" {
				dockerClient.Host().Name = hostname
			}

			_, err := dockerClient.ListContainers()

			if err != nil {
				log.Debugf("Could not connect to local Dockerd: %s", err)
			} else {
				log.Debugf("Connected to local Dockerd")

				clients[dockerClient.Host().ID] = dockerClient

				break
			}
		}
		if args.WaitForDockerSeconds > 0 {
			log.Infof("Waiting for Dockerd (attempt %d)", i)
			time.Sleep(5 * time.Second)
			args.WaitForDockerSeconds -= 2
		} else {
			log.Debugf("Local Dockerd not found")
			break
		}
	}

	for _, remoteHost := range args.RemoteHost {
		host, err := docker.ParseConnection(remoteHost)
		if err != nil {
			log.Fatalf("Could not parse remote host %s: %s", remoteHost, err)
		}

		log.Infof("Creating remote client for %s", host.URL.String())
		if client, err := docker.NewTLSClientWithOpts(args.Filter, host); err == nil {
			if _, err := client.ListContainers(); err == nil {
				log.Debugf("Connected to remote Dockerd")
				clients[client.Host().URL.Hostname()] = client
			} else {
				log.Warnf("Could not connect to remote host %s: %s", host.ID, err)
			}
		} else {
			log.Warnf("Could not create client for %s: %s", host.ID, err)
		}
	}

	return clients
}

func createServer(args args, clients map[string]docker.Client) {
	config := handler.Config{
		Addr:         args.Addr,
		Base:         args.Base,
		Version:      version,
		Hostname:     args.Hostname,
		AuthProvider: handler.NONE,
	}

	assets, err := fs.Sub(content, "web/dist")
	if err != nil {
		log.Fatalf("Could not open web content at dist folder: %v", err)
	}

	api.CreateServer(clients, &config, assets)
}

func parseArgs() args {
	var args args
	parser := arg.MustParse(&args)

	configureLogger(args.LogLevel)

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
	log.SetOutput(os.Stdout)

	if l, err := log.ParseLevel(level); err == nil {
		log.SetLevel(l)
	} else {
		panic(err)
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})
}
