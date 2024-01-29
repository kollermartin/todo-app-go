package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		logrus.WithFields(logrus.Fields{
			"method":  c.Request.Method,
			"path":    c.Request.RequestURI,
			"status":  c.Writer.Status(),
			"latency": time.Since(startTime),
			"ip":      c.ClientIP(),
		}).Info("Handled request")
	}
}