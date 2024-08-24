package models

type AddTriggerRequest struct {
	ContainerId string                    `json:"container_id" form:"container_id" query:"container_id"`
	Triggers    []AddTriggerDetailRequest `json:"triggers" form:"triggers" query:"triggers"`
}

type AddTriggerDetailRequest struct {
	Id       uint   `json:"id" form:"id" query:"id"`
	Type     string `json:"type" form:"type" query:"type"`
	Criteria string `json:"criteria" form:"criteria" query:"criteria"`
	Active   bool   `json:"active" form:"active" query:"active"`
}
