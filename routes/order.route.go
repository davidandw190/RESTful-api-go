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
	ID        uint64              `json:"id"`
	User      UserSerializer      `json:"user"`
	Products  []ProductSerializer `json:"products"`
	CreatedAt time.Time           `json:"order_date"`
}

// CreateResponseOrder creates a serialized representation of an order.
func CreateResponseOrder(order *models.Order, user UserSerializer, products []ProductSerializer) OrderSerializer {
	return OrderSerializer{
		ID:        order.ID,
		User:      user,
		Products:  products,
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

	// Parse a list of product IDs from the request body.
	var productIDs []uint64
	if err := c.BodyParser(&productIDs); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid product IDs"})
	}

	var products []models.Product
	for _, productID := range productIDs {
		var product models.Product
		if err := findProductByID(productID, &product); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		products = append(products, product)
	}

	// Associate the products with the order.
	order.Products = products
	db.Database.Db.Create(&order)

	responseUser := CreateResponseUser(&user)
	responseProducts := []ProductSerializer{}
	for _, product := range products {
		responseProducts = append(responseProducts, CreateResponseProduct(&product))
	}
	responseOrder := CreateResponseOrder(&order, responseUser, responseProducts)

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
		var products []ProductSerializer
		for _, product := range o.Products {
			products = append(products, CreateResponseProduct(&product))
		}
		rOrder := CreateResponseOrder(&o, responseUser, products)
		responseOrders = append(responseOrders, rOrder)
	}

	return c.Status(http.StatusOK).JSON(responseOrders)
}

// GetAllUserOrders retrieves a specific order associated with a user.
func GetUserOrder(c *fiber.Ctx) error {
	paramUserID := c.Params("user_id")
	userID, err := strconv.ParseUint(paramUserID, 10, 64)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	paramOrderID := c.Params("order_id")
	orderID, err := strconv.ParseUint(paramOrderID, 10, 64)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	var user models.User
	if err := findUserByID(userID, &user); err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	var order models.Order
	if err := db.Database.Db.Where("id = ? AND user_referer = ?", orderID, userID).Preload("Products").Find(&order).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Order not found for the specified user"})
	}

	var products []ProductSerializer
	for _, product := range order.Products {
		products = append(products, CreateResponseProduct(&product))
	}

	responseUser := CreateResponseUser(&user)
	responseOrder := CreateResponseOrder(&order, responseUser, products)

	return c.Status(http.StatusOK).JSON(responseOrder)
}
