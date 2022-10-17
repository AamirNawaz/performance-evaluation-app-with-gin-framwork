package routes

import (
	"github.com/gin-gonic/gin"
	"performance-evaluation-app-with-gin/controllers"
	middleware "performance-evaluation-app-with-gin/middlewares"
)

func SetupRoutes(router *gin.Engine) {

	//Auth Routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/", controllers.Login)
		auth.POST("/get-token", controllers.GetNewAccessToken)
		auth.GET("/logout", controllers.Logout)
	}

	//User Routes
	user := router.Group("/api/users")
	{
		user.GET("/", controllers.GetUsers)
		user.GET("/all/:id", controllers.GetAllDataWithRating)
		user.GET("/:id", middleware.CheckAuth, controllers.GetUserById)
		user.DELETE("/delete/:id", middleware.CheckAuth, middleware.CheckRole, controllers.DeleteUser)
		user.POST("/assign-role", middleware.CheckAuth, middleware.CheckRole, controllers.AssignRole)

	}

	//Roles Routes
	role := router.Group("/api/role")
	{
		role.GET("/", middleware.CheckAuth, middleware.CheckRole, controllers.GetRoles)
		role.GET("/:id", middleware.CheckAuth, middleware.CheckRole, controllers.GetRoleById)
		role.POST("/create", middleware.CheckAuth, middleware.CheckRole, controllers.CreateRole)
		role.PUT("/update/:id", middleware.CheckAuth, middleware.CheckRole, controllers.UpdateRole)
		role.DELETE("/delete/:id", middleware.CheckAuth, middleware.CheckRole, controllers.DeleteRole)
	}

	//Rating routes
	rating := router.Group("/api/rating")
	{
		rating.GET("/", controllers.GetRating)
		rating.POST("/thumbs-up", controllers.ThumbUp)
		rating.POST("/thumbs-down", controllers.ThumbDown)
		rating.GET("/positive-rating", controllers.GetPositiveRating)
		rating.GET("/negative-rating", controllers.GetNegativeRating)

	}

}
