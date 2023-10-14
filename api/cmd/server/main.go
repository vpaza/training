package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/http2"

	_ "github.com/vpaza/training/api/docs"
	"github.com/vpaza/training/api/internal/routes"
	"github.com/vpaza/training/api/pkg/config"
	"github.com/vpaza/training/api/pkg/database"
	"github.com/vpaza/training/api/pkg/logger"
	trainingmiddleware "github.com/vpaza/training/api/pkg/middleware"
	"github.com/vpaza/training/api/pkg/models"
	"github.com/vpaza/training/api/pkg/validator"
)

// @title ZAN Training API
// @version 1.0
// @description ZAN Training API
//
// @contact.name Daniel Hawton
// @contact.email daniel@hawton.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host training.zanartcc.org
// @BasePath /

// @securityDefinitions.apiKey XBearer
// @in header
// @name Authorization
// @description Enter token with the `Bearer ` prefix
func main() {
	_, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger.InitWithOptions(
		logger.WithConfigLevel(config.Cfg.LogLevel),
	)
	defer logger.Log().Sync()
	logger.Log().Infof("Setting up Training API...")
	logger.Log().Infof("Log Level set to %s", config.Cfg.LogLevel)

	logger.Log().Debugf("Config=%+v", config.Cfg)

	logger.Log().Info("Setting up server")
	e := echo.New()

	logger.Log().Infof("Setting up middleware")
	e.Use(trainingmiddleware.ZapLogger())
	e.Use(session.Middleware(
		sessions.NewCookieStore([]byte(config.Cfg.Cookies.Secret)),
	))
	e.Use(trainingmiddleware.Auth)
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(_ string) (bool, error) {
			return true, nil
		},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.Use(middleware.MethodOverride())
	e.Use(middleware.RemoveTrailingSlash())
	generateSecureMiddleware(e)
	e.HideBanner = true
	e.Validator = validator.Get()
	e.HidePort = true

	logger.Log().Infof("Setting up database")
	err = database.Connect(
		&database.DBOptions{
			Driver:   config.Cfg.Database.Driver,
			Host:     config.Cfg.Database.Hostname,
			Port:     fmt.Sprint(config.Cfg.Database.Port),
			User:     config.Cfg.Database.Username,
			Password: config.Cfg.Database.Password,
			Database: config.Cfg.Database.Database,
		},
	)
	if err != nil {
		logger.Log().Fatal(err)
	}

	logger.Log().Infof("Running migrations")
	err = database.DB.AutoMigrate(
		&models.Request{},
	)
	if err != nil {
		logger.Log().Fatal(err)
	}

	logger.Log().Infof("Setting up routes")
	routes.RegisterRoutes(e)

	logger.Log().Infof("Starting server on port %d", config.Cfg.ListenPort)
	go func() {
		switch config.Cfg.Mode {
		case "tls":
			err := e.StartTLS(
				fmt.Sprintf(":%d", config.Cfg.ListenPort),
				config.Cfg.TLSCert,
				config.Cfg.TLSKey,
			)
			if err != nil {
				logger.Log().Fatal(err)
			}
		case "h2c":
			sh2 := &http2.Server{
				MaxConcurrentStreams: 250,
				MaxReadFrameSize:     1048576,
				IdleTimeout:          10 * time.Second,
			}
			err := e.StartH2CServer(fmt.Sprintf(":%d", config.Cfg.ListenPort), sh2)
			if err != nil {
				logger.Log().Fatal(err)
			}
		case "plain":
			err := e.Start(fmt.Sprintf(":%d", config.Cfg.ListenPort))
			if err != nil {
				logger.Log().Fatal(err)
			}
		default:
			logger.Log().Fatal("Invalid mode")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Log().Infof("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Log().Fatal(err)
	}
}

func generateSecureMiddleware(e *echo.Echo) {
	c := &middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
	}

	if config.Cfg.Mode != "plain" {
		c.HSTSExcludeSubdomains = false
		c.HSTSMaxAge = 3600
	}

	e.Use(middleware.SecureWithConfig(*c))
}
