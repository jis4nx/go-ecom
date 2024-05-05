package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *ProductApp) LoadRoutes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route("/products", app.loadProductRoutes)
	return router
}

func (a *ProductApp) loadProductRoutes(router chi.Router) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
  router.Post("/create", a.CreateProduct)
}
