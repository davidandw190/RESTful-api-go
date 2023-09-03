package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/davidandw190/RESTful-api-go/db"
	"github.com/davidandw190/RESTful-api-go/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ProductSerializer is a struct used to serialize product data to JSON format.
type ProductSerializer struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_num"`
}

// CreateResponseProduct creates a ProductSerializer from a given product model.
func CreateResponseProduct(productModel *models.Product) ProductSerializer {
	return ProductSerializer{
		ID:           productModel.ID,
		Name:         productModel.Name,
		SerialNumber: productModel.SerialNumber,
	}
}

// CreateProduct handles the creation of a new product.
func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := db.Database.Db.Create(product).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
	}

	responseProduct := CreateResponseProduct(&product)
	return c.Status(http.StatusCreated).JSON(responseProduct)
}

// GetAllProducts retrieves all products from the database and returns them as JSON.
func GetAllProducts(c *fiber.Ctx) error {
	products := []models.Product{}

	if err := db.Database.Db.Find(&products); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve products"})
	}

	responseProducts := make([]ProductSerializer, len(products))
	for i, productModel := range products {
		responseProducts[i] = CreateResponseProduct(&productModel)
	}

	return c.Status(http.StatusOK).JSON(responseProducts)
}

// GetProduct retrieves a product by its ID and returns it as JSON.
func GetProduct(c *fiber.Ctx) error {
	paramID := c.Params("id")
	id, err := strconv.ParseUint(paramID, 10, 64)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var product models.Product

	if err := findProductByID(id, &product); err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	responseProduct := CreateResponseProduct(&product)
	return c.Status(http.StatusAccepted).JSON(responseProduct)
}

// UpdateProduct updates a product's information based on the provided JSON data.
func UpdateProduct(c *fiber.Ctx) error {
	paramID := c.Params("id")
	id, err := strconv.ParseUint(paramID, 10, 64)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var product models.Product

	if err := findProductByID(id, &product); err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	// to hold the updated data from the request
	type UpdatedProduct struct {
		Name         string `json:"name"`
		SerialNumber string `json:"serial_num"`
	}

	var updatedData UpdatedProduct

	// to the request body into updatedData
	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	// updating the product's data
	product.Name = updatedData.Name
	product.SerialNumber = updatedData.SerialNumber

	// saving the updated product to the database
	if err := db.Database.Db.Save(&product).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product"})
	}

	// serializing and returning the updated product as JSON
	responseProduct := CreateResponseProduct(&product)
	return c.Status(http.StatusAccepted).JSON(responseProduct)
}

func DeleteProduct(c *fiber.Ctx) error {
	paramID := c.Params("id")
	id, err := strconv.ParseUint(paramID, 10, 64)

	var product models.Product

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := findProductByID(id, &product); err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	if err := db.Database.Db.Delete(&models.User{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}

func findProductByID(id uint64, product *models.Product) error {
	db.Database.Db.Find(&product, "id = ?", id)

	if product.ID == 0 {
		return errors.New("Product does not exist")
	}

	return nil
}
