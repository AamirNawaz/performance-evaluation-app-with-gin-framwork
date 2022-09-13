package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	db "performance-evaluation-app/configs"
	"performance-evaluation-app/routes"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Static("/", "./public")

	//mongo connection
	db.Connect()

	//Logging middleware
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	//Getting host ip
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)

	routes.SetupRoutes(app)
	err := app.Listen((addrs[1].String()) + ":" + (os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}
}
