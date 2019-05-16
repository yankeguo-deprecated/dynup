package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Resource something that need be health checked
type Resource interface {
	HealthCheck() error // returns error if health check no passed
}

// New create a health check middleware
func New(rs ...Resource) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == "/_health" {
				for _, r := range rs {
					if err := r.HealthCheck(); err != nil {
						return c.String(http.StatusInternalServerError, err.Error())
					}
				}
				return c.String(http.StatusOK, "OK")
			}

			return next(c)
		}
	}
}
