package routes

import (
	"errors"
	"net/http"

	"github.com/davidandw190/RESTful-api-go/db"
	"github.com/davidandw190/RESTful-api-go/models"
	"github.com/gofiber/fiber/v2"
)

type ProductSerializer struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_num"`
}

func CreateResponseProduct(productModel *models.Product) ProductSerializer {
	return ProductSerializer{
		ID:           productModel.ID,
		Name:         productModel.Name,
		SerialNumber: productModel.SerialNumber,
	}
}

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

func GetProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var product models.Product

	if err := findProductById(id, &product); err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	responseProduct := CreateResponseProduct(&product)

	return c.Status(http.StatusAccepted).JSON(responseProduct)

}

func findProductById(id int, product *models.Product) error {
	db.Database.Db.Find(&product, "id=?", id)

	if product.ID == 0 {
		return errors.New("Product does not exist")
	}

	return nil
}
