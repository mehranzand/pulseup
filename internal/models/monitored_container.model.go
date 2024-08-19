package models

import "gorm.io/gorm"

type MonitoredContainer struct {
	gorm.Model
	Host        string    `json:"host"`
	ContainerId string    `json:"container_id"`
	Active      bool      `json:"active"`
	Triggers    []Trigger `json:"triggers" gorm:"foreignKey:MonitoredContainerID"`
}

func (MonitoredContainer) TableName() string {
	return "monitored_containers"
}
