package main

import (
	"log"

	"github.com/davidandw190/RESTful-api-go/db"
	"github.com/gofiber/fiber/v2"
)

func home(c *fiber.Ctx) error {
	return c.SendString("You are in home!")
}

func main() {

	db.ConnectDb()

	app := fiber.New()

	app.Get("/api", home)

	log.Fatal(app.Listen(":4001"))

}
