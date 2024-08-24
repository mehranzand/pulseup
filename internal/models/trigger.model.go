package models

import "gorm.io/gorm"

type TriggerType string

const (
	TriggerTypeLogTypeMatch TriggerType = "log_type_match"
	TriggerTypeTextMatch    TriggerType = "text_match"
	TriggerTypeObjectMatch  TriggerType = "object_match"
)

type Trigger struct {
	gorm.Model
	MonitoredContainerID uint            `json:"monitored_container_id" gorm:"not null"`
	Type                 string          `json:"type"`
	Criteria             string          `json:"criteria" gorm:"not null"`
	Active               bool            `json:"active"`
	Actions              []Action        `json:"actions" gorm:"foreignKey:TriggerID"`
	Logs                 []TriggerLog    `json:"logs" gorm:"foreignKey:TriggerID"`
	Results              []TriggerResult `json:"results" gorm:"foreignKey:TriggerID"`
}

func (Trigger) TableName() string {
	return "triggers"
}
