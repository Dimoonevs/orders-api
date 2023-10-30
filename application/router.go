package application

import (
	"context"
	"fmt"
	"github.com/Dimoonevs/orders-api.git/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func loadRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/orders", loadOrderRouter)

	return router
}

func loadOrderRouter(router chi.Router) {
	orderHandler := &handler.Order{}

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
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("ListenAndServe: %v", err)
	}
	return nil

}
