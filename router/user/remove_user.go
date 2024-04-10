// router/user/remove_user.go

package user

import (
	"context"
	"net/http"
	"os"
	"strings"
	"fmt"

	"github.com/Taker-Academy/kedubak-novaepitech/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func parseTokenRemove(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := strings.Split(c.Get("Authorization"), " ")[1]
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
}

func extractUserID(token *jwt.Token) (primitive.ObjectID, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return primitive.NilObjectID, fmt.Errorf("invalid token")
	}
	return primitive.ObjectIDFromHex(claims["id"].(string))
}

func fetchUser(ctx context.Context, client *mongo.Client, userID primitive.ObjectID) (models.User, error) {
	collection := client.Database("kedubak").Collection("users")
	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	return user, err
}

func removeUserFromDB(ctx context.Context, client *mongo.Client, userID primitive.ObjectID) error {
	collection := client.Database("kedubak").Collection("users")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": userID})
	return err
}

func RemoveUser(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	token, err := parseTokenRemove(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token"})
	}

	userID, err := extractUserID(token)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	user, err := fetchUser(ctx, client, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"ok": false, "message": "User not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	err = removeUserFromDB(ctx, client, userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"ok": true, "data": user})
}
