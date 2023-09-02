package routes

import (
	"net/http"
	"time"

	"github.com/davidandw190/RESTful-api-go/db"
	"github.com/davidandw190/RESTful-api-go/models"
	"github.com/gofiber/fiber/v2"
)

// OrderSerializer is a struct that defines the JSON representation of an order.
type OrderSerializer struct {
	ID        uint              `json:"id"`
	User      UserSerializer    `json:"user"`
	Product   ProductSerializer `json:"product"`
	CreatedAt time.Time         `json:"order_date"`
}

// CreateResponseOrder creates a serialized representation of an order.
func CreateResponseOrder(order *models.Order, user UserSerializer, product ProductSerializer) OrderSerializer {
	return OrderSerializer{
		ID:        order.ID,
		User:      user,
		Product:   product,
		CreatedAt: order.CreatedAt,
	}
}

// CreateOrder handles the creation of a new order.
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

// GetAllUserOrders retrieves all orders associated with a specific user.
func GetAllUserOrders(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	var user models.User
	if err := findUserById(userId, &user); err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	var orders []models.Order
	if err := db.Database.Db.Where("user_referer = ?", userId).Find(&orders).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user orders"})
	}

	responseOrders := []OrderSerializer{}
	responseUser := CreateResponseUser(&user)
	for _, o := range orders {
		var product models.Product
		db.Database.Db.Find(&product, "id = ?", o.ProductRefer)
		rOrder := CreateResponseOrder(&o, responseUser, CreateResponseProduct(&product))
		responseOrders = append(responseOrders, rOrder)
	}

	return c.Status(http.StatusOK).JSON(responseOrders)
}
