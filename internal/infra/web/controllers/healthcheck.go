package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/kameikay/rate-limiter/internal/infra/web/handlers"
)

type HelthCheckController struct {
	router chi.Router
}

func NewHeathCheckController(router chi.Router) *HelthCheckController {
	return &HelthCheckController{
		router: router,
	}
}

func (hc *HelthCheckController) Route() {
	hc.router.Route("/healthcheck", func(r chi.Router) {
		r.Get("/", handlers.HealthCheck)
	})
}
