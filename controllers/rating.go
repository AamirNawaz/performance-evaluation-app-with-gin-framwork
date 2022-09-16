package controllers

import "github.com/gin-gonic/gin"

func GetRating(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Gat all rating",
	})
}
func ThumbsUp(c *gin.Context) {
	c.JSON(200, "ThumbsUp function calling")
}

func ThumbsDown(c *gin.Context) {
	c.JSON(200, "ThumbsDown function calling")
}
