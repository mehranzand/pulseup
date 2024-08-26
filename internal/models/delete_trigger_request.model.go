package models

type DeleteTriggerRequest struct {
	TriggerId string `json:"trigger_id" form:"trigger_id" query:"trigger_id"`
}
