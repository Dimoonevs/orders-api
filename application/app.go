package application

import (
	"github.com/go-redis/redis"
	"net/http"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {
	app := &App{
		router: loadRouter(),
		rdb:    redis.NewClient(&redis.Options{}),
	}
	return app
}
