package main

import (
	"log"

	"github.com/davidandw190/RESTful-api-go/db"
	"github.com/davidandw190/RESTful-api-go/routes"
	"github.com/gofiber/fiber/v2"
)

func home(c *fiber.Ctx) error {
	return c.SendString("You are in home!")
}

// setupRoutes initializes the API routes.
func setupRoutes(app *fiber.App) {
	app.Get("/api", home)

	app.Post("/api/users", routes.CreateUser)
	app.Get("/api/users", routes.GetAllUsers)
	app.Get("/api/users/:id", routes.GetUser)
	app.Put("/api/users/:id", routes.UpdateUser)
	app.Delete("/api/users/:id", routes.DeleteUser)
}

func main() {
	// Connect to the database.
	db.ConnectDb()

	app := fiber.New()

	// Set up API routes.
	setupRoutes(app)

	// Start the server.
	log.Fatal(app.Listen(":4001"))

}
