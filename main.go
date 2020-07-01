package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"gitlab.citodi.com/coretech/esther/logging"
	"gitlab.citodi.com/coretech/esther/persistence"
)

func main() {
	logging.Logger.SetOutput(os.Stdout)
	logging.Logger.Info("Starting Esther")

	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	//	autoCheck()

	r.Run()
}

func setupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logging.GinLogHandler(), gin.Recovery())

	// TEMP
	gin.SetMode(gin.DebugMode)

	/* Monitoring */
	r.GET("/", healthCheck)
	r.GET("/ready", readyCheck)

	/* API */
	api := r.Group("/plans/:id")
	{
		api.GET("/eventCallbacks", getCallbacks)
		api.POST("/eventCallbacks", postOneCallback)
		api.GET("/eventCallbacks/:eventId", getOneCallback)
		api.PUT("/eventCallbacks/:eventId", putOneCallback)
		api.DELETE("/eventCallbacks/:eventId", deleteOneCallback)
		api.DELETE("/eventCallbacks", deleteCallbacks)
		api.PUT("/eventCallbacksToParent", putCallbacksToParent)
	}

	/* OpenAPI doc */
	r.Static("/openapi", "openapi/")

	return r
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func readyCheck(c *gin.Context) {
	errors := autoCheck()
	if len(errors) > 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"errors": errors,
		})
		return
	}
	c.JSON(http.StatusOK, "OK")
}

func autoCheck() map[string][]string {
	errors := make(map[string][]string)
	if persistenceErrors := persistence.ReadyCheck(); len(persistenceErrors) > 0 {
		errors["persistence"] = persistenceErrors
	}
	// TODO
	return errors
}
