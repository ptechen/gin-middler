package middler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"time"
)

func Logger(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// before request
		c.Next()
		// after request
		latency := time.Since(t)
		status := c.Writer.Status()
		res, _ := c.Get("res")
		params, _ := c.Get("params")
		logger.Info().
			Str("handler", c.HandlerName()).
			Int("status", status).
			Dur("latency", latency).
			Str("method", c.Request.Method).
			Str("client_ip", c.ClientIP()).
			Interface("params", params).
			Interface("res", res).
			Send()
	}
}