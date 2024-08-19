package models

import (
	"time"

	"gorm.io/gorm"
)

type TriggerResult struct {
	gorm.Model
	TriggerID       uint      `json:"trigger_id" gorm:"not null"`
	OccurrenceCount uint      `json:"occurrence_count" gorm:"not null"`
	StartDate       time.Time `json:"start_date" gorm:"not null"`
	EndDate         time.Time `json:"end_date" gorm:"not null"`
}

func (TriggerResult) TableName() string {
	return "trigger_results"
}
