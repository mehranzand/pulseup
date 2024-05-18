package models

import "gorm.io/gorm"

type WatchedContainer struct {
	gorm.Model
	Host        string      `json:"host"`
	ContainerId string      `json:"container_id"`
	Active      bool        `json:"active"`
	Actions     []LogAction `json:"actions" gorm:"foreignKey:WatchedContainerID"`
}

func (WatchedContainer) TableName() string {
	return "watched_containers"
}
