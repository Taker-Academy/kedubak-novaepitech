// router/user/get_user_infos.go

package user

import (
	"context"
	"net/http"
	"strings"
	"fmt"
	"os"

	"github.com/Taker-Academy/kedubak-novaepitech/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserInfos(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	fmt.Println("Get user infos")
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Missing authorization header"})
	}
	tokenStr := strings.Replace(authHeader, "Bearer ", "", -1)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		fmt.Printf("Invalid token 1: %v\n", err)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid token"})
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		fmt.Println("Invalid token 2")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid token"})
	}
	userID, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err != nil {
		fmt.Printf("Error converting userID to ObjectID: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}
	collection := client.Database("kedubak").Collection("users")
	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}
	fmt.Println("User found")
	return c.Status(http.StatusOK).JSON(fiber.Map{"ok": true, "data": user})
}
