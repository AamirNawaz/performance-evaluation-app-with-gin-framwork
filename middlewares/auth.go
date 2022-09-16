package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"performance-evaluation-app-with-gin/configs"
	"performance-evaluation-app-with-gin/models"
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

func CheckRole(c *gin.Context) {
	stringArray := c.Request.Header["Authorization"]
	justString := strings.Join(stringArray, " ")

	parsToken := strings.SplitAfter(justString, " ")

	tokenString := parsToken[1]
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ACCESS_TOKEN_SECRETE")), nil
	})
	// ... error handling
	if err != nil {
		fmt.Println(err)
	}

	var roleID string
	for key, val := range claims {
		//fmt.Printf("Key: %v, value: %v\n", key, val)
		if key == "role_id" {
			roleID = val.(string)
		}
	}

	var role models.Roles
	roleCollection := configs.MI.DB.Collection("roles")
	obId, _ := primitive.ObjectIDFromHex(roleID)
	query := bson.M{"_id": obId}

	roleCollection.FindOne(context.Background(), query).Decode(&role)
	if role.Name != "admin" {
		c.JSON(403, gin.H{
			"success": false,
			"message": "UnAuthorized needs admin access",
		})
		c.Abort()
		return
	}

	c.Next()
}
