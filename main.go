package main

import (
	"fmt"
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
	autoCheck()

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
	api := r.Group("/plans/:planId")
	{
		api.GET("/eventCallbacks", getCallbacks)
		api.POST("/eventCallbacks", postOneCallback)
		api.GET("/eventCallbacks/:eventId", getOneCallback)
		api.PUT("/eventCallbacks/:eventId", putOneCallback)
		api.DELETE("/eventCallbacks/:eventId", deleteOneCallback)
		api.DELETE("/eventCallbacks", deleteCallbacks)
		api.PUT("/eventCallbacksToParent", putCallbacksToParent)
	}

	/* Local commands */
	local := r.Group("", isLocalhost)
	{
		local.GET("/reset", doReset)
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

func isLocalhost(c *gin.Context) {
	clientIP := c.ClientIP()
	logging.Logger.WithField("clientIP", clientIP).Info("Checking localhost by IP")
	if clientIP != "127.0.0.1" && clientIP != "::1" {
		err := fmt.Errorf("%s is not authorized for this action", clientIP)
		abortWithError(c, http.StatusUnauthorized, err, "Unauthorized IP")
	}
}

func doReset(c *gin.Context) {
	logging.Logger.Info("Reset")
	errors := reset()
	logging.Logger.Info("End of reset")

	if len(errors) > 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"errors": errors,
		})
		return
	}
	c.JSON(http.StatusOK, "OK")
}

func abortWithError(c *gin.Context, status int, err error, errTitle string) {
	errDetail := ""
	if err != nil {
		errDetail = fmt.Sprintf("%s", err)
	}
	logging.Logger.WithField("error", err).Error(errTitle)
	c.AbortWithStatusJSON(status, gin.H{
		"error": gin.H{
			"title":  errTitle,
			"detail": errDetail,
		},
	})
}

func autoCheck() map[string][]string {
	errors := make(map[string][]string)
	if persistenceErrors := persistence.ReadyCheck(); len(persistenceErrors) > 0 {
		errors["persistence"] = persistenceErrors
	}
	return errors
}

func reset() []string {
	errors := []string{}
	errors = append(errors, persistence.Reset()...)
	return errors
}
