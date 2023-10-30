package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func main() {

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/hello", basicHandler)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("ListenAndServe: %v", err)
	}
}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
