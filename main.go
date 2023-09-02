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
	// home endponts
	app.Get("/api", home)

	// user endpoints
	app.Post("/api/users", routes.CreateUser)
	app.Get("/api/users", routes.GetAllUsers)
	app.Get("/api/users/:id", routes.GetUser)
	app.Put("/api/users/:id", routes.UpdateUser)
	app.Delete("/api/users/:id", routes.DeleteUser)

	// product endpoints
	app.Post("/api/products", routes.CreateProduct)
	app.Get("/api/products", routes.GetAllProducts)
	app.Get("/api/products/:id", routes.GetProduct)
	app.Put("/api/products/:id", routes.UpdateProduct)
	app.Delete("api/products/:id", routes.DeleteProduct)

	// order endpoints
	app.Post("/api/orders", routes.CreateOrder)
	app.Get("/api/users/:id/orders", routes.GetAllUserOrders)

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
