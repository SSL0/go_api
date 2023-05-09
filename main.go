package main

import (
	"go_api/auth"
	"go_api/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	api := app.Group("/auth")
	api.Post("/register", auth.Register)
	api.Post("/login", auth.Login)
	api.Post("/logout", auth.Logout)
	api.Post("/get-user", auth.GetUser)

	app.Listen(":8000")

}
