package main

import (
	"go_api/auth"
	"go_api/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	// database.AutoMigrate()
	defer database.Disconnect()

	app := fiber.New()
	app.Static("/", "./public/")
	app.Static("/auth/", "./public/auth/")
	app.Static("/home/", "./public/home/")

	api_auth := app.Group("/api/auth")
	api_auth.Post("/register", auth.Register)
	api_auth.Post("/login", auth.Login)
	api_auth.Post("/logout", auth.Logout)
	api_auth.Post("/get-user", auth.GetUser)

	api_user := app.Group("/api/user")
	api_user.Post("/money", auth.Register)
	api_user.Post("/addupgrade", auth.Register)

	app.Listen(":8000")
}
