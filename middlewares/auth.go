package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"os"
	"strings"
)

func CheckAuth(c *fiber.Ctx) error {
	if c.Get("Authorization") == "" {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Token is required",
		})
	}

	parsToken := strings.SplitAfter(c.Get("Authorization"), " ")

	_, err := jwt.Parse(parsToken[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ACCESS_TOKEN_SECRETE")), nil
	})

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	c.Next();
	return nil

}
