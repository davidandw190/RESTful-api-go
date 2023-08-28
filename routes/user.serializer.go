package routes

import (
	"net/http"

	"github.com/davidandw190/RESTful-api-go/db"
	"github.com/davidandw190/RESTful-api-go/models"
	"github.com/gofiber/fiber/v2"
)

// UserSerializer serializes a User model for response.
type UserSerializer struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// CreateResponseUser creates a UserSerializer from a User model.
func CreateResponseUser(userModel models.User) UserSerializer {
	return UserSerializer{
		ID:        userModel.ID,
		FirstName: userModel.FirstName,
		LastName:  userModel.LastName,
		Email:     userModel.Email,
	}
}

// CreateUser creates a new user.
func CreateUser(c *fiber.Ctx) error {
	var user models.User

	// Parse user data from the request body.
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	// Create the user in the database.
	db.Database.Db.Create(&user)
	responseUser := CreateResponseUser(user)

	return c.Status(http.StatusCreated).JSON(responseUser)
}
