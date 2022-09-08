package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"performance_valuation_app/controllers"
)

func SetupRoutes(app *fiber.App) {

	//middleware
	api := app.Group("/api", logger.New())

	//Auth group middleware
	auth := api.Group("/auth")
	auth.Post("/signup", controllers.Signup)
	auth.Post("/", controllers.Login)
	auth.Get("/logout", controllers.Logout)

	//User group middleware
	//user := api.Group("/user")
	//user.Get("/", middleware.CheckAuth, controllers.GetUsers)
	//user.Get("/:id", middlewares.AdminProtected, controllers.GetUser)

}
