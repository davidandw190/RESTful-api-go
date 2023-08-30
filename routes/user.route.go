package routes

import (
	"errors"
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

// GetAllUsers returns a slice of serialised users
func GetAllUsers(c *fiber.Ctx) error {
	users := []models.User{}

	db.Database.Db.Find(&users)
	respUsers := make([]UserSerializer, len(users), 0)

	for _, user := range users {
		rUser := CreateResponseUser(user)
		respUsers = append(respUsers, rUser)

	}

	return c.Status(http.StatusOK).JSON(respUsers)
}

// GetUser returns a user in serialised form by id
func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error)
	}

	if err := findUser(id, &user); err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}

	responseUser := CreateResponseUser(user)

	return c.Status(http.StatusAccepted).JSON(responseUser)

}

// UpdateUser updates a user entry and returns it in serialised form.
func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	if err := findUser(id, &user); err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}

	type UpdateUser struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	var updatedData UpdateUser

	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	user.FirstName = updatedData.FirstName
	user.LastName = updatedData.LastName
	user.Email = updatedData.Email

	db.Database.Db.Save(&user)

	responseUser := CreateResponseUser(user)

	return c.Status(http.StatusAccepted).JSON(responseUser)

}

func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	if err := findUser(id, &user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	if err := db.Database.Db.Delete(&user).Error; err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	return c.Status(http.StatusOK).SendString("User deleted successfully!")

}

func findUser(id int, user *models.User) error {
	db.Database.Db.Find(&user, "id=?", id)

	if user.ID == 0 {
		return errors.New("User does not exist")
	}

	return nil

}