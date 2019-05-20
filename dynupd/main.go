package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/virtcanhead/health"
	"os"
)

var (
	bind      = os.Getenv("BIND")
	redisAddr = os.Getenv("REDIS_ADDR")
)

var (
	storage Storage
)

func init() {
	if len(redisAddr) == 0 {
		redisAddr = "127.0.0.1:6379"
	}
	if len(bind) == 0 {
		bind = ":8080"
	}
}

func main() {
	storage = NewRedisStorage(redisAddr)

	e := echo.New()
	e.Debug = len(os.Getenv("DEBUG")) > 0
	e.HidePort = true
	e.HideBanner = true
	e.Use(health.New(storage))
	e.Static("/", "public")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	mountRoutes(e)

	e.Logger.Fatal(e.Start(bind))
}
