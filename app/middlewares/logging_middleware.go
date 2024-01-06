package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		logger.WithFields(logrus.Fields{
			"method":  c.Request.Method,
			"path":    c.Request.RequestURI,
			"status":  c.Writer.Status(),
			"latency": time.Since(startTime),
			"ip":      c.ClientIP(),
		}).Info("Handled request")
	}
}
