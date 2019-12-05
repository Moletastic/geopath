package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/swaggo/echo-swagger"
)

// New crea un nueva instancia de Echo
func New() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} ${status}\n",
	}))
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	return e
}
