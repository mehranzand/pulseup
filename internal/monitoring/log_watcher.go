package action

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/mehranzand/pulseup/internal/docker"
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

func (w *LogWatcher) AddContainer(host string, id string) {
	ctx, cancel := context.WithCancel(context.Background())
	go w.addContainer(ctx, host, id)
	w.observatory[id] = cancel
}

func (w *LogWatcher) RemoveContainer(host string, id string) {
	w.observatory[id].(context.CancelFunc)()
	delete(w.observatory, id)
}

func (w *LogWatcher) addContainer(ctx context.Context, host string, id string) {
	var stdTypes docker.StdType
	stdTypes |= docker.STDOUT
	stdTypes |= docker.STDERR
	since := time.Now().AddDate(0, 0, -10)
	reader, err := w.Clients[host].ContainerLogs(ctx, id, strconv.FormatInt(since.Unix(), 10), stdTypes)

	if err != nil {
		if err == io.EOF {
			fmt.Printf("event: container-stopped\ndata: end of stream")

		} else {
			fmt.Printf("watcher: %s", err.Error())
		}
	}

	lr := docker.NewLogReader(reader, false)
	for {
		select {
		case event := <-lr.Events:
			if buffer, err := json.Marshal(event); err != nil {
				log.Errorf("json encoding error while streaming %v", err.Error())
			} else {
				fmt.Printf("%s\n", buffer)
			}
		case <-ctx.Done():
			return
		}
	}
}