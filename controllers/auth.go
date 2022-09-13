package controllers

import (
	//go builtin imports
	"context"
	"fmt"
	"os"
	"performance-evaluation-app/helper"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	//App imports
	"performance-evaluation-app/configs"
	"performance-evaluation-app/models"
)

func Signup(c *fiber.Ctx) error {
	userCollection := configs.MI.DB.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	user := new(models.Users)

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
			"message": "Field contains whitespaces",
		})
	}

	//Multiple fields validation
	error := validation.Errors{
		"name":     validation.Validate(user.Name, validation.Required),
		"email":    validation.Validate(user.Email, validation.Required, is.Email),
		"password": validation.Validate(user.Password, validation.Required, validation.Length(4, 12)),
	}.Filter()

	if error != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": error.Error(),
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
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
	userCollection := configs.MI.DB.Collection("users")

	var input map[string]string
	c.BodyParser(&input)
	var user models.Users

	error := validation.Errors{
		"email":    validation.Validate(input["email"], validation.Required, is.Email),
		"password": validation.Validate(input["password"], validation.Required, validation.Length(4, 12)),
	}.Filter()

	if error != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": error.Error(),
		})
	}

	result := userCollection.FindOne(context.Background(), bson.M{"email": input["email"]})

	if err := result.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User Not found",
			"error":   err.Error(),
		})
	}
	err := result.Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Error occurred while decoding response",
			"error":   err.Error(),
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input["password"]))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Entered wrong password",
			//"error":   err.Error(),
		})
	}

	//************** Access Token
	claims := &jwt.RegisteredClaims{
		Issuer:    user.Name,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(25 * time.Second)),
		ID:        user.ID.Hex(),
	}

	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, error := tokenString.SignedString([]byte(os.Getenv("JWT_ACCESS_TOKEN_SECRETE")))
	if error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   error.Error(),
		})
	}

	//************** Refresh Token
	refreshTokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    user.Name,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		ID:        user.ID.Hex(),
	})
	rToken, error := refreshTokenString.SignedString([]byte(os.Getenv("JWT_REFRESH_TOKEN_SECRETE")))
	if error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   error.Error(),
		})
	}

	//Now storing refresh token in redis cache
	err = helper.SetExVal("Issuer", user.Name, 12*time.Hour)
	fmt.Println(err)
	err = helper.SetExVal("ID", user.ID.Hex(), 12*time.Hour)
	err = helper.SetExVal("refresh_token", rToken, 12*time.Hour)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  token,
		"refresh_token": rToken,
		"exp":           claims.ExpiresAt,
	})

}

func GetNewAccessToken(c *fiber.Ctx) error {
	var input map[string]string
	c.BodyParser(&input)

	err := validation.Validate(input["refresh_token"], validation.Required)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	//validate refresh token with redis if match user is verified and assign new access token
	if input["refresh_token"] != helper.GetExVal("refresh_token") {
		return c.Status(401).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid refresh token provided",
		})
	}
	//************** Access Token
	IssuerCache := helper.GetExVal("Issuer")
	IDCache := helper.GetExVal("ID")
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		Issuer:    IssuerCache,
		ID:        IDCache,
	}

	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, error := tokenString.SignedString([]byte(os.Getenv("JWT_ACCESS_TOKEN_SECRETE")))
	if error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": token,
	})

}
func Logout(c *fiber.Ctx) error {

	return c.SendString("Logout function calling")

}
