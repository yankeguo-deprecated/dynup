package main

import (
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/virtcanhead/health"
	"os"
)

var (
	bind       = os.Getenv("BIND")
	addrRedis  = os.Getenv("REDIS_ADDR")
	addrConsul = os.Getenv("CONSUL_ADDR")
)

func init() {
	if len(addrRedis) == 0 {
		addrRedis = "127.0.0.1:6379"
	}
	if len(addrConsul) == 0 {
		addrConsul = "127.0.0.1:8500"
	}
	if len(bind) == 0 {
		bind = ":9080"
	}
}

type redisResource struct {
	*redis.Client
}

func (r redisResource) HealthCheck() error {
	return r.Client.Ping().Err()
}

func main() {
	client := redis.NewClient(&redis.Options{Addr: addrRedis,})

	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Use(health.New(redisResource{client}))
	e.Static("/", "public")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(bind))
}
