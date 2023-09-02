package routes

import (
	"net/http"
	"time"

	"github.com/davidandw190/RESTful-api-go/db"
	"github.com/davidandw190/RESTful-api-go/models"
	"github.com/gofiber/fiber/v2"
)

type OrderSerializer struct {
	ID        uint              `json:"id"`
	User      UserSerializer    `json:"user"`
	Product   ProductSerializer `json:"product"`
	CreatedAt time.Time         `json:"order_date"`
}

func CreateResponseOrder(order *models.Order, user UserSerializer, product ProductSerializer) OrderSerializer {
	return OrderSerializer{
		ID:        order.ID,
		User:      user,
		Product:   product,
		CreatedAt: order.CreatedAt,
	}
}

func CreateOrder(c *fiber.Ctx) error {
	var order models.Order

	if err := c.BodyParser(&order); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	var user models.User
	if err := findUserById(int(order.UserReferer), &user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var product models.Product
	if err := findProductById(int(order.ProductRefer), &product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	db.Database.Db.Create(&order)

	responseUser := CreateResponseUser(&user)
	responseProduct := CreateResponseProduct(&product)
	responseOrder := CreateResponseOrder(&order, responseUser, responseProduct)

	return c.Status(http.StatusCreated).JSON(responseOrder)
}
