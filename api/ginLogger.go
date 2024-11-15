package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process the request
		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath()

		logger := log.Info()
		if statusCode >= http.StatusBadRequest {
			logger = log.Error()
		}

		logger.Str("protocol", "http").
			Str("method", method).
			Str("path", path).
			Int("status_code", statusCode).
			Dur("duration", duration).
			Msg("received an HTTP request")
	}
}
