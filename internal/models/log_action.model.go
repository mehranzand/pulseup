package models

import "gorm.io/gorm"

type LogActionType string

const (
	Error     LogActionType = "error"
	TextMatch LogActionType = "text_match"
)

type LogAction struct {
	gorm.Model
	Type               LogActionType `json:"type"`
	WatchedContainerID uint          `gorm:"not null"`
}

func (LogAction) TableName() string {
	return "log_actions"
}
