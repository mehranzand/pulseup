package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mehranzand/pulseup/internal/models"
)

// AddTrigger
// @Summary add new trigger
func (h *Handler) SaveTrigger(c echo.Context) error {
	//cc := c.(*middleware.DockerContext)
	r := new(models.AddTriggerRequest)
	if err := c.Bind(r); err != nil {
		return err
	}

	if len(r.ContainerId) == 0 {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "container_id is required"})
	}

	if len(r.Triggers) == 0 {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "no triggers were found to create"})
	}

	var monitoredContainer models.MonitoredContainer
	result := h.db.Preload("Triggers").First(&monitoredContainer, "container_id = ?", r.ContainerId)

	if result.Error != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{"err": result.Error})
	}

	if result.RowsAffected == 0 {
		monitoredContainer := models.MonitoredContainer{ContainerId: r.ContainerId, Host: "localhost", Active: true}
		h.db.Omit("Triggers").Create(&monitoredContainer)
	} else if !monitoredContainer.Active {
		h.db.Model(&monitoredContainer).Updates(map[string]interface{}{"active": 1, "updated_at": time.Now()})
	}

	for _, trigger := range r.Triggers {
		if trigger.Id != 0 {
			var triggerToUpdate *models.Trigger
			for _, t := range monitoredContainer.Triggers {
				if t.ID == trigger.Id {
					triggerToUpdate = &t
					break
				}
			}

			if triggerToUpdate != nil {
				triggerToUpdate.Type = trigger.Type
				triggerToUpdate.Criteria = trigger.Criteria
				triggerToUpdate.Active = trigger.Active
				triggerToUpdate.UpdatedAt = time.Now()

				h.db.Save(&triggerToUpdate)
			} else {
				return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "trigger not found to update"})
			}

		} else {
			h.db.Model(&monitoredContainer).Association("Triggers").Append(&models.Trigger{Type: trigger.Type, Criteria: trigger.Criteria, Active: trigger.Active})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
}

// DeleteTrigger
// @Summary delete a trigger
func (h *Handler) DeleteTrigger(c echo.Context) error {
	return c.String(http.StatusOK, "Delete!\n")
}

// EditTrigger
// @Summary edit a trigger
func (h *Handler) EditTrigger(c echo.Context) error {

	return c.String(http.StatusOK, "Edit!\n")
}
