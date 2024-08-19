package models

import "gorm.io/gorm"

type ActionType string

const (
	ActionTypeSendEmail ActionType = "send_email"
)

type Action struct {
	gorm.Model
	Type      ActionType `json:"type"`
	TriggerID uint       `json:"trigger_id" gorm:"not null"`
}

func (Action) TableName() string {
	return "actions"
}
