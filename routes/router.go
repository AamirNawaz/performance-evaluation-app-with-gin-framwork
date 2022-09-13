package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"performance-evaluation-app/controllers"
	middleware "performance-evaluation-app/middlewares"
)

func SetupRoutes(app *fiber.App) {

	//middleware
	api := app.Group("/api", logger.New())

	//Auth group middleware
	auth := api.Group("/auth")
	auth.Post("/signup", controllers.Signup)
	auth.Post("/", controllers.Login)
	auth.Get("/logout", middleware.CheckAuth, controllers.Logout)
	auth.Post("/get-token", controllers.GetNewAccessToken)

	//User group middleware
	user := api.Group("/users")
	user.Get("/", controllers.GetUsers)
	user.Get("/:id", controllers.GetUserById)
	user.Delete("/delete/:id", controllers.DeleteUser)
	//user.Get("/", middleware.CheckAuth, controllers.GetUsers)
	//user.Get("/:id", middlewares.AdminProtected, controllers.GetUser)

	//*****************Roles Routes*****************/
	role := api.Group("/role")
	role.Get("/", controllers.GetRoles)
	role.Get("/:id", controllers.GetRoleById)
	role.Post("/create", controllers.CreateRole)
	role.Put("/update/:id", controllers.UpdateRole)
	role.Delete("/delete/:id", controllers.DeleteRole)

}
