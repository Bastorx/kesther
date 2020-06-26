package main

import (
	"net/http"
	"os"

	"gitlab.citodi.com/coretech/esther/logging"

	"github.com/gin-gonic/gin"
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

	/* Monitoring */
	r.GET("/", healthCheck)
	r.GET("/ready", readyCheck)

	/* API */
	api := r.Group("/plans/:id")
	{
		// TODO
		api.GET("/", healthCheck)
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
	// TODO
	return errors
}
