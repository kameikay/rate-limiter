package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kameikay/rate-limiter/ratelimiter"
)

func main() {
	err := ratelimiter.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	router.Use(ratelimiter.NewRateLimiter())
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}

}
