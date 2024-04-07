package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Taker-Academy/kedubak-novaepitech/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
    }
	app := fiber.New()

	app.Use(cors.New(cors.Config{
        AllowOrigins: "*", // ou "*" pour autoriser toutes les origines
        AllowHeaders: "Origin, Content-Type, Accept",
    }))

	app.Post("/auth/register", func(c *fiber.Ctx) error {
		client, ctx := connectToMongoDB()
		defer disconnectFromMongoDB(client, ctx)

		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		insertUser(client, ctx, user)

		// Create a new token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  user.ID,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

		// Sign and get the complete encoded token as a string
		tokenString, err := token.SignedString([]byte("your-secret-key"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Create a new response
		response := models.Response{
			Ok: true,
			Data: models.Data{
				Token: tokenString,
				User:  *user,
			},
		}

		// Return the response as JSON
		return c.JSON(response)
	})

	app.Listen(":8080")
}

func connectToMongoDB() (*mongo.Client, context.Context) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		fmt.Println("'MONGODB_URI' environmental variable not found")
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	return client, context.TODO()
}

func disconnectFromMongoDB(client *mongo.Client, ctx context.Context) {
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func insertUser(client *mongo.Client, ctx context.Context, user *models.User) {
	collection := client.Database("kedubak").Collection("users")
	user.CreatedAt = time.Now()
	user.LastUpVote = time.Now().Add(-1 * time.Minute)
	insertResult, err := collection.InsertOne(ctx, user)
	if err != nil {
		panic(err)
	}
	objID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		panic("Cannot convert InsertedID to ObjectID")
	}
	user.ID = objID.Hex()
	fmt.Printf("User added with ID %v\n", insertResult.InsertedID)
}
