package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"performance-evaluation-app-with-gin/routes"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	cors "github.com/rs/cors/wrapper/gin"

	db "performance-evaluation-app-with-gin/configs"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func main() {
	r := gin.Default()
	//mongo connection
	db.Connect()

	// Getting host ip
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	// fmt.Println("address:", addrs)
	// fmt.Println("address:", addrs[2])
	// fmt.Println("host", host)

	r.Use(cors.Default())
	routes.SetupRoutes(r)

	//Method 1 : Dynamically extracting system ip
	prefix := "192."
	found := false
	var hostIpIndex = 0

	// Loop through the array and check for elements with the specified prefix
	for index, ipValue := range addrs {
		ipString := ipValue.String() // Convert net.IP to string
		if strings.HasPrefix(ipString, prefix) {
			hostIpIndex = index
			found = true
			break
		}
	}
	if found {
		err := r.Run((addrs[hostIpIndex].String()) + ":" + (os.Getenv("PORT")))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("IP address of the system not found")
	}

	// Method 2: Static way to start server with http://localhost
	// err := r.Run(("http://localhost"+(os.Getenv("PORT")))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
}
