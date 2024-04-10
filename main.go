// main.go

package main

import (
	"fmt"
	"os"

	"github.com/Taker-Academy/kedubak-novaepitech/common"
	"github.com/Taker-Academy/kedubak-novaepitech/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	app := fiber.New()

	app.Use(cors.New())

	client, ctx := common.ConnectToMongoDB()
	defer common.DisconnectFromMongoDB(client, ctx)

	router.AddGroups(app, client, ctx)

	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
}
