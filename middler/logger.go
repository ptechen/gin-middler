package middler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"time"
)

//func MiddleLogger(logger *zerolog.Logger) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		t := time.Now()
//		// before request
//
//		c.Next()
//		// after request
//
//		latency := time.Since(t)
//		status := c.Writer.Status()
//		res, _ := c.Get("res")
//		params, _ := c.Get("params")
//		errNo, _ := c.Get("err_no")
//		if errNo == nil {
//			errNo = 0
//		}
//		err, _ := c.Get("err")
//
//		if err != nil {
//			logger.Error().
//				Str("router", c.Request.URL.String()).
//				Int("status", status).
//				Interface("err_no", errNo).
//				Dur("latency", latency).
//				Str("method", c.Request.Method).
//				Str("client_ip", c.ClientIP()).
//				Interface("params", params).
//				Interface("res", res).
//				Msgf("%#v", err)
//
//		} else {
//			logger.Info().
//				Str("router", c.Request.URL.String()).
//				Int("status", status).
//				Interface("err_no", errNo).
//				Dur("latency", latency).
//				Str("method", c.Request.Method).
//				Str("client_ip", c.ClientIP()).
//				Interface("params", params).
//				Interface("res", res).
//				Msg("success")
//		}
//	}
//}

func MiddleLogger(logger *zerolog.Logger) gin.HandlerFunc {
	return LoggerWithConfig(logger, gin.LoggerConfig{})
}

func LoggerWithConfig(logger *zerolog.Logger, conf gin.LoggerConfig) gin.HandlerFunc {
	notlogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			keys := c.Keys
			// Stop timer
			timeStamp := time.Now()
			latency := timeStamp.Sub(start)
			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()
			errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
			bodySize := c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			params, _ := c.Get("params")
			res, _ := c.Get("res")
			if errorMessage != "" {
				logger.Error().
					Str("router", c.Request.URL.String()).
					Interface("keys", keys).
					Int("status", statusCode).
					Dur("latency", latency).
					Str("method", method).
					Str("path", path).
					Str("client_ip", clientIP).
					Interface("params", params).
					Int("body_size", bodySize).
					Interface("skip", skip).
					Interface("res", res).
					Msg(errorMessage)

			} else {
				logger.Info().
					Str("router", c.Request.URL.String()).
					Interface("keys", keys).
					Int("status", statusCode).
					Dur("latency", latency).
					Str("method", method).
					Str("client_ip", clientIP).
					Interface("params", params).
					Int("body_size", bodySize).
					Interface("skip", skip).
					Interface("res", res).
					Msg("success")
			}
		}
	}
}
