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
	clients     map[string]docker.Client
	objectPool  *Observable
	contextPool *Observable
	db          *gorm.DB
}

func NewLogWatcher(clients map[string]docker.Client, db *gorm.DB) *LogWatcher {
	logWatcher := &LogWatcher{
		clients:     clients,
		db:          db,
		objectPool:  NewObservable(),
		contextPool: NewObservable(),
	}

	logWatcher.objectPool.RegisterListener(func(key string, value interface{}, action string) {
		if action == "save" {
			if ctx, exists := logWatcher.contextPool.Get(key); exists {
				ctx.(context.CancelFunc)()
			}

			ctx, cancel := context.WithCancel(context.Background())
			go logWatcher.watch(ctx, value.(models.MonitoredContainer))
			logWatcher.contextPool.Add(key, cancel)
		} else if action == "remove" {
			if value, exists := logWatcher.contextPool.Get(key); exists {
				value.(context.CancelFunc)()
				logWatcher.contextPool.Remove(key)
			}
		}
	})

	return logWatcher
}

func (w *LogWatcher) WatchContainer(ctr models.MonitoredContainer) {
	w.objectPool.Add(ctr.ContainerId, ctr)
}

func (w *LogWatcher) RemoveTrigger(cid string, tid uint) {
	if value, exists := w.objectPool.Get(cid); exists {
		var monitoredContainer = value.(models.MonitoredContainer)
		var index int
		for i, t := range monitoredContainer.Triggers {
			if t.ID == tid {
				index = i
				break
			}
		}

		triggers := monitoredContainer.Triggers
		triggers = append(triggers[:index], triggers[index+1:]...)
		monitoredContainer.Triggers = triggers

		if len(monitoredContainer.Triggers) == 0 {
			w.objectPool.Remove(cid)
		} else {

			w.WatchContainer(monitoredContainer)
		}
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
