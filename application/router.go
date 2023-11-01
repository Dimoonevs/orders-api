package application

import (
	"context"
	"fmt"
	"github.com/Dimoonevs/orders-api.git/handler"
	"github.com/Dimoonevs/orders-api.git/repository/order"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"time"
)

func (a *App) loadRouter() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/orders", a.loadOrderRouter)

	a.router = router
}

func (a *App) loadOrderRouter(router chi.Router) {
	orderHandler := &handler.Order{
		Repo: &order.RedisRepo{
			RedisClient: a.rdb,
		},
	}

	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.List)
	router.Get("/{id}", orderHandler.GetById)
	router.Put("/{id}", orderHandler.UpdateById)
	router.Delete("/{id}", orderHandler.DeleteById)
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("Ping: %v", err)
	}
	fmt.Println("Redis connected")

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("ListenAndServe: %v", err)
		}
		close(ch)
	}()
	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(timeout)
	}

	return nil

}
