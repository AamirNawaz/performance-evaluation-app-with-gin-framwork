package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"performance_valuation_app/configs"
	"performance_valuation_app/models"
	"strings"
	"time"
)

func Signup(c *fiber.Ctx) error {
	userCollection := configs.MI.DB.Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}
	
	//******************  Validation ***************/
	if user.Name == "" || strings.TrimSpace(user.Name) == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Name is required",
		})
	}

	if user.Email == "" || strings.TrimSpace(user.Email) == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Email is required",
		})
	}

	if user.Password == "" || strings.TrimSpace(user.Password) == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Password is required",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	// Comparing the password with the hash
	//err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("aamir"))
	//fmt.Println("matching password:::::", err)

	user.Password = string(hashedPassword)
	//************ DB Query **********/
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "User registration failed.",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "user inserted successfully",
	})

	return nil
}

// Login get user and password
func Login(c *fiber.Ctx) error {
	var data map[string]string
	c.BodyParser(&data)

	return c.JSON(fiber.Map{
		"data": data,
	})
}

func Logout(c *fiber.Ctx) error {
	return c.SendString("Logout function calling")

}
