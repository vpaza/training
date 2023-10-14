package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/vpaza/training/api/internal/routes/v1/auth"
)

var routeGroups map[string]map[string]func(*echo.Group)

func init() {
	routeGroups = make(map[string]map[string]func(*echo.Group))
	routeGroups["/v1"] = map[string]func(*echo.Group){
		"/auth": auth.Routes,
	}
}

func RegisterRoutes(e *echo.Echo) {
	for prefix, group := range routeGroups {
		g := e.Group(prefix)
		for path, fn := range group {
			fn(g.Group(path))
		}
	}
}
