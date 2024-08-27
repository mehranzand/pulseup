package monitoring

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/mehranzand/pulseup/internal/docker"
	"github.com/mehranzand/pulseup/internal/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Inspector struct {
	Events   chan *docker.LogEvent
	ctx      context.Context
	db       *gorm.DB
	triggers []models.Trigger
}

func NewInspector(events chan *docker.LogEvent, ctx context.Context, db *gorm.DB, triggers []models.Trigger) *Inspector {
	inspector := &Inspector{
		Events:   events,
		ctx:      ctx,
		db:       db,
		triggers: triggers,
	}

	go inspector.startListening()

	return inspector
}

func (i *Inspector) startListening() {
	for {
		select {
		case event := <-i.Events:
			if buffer, err := json.Marshal(event); err != nil {
				log.Errorf("json encoding error while watching %v", err.Error())
			} else {
				i.matchCriteria(buffer)

			}
		case <-i.ctx.Done():
			return
		}
	}
}

func (i *Inspector) matchCriteria(buffer []byte) {
	for _, trigger := range i.triggers {
		result, _ := regexp.Match(trigger.Criteria, buffer)

		if result {
			log.Println(string(buffer))
			log.Infof("Result is: %s for triggerId => %b", trigger.Criteria, trigger.ID)
		}
	}
}
