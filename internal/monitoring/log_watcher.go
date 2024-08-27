package monitoring

import (
	"context"
	"io"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mehranzand/pulseup/internal/docker"
	"github.com/mehranzand/pulseup/internal/models"
	"gorm.io/gorm"
)

type LogWatcher struct {
	clients map[string]docker.Client
	pool    *Observable
	db      *gorm.DB
}

func NewLogWatcher(clients map[string]docker.Client, db *gorm.DB) *LogWatcher {
	logWatcher := &LogWatcher{
		clients: clients,
		db:      db,
		pool:    NewObservable(),
	}

	logWatcher.pool.RegisterListener(func(key string, value interface{}, action string) {
	})

	return logWatcher
}

func (w *LogWatcher) AddContainer(ctr models.MonitoredContainer) {
	if value, exists := w.pool.Get(ctr.ContainerId); exists {
		value.(context.CancelFunc)()
	}

	ctx, cancel := context.WithCancel(context.Background())
	go w.watch(ctx, ctr)
	w.pool.Add(ctr.ContainerId, cancel)
}

func (w *LogWatcher) RemoveContainer(host string, id string) {
	if value, exists := w.pool.Get(id); exists {
		value.(context.CancelFunc)()
		w.pool.Remove(id)
	}
}

func (w *LogWatcher) watch(ctx context.Context, ctr models.MonitoredContainer) {
	var stdTypes docker.StdType
	stdTypes |= docker.STDOUT
	stdTypes |= docker.STDERR
	since := time.Now()
	reader, err := w.clients[ctr.Host].ContainerLogs(ctx, ctr.ContainerId, strconv.FormatInt(since.Unix(), 10), stdTypes)

	if err != nil {
		if err == io.EOF {
			log.Debugf("watcher: container %s was stopped", ctr.ContainerId)

		} else {
			log.Debugf("watcher: container %s error %s", ctr.ContainerId, err.Error())
		}
	}

	lr := docker.NewLogReader(reader, false)

	select {
	case err := <-lr.Errors:
		if err != nil {
			if err == io.EOF {
				log.Debugf("container stopped: %v", ctr.ContainerId)

			} else if err != context.Canceled {
				log.Errorf("unknown error while watching %v", err.Error())
			}
		}
	default:
	}

	NewInspector(lr.Events, ctx, w.db, ctr.Triggers)

	//TODO: Handle inpector errors
}
