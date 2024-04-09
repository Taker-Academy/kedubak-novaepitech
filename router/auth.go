// router/auth.go

package router

import (
	"fmt"
	"time"
	"context"

	"github.com/Taker-Academy/kedubak-novaepitech/models"
	"github.com/Taker-Academy/kedubak-novaepitech/common"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func AddAuthGroup(app *fiber.App, client *mongo.Client, ctx context.Context) {
	authGroup := app.Group("/auth")

	authGroup.Post("/register", func(c *fiber.Ctx) error {
		return registerHandler(c, client, ctx)
	})
	authGroup.Post("/login", func(c *fiber.Ctx) error {
		return loginHandler(c, client, ctx)
	})
}

func loginHandler(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	fmt.Println("Logging a user")
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	dbUser, err := common.GetUserByEmail(client, ctx, user.Email)
	if err != nil {
		return c.Status(401).SendString("Invalid credentials")
	}
	if !common.ComparePasswords(dbUser.Password, user.Password) {
		return c.Status(401).SendString("Invalid credentials")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  dbUser.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	response := models.Response{
		Ok: true,
		Data: models.Data{
			Token: tokenString,
			User:  *dbUser,
		},
	}
	return c.JSON(response)
}

func registerHandler(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	fmt.Println("Registering a new user")
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// Check if a user with the same email already exists
	existingUser, err := common.GetUserByEmail(client, ctx, user.Email)
	if err == nil && existingUser != nil {
		return c.Status(400).SendString("A user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).SendString("Failed to hash password")
	}
	user.Password = string(hashedPassword)
	common.InsertUser(client, ctx, user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	response := models.Response{
		Ok: true,
		Data: models.Data{
			Token: tokenString,
			User:  *user,
		},
	}
	return c.JSON(response)
}
