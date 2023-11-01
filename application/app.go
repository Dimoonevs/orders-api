package application

import (
	"github.com/redis/go-redis/v9"
	"net/http"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{}),
	}
	app.loadRouter()
	return app
}
