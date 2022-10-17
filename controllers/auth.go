package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"performance-evaluation-app-with-gin/configs"
	"performance-evaluation-app-with-gin/helper"
	"performance-evaluation-app-with-gin/models"
	"time"
)

func Signup(c *gin.Context) {
	userCollection := configs.MI.DB.Collection("users")
	var user models.Users

	c.ShouldBindJSON(&user)

	//Multiple fields validation
	err2 := validation.Errors{
		"name":     validation.Validate(user.Name, validation.Required),
		"email":    validation.Validate(user.Email, validation.Required, is.Email),
		"password": validation.Validate(user.Password, validation.Required, validation.Length(4, 12)),
	}.Filter()

	if err2 != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err2.Error(),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Password = string(hashedPassword)

	//************ DB Query **********/
	result, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "User registration failed.",
			"err":     err,
		})
		return
	}

	c.JSON(200,
		gin.H{"data": result,
			"success": true,
			"message": "user inserted successfully",
		})

}

func Login(c *gin.Context) {

	var input map[string]string
	var user models.Users

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	//Multiple fields validation
	error := validation.Errors{
		"email":    validation.Validate(input["email"], validation.Required, is.Email),
		"password": validation.Validate(input["password"], validation.Required, validation.Length(4, 12)),
	}.Filter()

	if error != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": error.Error(),
		})
		return
	}

	userCollection := configs.MI.DB.Collection("users")
	result := userCollection.FindOne(context.Background(), bson.M{"email": input["email"]})

	if err := result.Err(); err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "User Not found",
			"error":   err.Error(),
		})
		return
	}

	err := result.Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "Error occurred while decoding response",
			"error":   err.Error(),
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input["password"]))
	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "Entered wrong password",
			//"error":   err.Error(),
		})
		return

	}

	//************** Access Token
	claims := jwt.MapClaims{}
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(1 * time.Minute))
	claims["issuer"] = user.Name
	claims["user_id"] = user.ID.Hex()
	if user.UserRole != nil {
		claims["role_id"] = user.UserRole.Hex()
	}

	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, error := tokenString.SignedString([]byte(os.Getenv("JWT_ACCESS_TOKEN_SECRETE")))
	if error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   error.Error(),
		})
		return
	}

	//************** Refresh Token
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(12 * time.Hour))
	refreshTokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	rToken, error := refreshTokenString.SignedString([]byte(os.Getenv("JWT_REFRESH_TOKEN_SECRETE")))
	if error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   error.Error(),
		})
		return
	}

	if os.Getenv("IS_REDIS_CACHE_ENABLED") == "true" {
		//Now storing refresh token in redis cache
		err = helper.SetExVal("Issuer", user.Name, 12*time.Hour)
		err = helper.SetExVal("ID", user.ID.Hex(), 12*time.Hour)
		err = helper.SetExVal("refresh_token", rToken, 12*time.Hour)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	} else {
		//db to store refresh token
		filter := bson.M{"_id": user.ID}
		fields := bson.M{"$set": bson.M{"refresh_token": rToken, "updated_at": time.Now()}}
		_, err := userCollection.UpdateOne(context.Background(), filter, fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to updated refresh token in db",
				"error":   err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  token,
		"refresh_token": rToken,
		"exp":           jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
	})

}

func GetNewAccessToken(c *gin.Context) {
	var user models.Users
	var input map[string]string
	c.ShouldBindJSON(&input)

	err := validation.Validate(input["refresh_token"], validation.Required)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	var IssuerCache = ""
	var IDCache = ""
	if os.Getenv("IS_REDIS_CACHE_ENABLED") == "true" {
		//validate refresh token with redis if match user is verified and assign new access token
		if input["refresh_token"] != helper.GetExVal("refresh_token") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Invalid refresh token provided",
			})
			return
		}

		//************** Access Token
		IssuerCache = helper.GetExVal("Issuer")
		IDCache = helper.GetExVal("ID")

	} else {
		userCollection := configs.MI.DB.Collection("users")
		result := userCollection.FindOne(context.Background(), bson.M{"refresh_token": input["refresh_token"]})

		if err := result.Err(); err != nil {
			c.JSON(404, gin.H{
				"success": false,
				"message": "User Not found",
				"error":   err.Error(),
			})
			return
		}

		result.Decode(&user)
		if input["refresh_token"] != user.RefreshToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Invalid refresh token provided",
			})
			return
		}

		//************** Access Token
		IssuerCache = user.Name
		IDCache = user.ID.Hex()

	}

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		Issuer:    IssuerCache,
		ID:        IDCache,
	}

	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, error := tokenString.SignedString([]byte(os.Getenv("JWT_ACCESS_TOKEN_SECRETE")))
	if error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})

}

func Logout(c *gin.Context) {
	c.JSON(200, "Logout function")
}
