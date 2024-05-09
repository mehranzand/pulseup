package models

import "gorm.io/gorm"

type LogWatcher struct {
	gorm.Model
	ContainerId string
	Host        string
}
