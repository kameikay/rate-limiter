package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/kameikay/rate-limiter/internal/infra/web/handlers"
)

type HelloController struct {
	router       chi.Router
	HelloHandler *handlers.HelloHandler
}

func NewHelloController(
	router chi.Router,
	HelloHandler *handlers.HelloHandler,
) *HelloController {
	return &HelloController{
		router:       router,
		HelloHandler: HelloHandler,
	}
}

func (uc *HelloController) Route() {
	uc.router.Route("/Hello", func(r chi.Router) {

		// r.With(middlewares.CheckRole([]string{"admin"})).Get("/", uc.HelloHandler.GetAllHelloHandler)
	})
}
