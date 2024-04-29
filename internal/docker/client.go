package docker

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

type StdType int

const (
	UNKNOWN StdType = 1 << iota
	STDOUT
	STDERR
)
const STDALL = STDOUT | STDERR

func (s StdType) String() string {
	switch s {
	case STDOUT:
		return "stdout"
	case STDERR:
		return "stderr"
	case STDALL:
		return "all"
	default:
		return "unknown"
	}
}

type DockerCLI interface {
	ContainerList(context.Context, container.ListOptions) ([]types.Container, error)
	ContainerLogs(context.Context, string, container.LogsOptions) (io.ReadCloser, error)
	Events(context.Context, types.EventsOptions) (<-chan events.Message, <-chan error)
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
	ContainerStats(ctx context.Context, containerID string, stream bool) (types.ContainerStats, error)
	Ping(ctx context.Context) (types.Ping, error)
	ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
	ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error
}

type Client interface {
	Host() *Host
	ListContainers() ([]Container, error)
	FindContainer(string) (Container, error)
	ContainerLogs(context.Context, string, string, StdType) (io.ReadCloser, error)
	ContainerLogsRange(context.Context, string, time.Time, time.Time, StdType) (io.ReadCloser, error)
	ContainerActions(action string, containerID string) error
	Events(context.Context, chan<- ContainerEvent)
	Ping(context.Context) (types.Ping, error)
}
type _client struct {
	cli       DockerCLI
	filters   filters.Args
	host      *Host
	SwarmMode bool
}

func NewClient(cli DockerCLI, filters filters.Args, host *Host, swarm bool) Client {
	return &_client{cli, filters, host, swarm}
}

func NewLocalClientWithOpts(f map[string][]string) (Client, error) {
	filterArgs := filters.NewArgs()

	for key, values := range f {
		for _, value := range values {
			filterArgs.Add(key, value)
		}
	}

	log.Debugf("localhost filterArgs=%v", filterArgs)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, err
	}

	info, _ := cli.Info(context.Background())
	swarm := info.Swarm.LocalNodeState != swarm.LocalNodeStateInactive

	return NewClient(cli, filterArgs, &Host{Name: "localhost", ID: "localhost"}, swarm), nil
}

func NewTLSClientWithOpts(f map[string][]string, host Host) (Client, error) {
	filterArgs := filters.NewArgs()
	for key, values := range f {
		for _, value := range values {
			filterArgs.Add(key, value)
		}
	}

	log.Debugf("%s filterArgs = %v", host.Name, filterArgs)

	if host.URL.Scheme != "tcp" {
		log.Fatal("Only tcp scheme is supported")
	}

	opts := []client.Opt{
		client.WithHost(host.URL.String()),
	}

	if host.ValidCerts {
		log.Debugf("Using TLS client config with certs at: %s", filepath.Dir(host.CertPath))
		opts = append(opts, client.WithTLSClientConfig(host.CACertPath, host.CertPath, host.KeyPath))
	} else {
		log.Debugf("No valid certs found, using plain TCP")
	}

	opts = append(opts, client.WithAPIVersionNegotiation())

	cli, err := client.NewClientWithOpts(opts...)

	if err != nil {
		return nil, err
	}

	info, _ := cli.Info(context.Background())
	swarm := info.Swarm.LocalNodeState != swarm.LocalNodeStateInactive

	return NewClient(cli, filterArgs, &host, swarm), nil
}

func (d *_client) FindContainer(id string) (Container, error) {
	var container Container
	containers, err := d.ListContainers()
	if err != nil {
		return container, err
	}

	found := false
	for _, c := range containers {
		if c.ID == id {
			container = c
			found = true
			break
		}
	}
	if !found {
		return container, fmt.Errorf("unable to find container with id: %s", id)
	}

	if json, err := d.cli.ContainerInspect(context.Background(), container.ID); err == nil {
		container.Tty = json.Config.Tty
	} else {
		return container, err
	}

	return container, nil
}

func (d *_client) ContainerActions(action string, containerID string) error {
	switch action {
	case "start":
		return d.cli.ContainerStart(context.Background(), containerID, container.StartOptions{})
	case "stop":
		return d.cli.ContainerStop(context.Background(), containerID, container.StopOptions{})
	case "restart":
		return d.cli.ContainerRestart(context.Background(), containerID, container.StopOptions{})
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

func (d *_client) ListContainers() ([]Container, error) {
	containerListOptions := container.ListOptions{
		Filters: d.filters,
		All:     true,
	}
	list, err := d.cli.ContainerList(context.Background(), containerListOptions)
	if err != nil {
		return nil, err
	}

	var containers = make([]Container, 0, len(list))
	for _, c := range list {
		name := "no name"
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}
		container := Container{
			ID:      c.ID[:12],
			Names:   c.Names,
			Name:    name,
			Image:   c.Image,
			ImageID: c.ImageID,
			Command: c.Command,
			Created: c.Created,
			State:   c.State,
			Status:  c.Status,
			Host:    d.host.ID,
			Labels:  c.Labels,
		}
		containers = append(containers, container)
	}

	sort.Slice(containers, func(i, j int) bool {
		return strings.ToLower(containers[i].Name) < strings.ToLower(containers[j].Name)
	})

	return containers, nil
}

func (d *_client) ContainerLogs(ctx context.Context, id string, since string, stdType StdType) (io.ReadCloser, error) {
	log.WithField("id", id).WithField("since", since).WithField("stdType", stdType).Debug("streaming logs for container")

	if since != "" {
		if m, err := strconv.ParseInt(since, 10, 64); err == nil {
			since = time.Unix(m, 0).Format(time.RFC3339Nano)
		} else {
			log.WithError(err).Debug("unable to parse since")
		}
	}

	options := container.LogsOptions{
		ShowStdout: stdType&STDOUT != 0,
		ShowStderr: stdType&STDERR != 0,
		Follow:     true,
		Tail:       "400",
		Timestamps: true,
		Since:      since,
	}

	reader, err := d.cli.ContainerLogs(ctx, id, options)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (d *_client) Events(ctx context.Context, messages chan<- ContainerEvent) {
	dockerMessages, errors := d.cli.Events(ctx, types.EventsOptions{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-errors:
				log.Fatalf("Exiting while listening to docker events an error occurred: %v.", err)
			case message, ok := <-dockerMessages:
				if !ok {
					log.Errorf("docker events channel closed")
					return
				}

				if message.Type == "container" && len(message.Actor.ID) > 0 {
					messages <- ContainerEvent{
						ActorID: message.Actor.ID[:12],
						Action:  string(message.Action),
						Host:    d.host.ID,
					}
				}
			}
		}
	}()
}

func (d *_client) ContainerLogsRange(ctx context.Context, id string, from time.Time, to time.Time, stdType StdType) (io.ReadCloser, error) {
	options := container.LogsOptions{
		ShowStdout: stdType&STDOUT != 0,
		ShowStderr: stdType&STDERR != 0,
		Timestamps: true,
		Since:      from.Format(time.RFC3339Nano),
		Until:      to.Format(time.RFC3339Nano),
	}

	log.Debugf("fetching logs from Docker with option: %+v", options)

	reader, err := d.cli.ContainerLogs(ctx, id, options)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (d *_client) Ping(ctx context.Context) (types.Ping, error) {
	return d.cli.Ping(ctx)
}

func (d *_client) Host() *Host {
	return d.host
}
