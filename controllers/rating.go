package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"performance-evaluation-app-with-gin/configs"
	"performance-evaluation-app-with-gin/models"
	"strconv"
	"time"
)

type RatingResponse struct {
	ID      *primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Count   int64                 `json:"count"`
	RatedBy []*primitive.ObjectID `json:"rated_by" bson:"rated_by"`
}

func GetRating(c *gin.Context) {
	ratingCollection := configs.MI.DB.Collection("rating")
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))

	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))

	query := bson.M{}
	result, err := ratingCollection.Find(context.Background(), query, &opt)

	if err != nil {
		c.JSON(500, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}
	defer result.Close(context.Background())
	var rating []models.Rating

	err = result.All(context.TODO(), &rating)
	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"status": true,
		"data":   rating,
	})
}

func ThumbUp(c *gin.Context) {
	ratingCollection := configs.MI.DB.Collection("rating")
	var rating models.Rating
	c.ShouldBindJSON(&rating)

	fmt.Println(rating)
	//Multiple fields validation
	err := validation.Errors{
		"user_id":  validation.Validate(rating.UserId, validation.Required),
		"rated_by": validation.Validate(rating.RatedBy, validation.Required),
	}.Filter()

	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	rating.ThumbsUp = 1
	rating.CreatedAt = time.Now().UTC()
	rating.UpdatedAt = time.Now().UTC()

	result, Err := ratingCollection.InsertOne(context.Background(), rating)
	if Err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to rate user.",
			"error":   Err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    result,
		"status":  "success",
		"message": "user rated successfully",
	})

}

func ThumbDown(c *gin.Context) {
	ratingCollection := configs.MI.DB.Collection("rating")
	var rating models.Rating
	c.ShouldBindJSON(&rating)

	//Multiple fields validation
	err := validation.Errors{
		"user_id":  validation.Validate(rating.UserId, validation.Required),
		"rated_by": validation.Validate(rating.RatedBy, validation.Required),
	}.Filter()

	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	rating.ThumbsDown = 1
	rating.CreatedAt = time.Now().UTC()
	rating.UpdatedAt = time.Now().UTC()

	result, Err := ratingCollection.InsertOne(context.Background(), rating)
	if Err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to rate user.",
			"error":   Err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    result,
		"status":  "success",
		"message": "user rated successfully",
	})

}

func GetPositiveRating(c *gin.Context) {
	ratingCollection := configs.MI.DB.Collection("rating")
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))

	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))

	//query := bson.M{}
	id, _ := primitive.ObjectIDFromHex(c.Query("id"))
	fmt.Println(id)

	matchStage := bson.D{{"$match", bson.D{{"user_id", id}, {"thumbs_up", bson.D{{"$gt", 0}}}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "$user_id"}, {"count", bson.D{{"$sum", "$thumbs_up"}}}, {"rated_by", bson.D{{"$push", "$rated_by"}}}}}}
	result, err := ratingCollection.Aggregate(context.Background(), mongo.Pipeline{matchStage, groupStage})

	if err != nil {
		c.JSON(500, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	defer result.Close(context.Background())
	var rating []RatingResponse

	err = result.All(context.TODO(), &rating)
	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"status": true,
		"data":   rating,
	})
}

func GetNegativeRating(c *gin.Context) {
	ratingCollection := configs.MI.DB.Collection("rating")
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))

	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))

	//query := bson.M{}
	id, _ := primitive.ObjectIDFromHex(c.Query("id"))
	fmt.Println(id)

	matchStage := bson.D{{"$match", bson.D{{"user_id", id}, {"thumbs_down", bson.D{{"$gt", 0}}}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "$user_id"}, {"count", bson.D{{"$sum", "$thumbs_down"}}}, {"rated_by", bson.D{{"$push", "$rated_by"}}}}}}
	result, err := ratingCollection.Aggregate(context.Background(), mongo.Pipeline{matchStage, groupStage})

	if err != nil {
		c.JSON(500, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	defer result.Close(context.Background())
	var rating []RatingResponse

	err = result.All(context.TODO(), &rating)
	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"status": true,
		"data":   rating,
	})
}
