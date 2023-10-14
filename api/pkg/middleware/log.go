package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vpaza/training/api/pkg/logger"
)

func ZapLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			logger.Log().Infof("%d %s %s %s %s %s %s %d %d (%s)",
				res.Status,
				time.Since(start).String(),
				id,
				req.Method,
				req.RequestURI,
				req.Host,
				c.RealIP(),
				req.ContentLength,
				res.Size,
				req.UserAgent(),
			)

			return nil
		}
	}
}
