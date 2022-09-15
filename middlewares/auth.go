package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"strings"
)

func CheckAuth(c *gin.Context) {

	stringArray := c.Request.Header["Authorization"]
	justString := strings.Join(stringArray, " ")

	parsToken := strings.SplitAfter(justString, " ")

	if parsToken[0] == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Token is required",
		})
	}

	if parsToken[1] == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Token is required",
		})
	}
	_, err := jwt.Parse(parsToken[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ACCESS_TOKEN_SECRETE")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})

		c.Abort()
		return

	}
	c.Next()

}
