package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/davidandw190/RESTful-api-go/db"
	"github.com/davidandw190/RESTful-api-go/models"
	"github.com/gofiber/fiber/v2"
)

// OrderSerializer is a struct that defines the JSON representation of an order.
type OrderSerializer struct {
	ID        uint64            `json:"id"`
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
	if err := findUserByID(order.UserRefer, &user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var product models.Product
	if err := findProductByID(order.ProductRefer, &product); err != nil {
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
	paramID := c.Params("id")
	userID, err := strconv.ParseUint(paramID, 10, 64)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	var user models.User
	if err := findUserByID(userID, &user); err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	var orders []models.Order
	if err := db.Database.Db.Where("user_referer = ?", userID).Find(&orders).Error; err != nil {
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

func GetUserOrder(c *fiber.Ctx) error {
	paramID := c.Params("user_id")
	userID, err := strconv.ParseUint(paramID, 10, 64)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	paramID = c.Params("order_id")
	orderID, err := strconv.ParseUint(paramID, 10, 64)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalID request"})
	}

	var user models.User
	if err := findUserByID(userID, &user); err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	var order models.Order
	if err := db.Database.Db.Where("id = ? AND user_referer = ?", orderID, userID).Find(&order).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Order not found for the specified user"})
	}

	var product models.Product
	if err := findProductByID(order.ProductRefer, &product); err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	responseUser := CreateResponseUser(&user)
	responseProduct := CreateResponseProduct(&product)
	responseOrder := CreateResponseOrder(&order, responseUser, responseProduct)

	return c.Status(http.StatusOK).JSON(responseOrder)
}
