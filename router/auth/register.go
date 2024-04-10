// router/auth/register.go

package auth

import (
	"context"
	"fmt"
	"time"
	"os"

	"github.com/Taker-Academy/kedubak-novaepitech/common"
	"github.com/Taker-Academy/kedubak-novaepitech/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
}

func parseUser(c *fiber.Ctx) (*models.User, error) {
	user := new(models.User)
	err := c.BodyParser(user)
	return user, err
}

func checkExistingUser(client *mongo.Client, ctx context.Context, email string) (*models.User, error) {
	return common.GetUserByEmail(client, ctx, email)
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func generateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_TOKEN")))
}

func RegisterHandler(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	user, err := parseUser(c)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	existingUser, err := checkExistingUser(client, ctx, user.Email)
	if err == nil && existingUser != nil {
		return c.Status(400).SendString("A user with this email already exists")
	}
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(500).SendString("Failed to hash password")
	}
	user.Password = hashedPassword
	insertUser(client, ctx, user)
	tokenString, err := generateToken(user)
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
	fmt.Printf("User %s registered\n", user.Email);
	return c.JSON(response)
}
