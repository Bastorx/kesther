package logging

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Logger logrus.Logger

type LogFields = map[string]interface{}

func init() {
	var level logrus.Level
	var logLevel = os.Getenv("LOG_LEVEL")

	level, e := logrus.ParseLevel(logLevel)
	if e != nil {
		level = logrus.InfoLevel
	}

	var formatter = new(logrus.JSONFormatter)
	formatter.DisableTimestamp = true

	Logger = logrus.Logger{
		Out:       os.Stdout,
		Formatter: formatter,
		Level:     level,
	}
}

func GinLogHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		status := c.Writer.Status()
		fields := LogFields{
			"httpCode":   status,
			"path":       c.Request.URL.Path,
			"method":     c.Request.Method,
			"clientIP":   c.ClientIP(),
			"dataLength": c.Writer.Size(),
			"latency":    stop.Seconds(),
		}
		msg := http.StatusText(status)
		if c.Request.URL.Path != "/" {
			if len(c.Errors) > 0 {
				fields["error"] = c.Errors.ByType(gin.ErrorTypePrivate)
				Logger.WithFields(fields).Error(msg)
			} else if status > 499 {
				Logger.WithFields(fields).Error(msg)
			} else if status > 399 {
				Logger.WithFields(fields).Warn(msg)
			} else if c.Request.URL.Path != "/" {
				Logger.WithFields(fields).Info(msg)
			}
		}
	}
}
