package http

import (
	"log"
	"sika/api/http/handlers"
	"sika/config"
	"sika/service"

	"github.com/gofiber/fiber/v2"
)

func Run(cfg config.Config, app *service.AppContainer) {
	fiberApp := fiber.New()
	app, err := service.NewAppContainer(cfg)
	if err != nil {
		panic(err)
	}
	fiberApp.Get("/users/:UserID", handlers.GetUserByID(app.UserService()))
	log.Fatal(fiberApp.Listen("localhost:8080"))
}
