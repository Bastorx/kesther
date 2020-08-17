package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab.citodi.com/coretech/esther/model"
)

func getCallbacks(c *gin.Context) {
	planID := c.Param("planId")

	eventCallbacks, err := model.FindEventCallbacksByPlanId(planID)
	if err != nil {
		c.Status(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, eventCallbacks)
	}
}

func getOneCallback(c *gin.Context) {
	planID := c.Param("planId")
	eventID := c.Param("eventId")

	eventCallback, err := model.FindEventCallbackById(planID, eventID)
	if err != nil {
		c.Status(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, eventCallback)
	}
}

func postOneCallback(c *gin.Context) {
	planID := c.Param("planId")

	var eventCallback model.EventCallback
	if err := c.ShouldBindBodyWith(&eventCallback, binding.JSON); err != nil {
		abortWithError(c, http.StatusBadRequest, err, "The callback input payload could not be bound")
		return
	}
	createdEventCallback, err := model.CreateEventCallback(planID, eventCallback)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err, "The callback could not be created")
		return
	}
	c.JSON(http.StatusCreated, createdEventCallback)
}

func putOneCallback(c *gin.Context) {
	planID := c.Param("planId")
	eventID := c.Param("eventId")

	var eventCallback model.EventCallback
	if err := c.ShouldBindBodyWith(&eventCallback, binding.JSON); err != nil {
		abortWithError(c, http.StatusBadRequest, err, "The callback input payload could not be bound")
		return
	}
	updatedEventCallback, err := model.UpdateEventCallback(planID, eventID, eventCallback)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err, "The callback could not be updated")
		return
	}
	c.JSON(http.StatusOK, updatedEventCallback)
}

func putCallbacksToParent(c *gin.Context) {
	planID := c.Param("planId")
	c.JSON(http.StatusOK, fmt.Sprintf("PUT callbacks list to parent %s", planID))
}

func deleteOneCallback(c *gin.Context) {
	planID := c.Param("planId")
	eventID := c.Param("eventId")

	if _, err := model.FindEventCallbackById(planID, eventID); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if err := model.DeleteEventCallbacksById(planID, eventID); err != nil {
		abortWithError(c, http.StatusBadRequest, err, "The event callback could not be deleted")
		return
	}
	c.Status(http.StatusNoContent)
}

func deleteCallbacks(c *gin.Context) {
	planID := c.Param("planId")

	if _, err := model.FindEventCallbacksByPlanId(planID); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if err := model.DeleteEventCallbacksByPlanId(planID); err != nil {
		abortWithError(c, http.StatusBadRequest, err, "The events could not be deleted")
		return
	}
	c.Status(http.StatusNoContent)
}
