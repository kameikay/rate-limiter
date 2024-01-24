package main

import (
	"github.com/kameikay/rate-limiter/configs"
	"github.com/kameikay/rate-limiter/internal/infra/redis"
	"github.com/kameikay/rate-limiter/internal/infra/web/controllers"
	"github.com/kameikay/rate-limiter/internal/infra/web/handlers"
	"github.com/kameikay/rate-limiter/internal/infra/web/webserver"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	rdb := redis.NewRedis()

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webserver.MountMiddlewares()
	// webserver.Router.Use(middleware.WithValue("RefreshTokenDuration", configs.RefreshTokenDuration))

	// HEALTH CHECK
	healthCheckController := controllers.NewHeathCheckController(webserver.Router)
	healthCheckController.Route()

	// HELLO
	helloHandler := handlers.NewHelloHandler(rdb)
	helloController := controllers.NewHelloController(webserver.Router, helloHandler)
	helloController.Route()

	// USERS
	// userRepo := userRepository.NewUserRepository(dbConnection)
	// confirmationCodeRepo := confirmationCodeRepository.NewConfirmationCodesRepository(dbConnection)
	// userHandler := handlers.NewUserHandler(uow, userRepo, confirmationCodeRepo, rdb)
	// userController := controllers.NewUsersController(webserver.Router, userHandler, configs.TokenAuth)
	// userController.Route()

	// STARTER
	webserver.Start()
}
