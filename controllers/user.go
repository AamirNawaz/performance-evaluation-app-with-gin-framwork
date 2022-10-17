package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"performance-evaluation-app-with-gin/configs"
	"performance-evaluation-app-with-gin/models"
	"strconv"
	"time"
)

func GetUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))

	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))

	//with projection, we can skip specific field to retrieve
	opt.SetProjection(bson.D{{"password", 0}, {"refresh_token", 0}})

	query := bson.M{}

	userCollection := configs.MI.DB.Collection("users")
	result, err := userCollection.Find(context.Background(), query, &opt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	defer result.Close(context.Background())
	var users []models.UsersResponse

	err = result.All(context.TODO(), &users)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   users,
	})

}

func GetUserById(c *gin.Context) {
	userId := c.Params[0].Value
	obId, _ := primitive.ObjectIDFromHex(userId)
	query := bson.M{"_id": obId}

	userCollection := configs.MI.DB.Collection("users")

	var users models.Users
	err := userCollection.FindOne(context.Background(), query).Decode(&users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   users,
	})
}

func DeleteUser(c *gin.Context) {
	userId := c.Params[0].Value
	obId, _ := primitive.ObjectIDFromHex(userId)
	query := bson.M{"_id": obId}
	userCollection := configs.MI.DB.Collection("users")

	result, err := userCollection.DeleteOne(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
		})
	}

	if result.DeletedCount > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "User deleted successfully",
			"data":    result,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "User not found",
	})
}

func AssignRole(c *gin.Context) {
	userCollection := configs.MI.DB.Collection("users")
	var input map[string]string
	c.ShouldBindJSON(&input)

	error := validation.Errors{
		"userId":    validation.Validate(input["userId"], validation.Required),
		"user_role": validation.Validate(input["user_role"], validation.Required),
	}.Filter()

	if error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": error.Error(),
		})
		return
	}

	userObjID, err := primitive.ObjectIDFromHex(input["userId"])
	filter := bson.M{"_id": userObjID}
	fields := bson.M{"$set": bson.M{"user_role": input["user_role"], "updated_at": time.Now()}}
	result, err := userCollection.UpdateOne(context.Background(), filter, fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  true,
			"error":   err.Error(),
			"message": "User not exists",
		})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  true,
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Role assigned successfully",
	})

}

func GetAllDataWithRating(c *gin.Context) {
	userId := c.Params[0].Value
	obId, _ := primitive.ObjectIDFromHex(userId)
	query := bson.M{"_id": obId}

	userCollection := configs.MI.DB.Collection("users")

	var users models.UsersResponse
	err := userCollection.FindOne(context.Background(), query).Decode(&users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  true,
			"message": err.Error(),
		})
		return
	}
	//fmt.Println(users)
	//fmt.Println(users.ID)

	var rating models.Rating
	//userObjId, _ := primitive.ObjectIDFromHex(users.ID)
	ratingCollection := configs.MI.DB.Collection("rating")
	err1 := ratingCollection.FindOne(context.Background(), bson.M{"user_id": users.ID}).Decode(&rating)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  true,
			"message": err1.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data": gin.H{
			"rating": rating,
			"user":   users,
		},
	})
}
