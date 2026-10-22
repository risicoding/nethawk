package echo

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/risicoding/nethawk/internal/logger"
)

func LoggerMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(
		middleware.RequestLoggerConfig{
			LogMethod:        true,
			LogURIPath:       true,
			LogStatus:        true,
			LogLatency:       true,
			LogRemoteIP:      true,
			LogContentLength: true,
			LogResponseSize:  true,
			HandleError:      true,

			LogValuesFunc: func(
				c *echo.Context,
				v middleware.RequestLoggerValues,
			) error {
				fields := []any{
					"method", v.Method,
					"path", v.URIPath,
					"status", v.Status,
					"latency", v.Latency.String(),
					"remote_ip", v.RemoteIP,
					"content_length", v.ContentLength,
					"response_size", v.ResponseSize,
				}

				if v.Error != nil {
					logger.Error("HTTP request failed",
						append(fields, "error", v.Error.Error())...,
					)
					return nil
				}

				logger.Info("HTTP request", fields...)
				return nil
			},
		},
	)
}
