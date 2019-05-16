# health

<a href="https://www.buymeacoffee.com/virtcanhead" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;" ></a>
[![Build Status](https://travis-ci.org/virtcanhead/health.svg?branch=master)](https://travis-ci.org/virtcanhead/health)

health check middleware for Echo web framework

## Usage

### Integration

```go
type redisResource struct {
  *redis.Client
}

func (r redisResource) HealthCheck() error {
  return r.Client.Ping().Err()
}

func main() {
  // ...

  e := echo.New()
  e.Use(health.New(redisResource{Client: redisClient}))

  // ...
}

```

### Check

For HTTP service, liveness and readiness should be identical.

```bash
curl http://127.0.0.1:1234/_health
```

This will check all resources with `HealthCheck()` method

## License

canhead <hi@canhead.xyz> MIT License
