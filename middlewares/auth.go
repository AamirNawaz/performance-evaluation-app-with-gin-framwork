package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"os"
	"strings"
)

func SplitToken(headerToken string) string {
	parsToken := strings.SplitAfter(headerToken, " ")
	tokenString := parsToken[1]
	return tokenString
}

func CheckAuth(c *fiber.Ctx) error {
	_, err := jwt.Parse(SplitToken(c.Get("Authorization")), func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRETE")), nil
	})

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	} else {
		return nil
	}

}
