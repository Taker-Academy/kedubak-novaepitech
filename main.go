// main.go

package main

import (
	"fmt"

	"github.com/Taker-Academy/kedubak-novaepitech/router"
	"github.com/Taker-Academy/kedubak-novaepitech/common"
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

	app.Listen(":8080")
}
