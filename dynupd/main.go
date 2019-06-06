package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/virtcanhead/health"
	"os"
)

var (
	optBind     = os.Getenv("BIND")
	optRedisURL = os.Getenv("REDIS_URL")
)

var (
	storage Storage
)

func init() {
	if len(optRedisURL) == 0 {
		optRedisURL = "redis://127.0.0.1:6379"
	}
	if len(optBind) == 0 {
		optBind = ":8080"
	}
}

func main() {
	storage = NewRedisStorage(optRedisURL)

	e := echo.New()
	e.Debug = len(os.Getenv("DEBUG")) > 0
	e.HidePort = true
	e.HideBanner = true
	e.Use(health.New(storage))
	e.Static("/dynup", "public/dynup")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	mountRoutes(e)

	e.Logger.Fatal(e.Start(optBind))
}
