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
	Clients     map[string]docker.Client
	observatory map[string]interface{}
	db          *gorm.DB
}

func NewLogWatcher(clients map[string]docker.Client, db *gorm.DB) *LogWatcher {
	logWatcher := &LogWatcher{
		Clients:     clients,
		observatory: make(map[string]interface{}),
		db:          db,
	}

	return logWatcher
}

func (w *LogWatcher) AddContainer(ctr models.MonitoredContainer) {
	ctx, cancel := context.WithCancel(context.Background())
	go w.addContainer(ctx, ctr)
	w.observatory[ctr.ContainerId] = cancel
}

func (w *LogWatcher) RemoveContainer(host string, id string) {
	w.observatory[id].(context.CancelFunc)()
	delete(w.observatory, id)
}

func (w *LogWatcher) addContainer(ctx context.Context, ctr models.MonitoredContainer) {
	var stdTypes docker.StdType
	stdTypes |= docker.STDOUT
	stdTypes |= docker.STDERR
	since := time.Now().AddDate(0, 0, -1)
	reader, err := w.Clients[ctr.Host].ContainerLogs(ctx, ctr.ContainerId, strconv.FormatInt(since.Unix(), 10), stdTypes)

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
