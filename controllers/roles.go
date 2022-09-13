package controllers

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"performance-evaluation-app/configs"
	"performance-evaluation-app/models"
	"strconv"
	"time"
)

func CreateRole(c *fiber.Ctx) error {
	roleCollection := configs.MI.DB.Collection("roles")
	var role models.Roles
	c.BodyParser(&role)

	//validation
	error := validation.Validate(role.Name, validation.Required)
	if error != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": error.Error(),
		})
	}

	role.CreatedAt = time.Now().UTC()
	role.UpdatedAt = time.Now().UTC()

	result, err := roleCollection.InsertOne(context.Background(), role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "role creating failed.",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"status":  "success",
		"message": "role created successfully",
	})

}

func UpdateRole(c *fiber.Ctx) error {
	roleCollection := configs.MI.DB.Collection("roles")

	var role models.Roles
	c.BodyParser(&role)

	error := validation.Errors{
		"roleId": validation.Validate(c.Params("id"), validation.Required),
		"name":   validation.Validate(role.Name, validation.Required),
	}.Filter()

	if error != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": error.Error(),
		})
	}

	roleObjID, err := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{"_id": roleObjID}
	fields := bson.M{"$set": bson.M{"name": role.Name, "updated_at": time.Now()}}
	result, err := roleCollection.UpdateOne(context.Background(), filter, fields)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": true,
			"error":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Role updated successfully",
		"role":    result,
	})
}

func GetRoles(c *fiber.Ctx) error {
	roleCollection := configs.MI.DB.Collection("roles")

	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))

	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))

	query := bson.M{}

	result, err := roleCollection.Find(context.Background(), query, &opt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"error":  err.Error(),
		})
	}
	defer result.Close(context.Background())
	var roles []models.Roles

	err = result.All(context.TODO(), &roles)
	if err != nil {
		panic(err)
	}

	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data":   roles,
	})
}

func GetRoleById(c *fiber.Ctx) error {
	roleId := c.Params("id")
	obId, _ := primitive.ObjectIDFromHex(roleId)
	query := bson.M{"_id": obId}

	userCollection := configs.MI.DB.Collection("roles")

	var result models.Roles
	userCollection.FindOne(context.Background(), query).Decode(&result)
	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data":   result,
	})
}

func DeleteRole(c *fiber.Ctx) error {
	roleId := c.Params("id")
	obId, _ := primitive.ObjectIDFromHex(roleId)
	query := bson.M{"_id": obId}
	userCollection := configs.MI.DB.Collection("roles")

	result, err := userCollection.DeleteOne(context.Background(), query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  true,
		"message": "Role deleted successfully",
		"data":    result,
	})
}
