package middleware

import (
	"net/http"
	"time"

	adh "github.com/adh-partnership/api/pkg/database/models"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/vpaza/training/api/pkg/config"
	"github.com/vpaza/training/api/pkg/models"
)

type CustomContext struct {
	echo.Context

	XAuth   bool
	XUserID string
	XUser   *adh.User
}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &CustomContext{Context: c}

		sess, _ := session.Get(config.Cfg.Cookies.Name, cc)
		sess.Options = &sessions.Options{
			Path:     config.Cfg.Cookies.Path,
			MaxAge:   config.Cfg.Cookies.MaxAge,
			Domain:   config.Cfg.Cookies.Domain,
			Secure:   config.Cfg.Cookies.Secure,
			HttpOnly: config.Cfg.Cookies.HttpOnly,
		}
		userID, ok := sess.Values["user_id"].(string)
		if ok {
			cc.XAuth = true
			cc.XUserID = userID
			u, err := models.FindUser(userID)
			if err != nil {
				return &echo.HTTPError{
					Code:     http.StatusInternalServerError,
					Message:  "Failed to find user",
					Internal: err,
				}
			}
			cc.XUser = u
		} else {
			cc.XAuth = false
		}

		sess.Values["updated"] = time.Now()
		sess.Save(cc.Request(), cc.Response())

		return next(cc)
	}
}
