package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	adh "github.com/adh-partnership/api/pkg/database/models"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/vpaza/training/api/pkg/config"
	"golang.org/x/oauth2"
)

func Routes(e *echo.Group) {
	e.GET("/login", getLogin)
	e.GET("/logout", getLogout)
	e.GET("/callback", getCallback)
}

// Login to Account
// @Summary Login to Account
// @Tags user, oauth
// @Param redirect query string false "Redirect URL"
// @Success 307
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/auth/login [get]
func getLogin(c echo.Context) error {
	state, _ := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 64)
	s, _ := session.Get(config.Cfg.Cookies.Name, c)
	verifier := oauth2.GenerateVerifier()
	s.Values["state"] = state
	s.Values["verifier"] = verifier
	s.Values["redirect"] = c.QueryParam("redirect")
	s.Save(c.Request(), c.Response())

	return c.Redirect(
		http.StatusTemporaryRedirect,
		config.Cfg.OAuth.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier)),
	)
}

// Logout of Account
// @Summary Logout of Account
// @Tags user, oauth
// @Success 307
// @Failure 500 {object} HTTPError
// @Router /v1/auth/logout [get]
func getLogout(c echo.Context) error {
	s, _ := session.Get(config.Cfg.Cookies.Name, c)
	s.Values["user_id"] = ""
	s.Save(c.Request(), c.Response())

	if c.QueryParam("redirect") == "" {
		return c.String(http.StatusNoContent, "")
	}

	return c.Redirect(
		http.StatusTemporaryRedirect,
		c.QueryParam("redirect"),
	)
}

type SSOUserResponse struct {
	Message string    `json:"message" yaml:"message" xml:"message"`
	User    *adh.User `json:"user" yaml:"user" xml:"user"`
}

// Callback from OAuth2 Provider
// @Summary Callback from OAuth2 Provider
// @Tags user, oauth
// @Param code query string true "OAuth2 Code"
// @Param state query string true "OAuth2 State"
// @Success 307
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/auth/callback [get]
func getCallback(c echo.Context) error {
	s, _ := session.Get(config.Cfg.Cookies.Name, c)
	s.Values["state"] = ""
	redirect := s.Values["redirect"].(string)
	s.Values["redirect"] = ""
	s.Values["verifier"] = ""
	s.Save(c.Request(), c.Response())
	if c.QueryParam("state") == "" || s.Values["state"] != c.QueryParam("state") {
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "Invalid state",
		}
	}

	token, err := config.Cfg.OAuth.Exchange(
		c.Request().Context(),
		c.QueryParam("code"),
		oauth2.VerifierOption(s.Values["verifier"].(string)),
	)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "failed_exchange",
		}
	}

	user, httperror := getVATSIMDetails(token.AccessToken)
	if httperror != nil {
		return err
	}

	if user == nil || user.User.ControllerType == "none" {
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "You are not a controller at ZAN",
		}
	}

	s.Values["user_id"] = fmt.Sprint(user.User.CID)
	s.Save(c.Request(), c.Response())

	if redirect == "" {
		return c.String(http.StatusOK, "Logged in.")
	}

	return c.Redirect(
		http.StatusTemporaryRedirect,
		redirect,
	)
}

func getVATSIMDetails(token string) (*SSOUserResponse, *echo.HTTPError) {
	res, err := http.NewRequest("GET", config.Cfg.OAuthUserInfo, nil)
	res.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	res.Header.Add("Accept", "application/json")
	res.Header.Add("User-Agent", "zan/training-api")
	if err != nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "Internal Server Error",
			Internal: err,
		}
	}

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "Internal Server Error",
			Internal: err,
		}
	}
	defer resp.Body.Close()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "Internal Server Error",
			Internal: err,
		}
	}

	if resp.StatusCode > 299 {
		return nil, &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "failed_exchange",
		}
	}

	data := &SSOUserResponse{}
	if err := json.Unmarshal(contents, &data); err != nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "Internal Server Error",
			Internal: err,
		}
	}

	return data, nil
}
