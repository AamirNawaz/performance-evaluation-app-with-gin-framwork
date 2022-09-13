package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"performance-evaluation-app/configs"
	"performance-evaluation-app/models"
	"strconv"
)

func GetUsers(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))

	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))

	//with projection we can skip specific field to retrive
	opt.SetProjection(bson.D{{"password", 0}})

	query := bson.M{}

	userCollection := configs.MI.DB.Collection("users")
	result, err := userCollection.Find(context.Background(), query, &opt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"error":  err.Error(),
		})
	}

	defer result.Close(context.Background())
	var users []models.Users

	err = result.All(context.TODO(), &users)
	if err != nil {
		panic(err)
	}

	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data":   users,
	})

}

func GetUserById(c *fiber.Ctx) error {
	userId := c.Params("id")
	obId, _ := primitive.ObjectIDFromHex(userId)
	query := bson.M{"_id": obId}

	userCollection := configs.MI.DB.Collection("users")

	var users models.Users
	err := userCollection.FindOne(context.Background(), query).Decode(&users)
	if err != nil {
		return c.Status(200).JSON(fiber.Map{
			"status":  true,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data":   users,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	obId, _ := primitive.ObjectIDFromHex(userId)
	query := bson.M{"_id": obId}
	userCollection := configs.MI.DB.Collection("users")

	result, err := userCollection.DeleteOne(context.Background(), query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	if result.DeletedCount > 0 {
		return c.Status(200).JSON(fiber.Map{
			"status":  true,
			"message": "User deleted successfully",
			"data":    result,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  true,
		"message": "User not found",
	})
}
