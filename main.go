package main

import (
	"go_api/auth"
	"go_api/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	defer database.Disconnect()

	app := fiber.New()
	app.Static("/", "./public/")
	app.Static("/auth/", "./public/auth")

	api_auth := app.Group("/api/auth")
	api_auth.Post("/register", auth.Register)
	api_auth.Post("/login", auth.Login)
	api_auth.Post("/logout", auth.Logout)
	api_auth.Get("/get-user", auth.GetUser)

	app.Listen(":8000")
}
