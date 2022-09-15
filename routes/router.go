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
		auth.GET("/logout", controllers.Logout())
	}

	//User Routes
	user := router.Group("/api/users")
	{
		user.GET("/", middleware.CheckAuth, controllers.GetUsers)
		user.GET("/:id", controllers.GetUserById)
		user.DELETE("/delete/:id", controllers.DeleteUser)
		user.POST("/assign-role", controllers.AssignRole)

	}

	//Roles Routes
	role := router.Group("/api/role")
	{
		role.GET("/", controllers.GetRoles)
		role.GET("/:id", controllers.GetRoleById)
		role.POST("/create", controllers.CreateRole)
		role.PUT("/update/:id", controllers.UpdateRole)
		role.DELETE("/delete/:id", controllers.DeleteRole)
	}

}
