package models

import (
	"gorm.io/gorm"
)

type TriggerLog struct {
	gorm.Model
	TriggerID uint `json:"trigger_id" gorm:"not null"`
}

func (TriggerLog) TableName() string {
	return "trigger_logs"
}
