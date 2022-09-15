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

func CreateRole(c *gin.Context) {
	roleCollection := configs.MI.DB.Collection("roles")
	var role models.Roles
	c.ShouldBindJSON(&role)

	//validation
	error := validation.Validate(role.Name, validation.Required)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": error.Error(),
		})
		return
	}

	role.CreatedAt = time.Now().UTC()
	role.UpdatedAt = time.Now().UTC()

	result, err := roleCollection.InsertOne(context.Background(), role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "role creating failed.",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    result,
		"status":  "success",
		"message": "role created successfully",
	})
}

func UpdateRole(c *gin.Context) {
	roleCollection := configs.MI.DB.Collection("roles")

	var role models.Roles
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	error := validation.Errors{
		"roleId": validation.Validate(c.Params[0].Value, validation.Required),
		"name":   validation.Validate(role.Name, validation.Required),
	}.Filter()

	if error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": error.Error(),
		})
		return
	}

	roleObjID, err := primitive.ObjectIDFromHex(c.Params[0].Value)
	filter := bson.M{"_id": roleObjID}
	fields := bson.M{"$set": bson.M{"name": role.Name, "updated_at": time.Now()}}
	result, err := roleCollection.UpdateOne(context.Background(), filter, fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": true,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Role updated successfully",
		"role":    result,
	})
}

func GetRoles(c *gin.Context) {
	roleCollection := configs.MI.DB.Collection("roles")
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))

	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))

	query := bson.M{}
	result, err := roleCollection.Find(context.Background(), query, &opt)

	if err != nil {
		c.JSON(500, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}
	defer result.Close(context.Background())
	var roles []models.Roles

	err = result.All(context.TODO(), &roles)
	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"status": true,
		"data":   roles,
	})
}

func GetRoleById(c *gin.Context) {
	roleId := c.Params[0].Value
	obId, _ := primitive.ObjectIDFromHex(roleId)
	query := bson.M{"_id": obId}

	userCollection := configs.MI.DB.Collection("roles")

	var result models.Roles
	userCollection.FindOne(context.Background(), query).Decode(&result)
	c.JSON(200, gin.H{
		"status": true,
		"data":   result,
	})
}

func DeleteRole(c *gin.Context) {
	roleId := c.Params[0].Value
	obId, _ := primitive.ObjectIDFromHex(roleId)
	query := bson.M{"_id": obId}
	userCollection := configs.MI.DB.Collection("roles")

	result, err := userCollection.DeleteOne(context.Background(), query)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  true,
		"message": "Role deleted successfully",
		"data":    result,
	})
}
