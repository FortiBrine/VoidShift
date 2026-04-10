package logger

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v5"
)

func RequestLogger(base *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			reqID := c.Response().Header().Get(echo.HeaderXRequestID)
			reqLogger := base.With(
				slog.String("request_id", reqID),
			)

			err := next(c)

			resp, unwrapErr := echo.UnwrapResponse(c.Response())
			if unwrapErr != nil {
				reqLogger.Error("failed to unwrap response", "error", unwrapErr)
				return err
			}

			status := resp.Status
			latency := time.Since(start)

			reqLogger.Info("http_request",
				slog.String("method", c.Request().Method),
				slog.String("path", c.Path()),
				slog.String("uri", c.Request().RequestURI),
				slog.Int("status", status),
				slog.Duration("latency", latency),
				slog.String("remote_ip", c.RealIP()),
				slog.Int64("bytes_out", resp.Size),
			)

			return err
		}
	}
}
