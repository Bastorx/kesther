package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getCallbacks(c *gin.Context) {
	c.JSON(http.StatusOK, "GET callbacks list")
}

func postOneCallback(c *gin.Context) {
	c.JSON(http.StatusOK, "POST callback")
}

func getOneCallback(c *gin.Context) {
	eventID := c.Param("eventId")
	c.JSON(http.StatusOK, fmt.Sprintf("GET callback %s", eventID))
}

func putOneCallback(c *gin.Context) {
	eventID := c.Param("eventId")
	c.JSON(http.StatusOK, fmt.Sprintf("PUT callback %s", eventID))
}

func deleteOneCallback(c *gin.Context) {
	eventID := c.Param("eventId")
	c.JSON(http.StatusOK, fmt.Sprintf("DELETE callback %s", eventID))
}

func deleteCallbacks(c *gin.Context) {
	c.JSON(http.StatusOK, "DELETE callbacks list")
}

func putCallbacksToParent(c *gin.Context) {
	c.JSON(http.StatusOK, "PUT callbacks list to parent")
}
