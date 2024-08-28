package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mehranzand/pulseup/internal/api/middleware"
	"github.com/mehranzand/pulseup/internal/models"
)

// SaveTrigger
// @Summary add or edit triggers
func (h *Handler) SaveTrigger(c echo.Context) error {
	cc := c.(*middleware.DockerContext)
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
	result := h.db.Preload("Triggers").Find(&monitoredContainer, "container_id = ?", r.ContainerId)

	if result.Error != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{"err": result.Error.Error()})
	}

	if result.RowsAffected == 0 {
		monitoredContainer = models.MonitoredContainer{ContainerId: r.ContainerId, Host: cc.Client.Host().ID, Active: true}
		h.db.Create(&monitoredContainer)
	} else if !monitoredContainer.Active {
		h.db.Model(&monitoredContainer).Updates(map[string]interface{}{"active": 1, "updated_at": time.Now()})
	}

	for _, trigger := range r.Triggers {
		var duplicate *models.Trigger
		for _, t := range monitoredContainer.Triggers {
			if t.Type == trigger.Type && t.Criteria == trigger.Criteria {
				duplicate = &t
				break
			}
		}

		if duplicate != nil {
			continue
		}

		if trigger.Id != 0 {
			var triggerToUpdate *models.Trigger
			for _, t := range monitoredContainer.Triggers {
				if t.ID == trigger.Id && t.Type == trigger.Type && t.Criteria != trigger.Criteria {
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
			trigger := models.Trigger{Type: trigger.Type, Criteria: trigger.Criteria, Active: trigger.Active}
			h.db.Model(&monitoredContainer).Association("Triggers").Append(&trigger)
		}
	}

	h.watcher.TrackContainer(monitoredContainer)

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
}

// DeleteTrigger
// @Summary delete a trigger
func (h *Handler) DeleteTrigger(c echo.Context) error {
	id := c.Param("id")

	if len(id) == 0 {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "trigger id is required"})
	}

	trigger := models.Trigger{}
	result := h.db.First(&trigger, id)
	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"err": result.Error.Error()})
	}

	monitoredContainer := models.MonitoredContainer{}
	result = h.db.First(&monitoredContainer, trigger.MonitoredContainerID)
	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"err": result.Error.Error()})
	}

	deleteResult := h.db.Delete(&trigger, id)
	if deleteResult.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"err": deleteResult.Error})
	} else {
		i, _ := strconv.Atoi(id)
		h.watcher.RemoveTrigger(monitoredContainer.ContainerId, uint(i))
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
	}
}
