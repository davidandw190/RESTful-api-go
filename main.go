package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func home(c *fiber.Ctx) error {
	return c.SendString("You are in home!")
}

func main() {
	app := fiber.New()

	app.Get("/api", home)

	log.Fatal(app.Listen(":3000"))

}
