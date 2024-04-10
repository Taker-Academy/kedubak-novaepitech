// router/user/get_user_infos.go

package user

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/Taker-Academy/kedubak-novaepitech/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getAuthorizationToken(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}
	return strings.Replace(authHeader, "Bearer ", "", -1), nil
}

func parseToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
}

func getUserIDFromToken(token *jwt.Token) (primitive.ObjectID, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return primitive.NilObjectID, errors.New("invalid token")
	}
	return primitive.ObjectIDFromHex(claims["id"].(string))
}

func fetchUserFromDatabase(client *mongo.Client, ctx context.Context, userID primitive.ObjectID) (models.User, error) {
	collection := client.Database("kedubak").Collection("users")
	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	return user, err
}

func GetUserInfos(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	tokenStr, err := getAuthorizationToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": err.Error()})
	}

	token, err := parseToken(tokenStr)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid token"})
	}

	userID, err := getUserIDFromToken(token)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	user, err := fetchUserFromDatabase(client, ctx, userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"ok": true, "data": user})
}
