// router/user/edit_user_infos.go

package user

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/Taker-Academy/kedubak-novaepitech/models"
	"github.com/Taker-Academy/kedubak-novaepitech/router/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func parseAndValidateToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := strings.Split(c.Get("Authorization"), " ")[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	return token, err
}

func extractUserIDFromToken(token *jwt.Token) (primitive.ObjectID, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return primitive.NilObjectID, errors.New("invalid JWT token")
	}
	return primitive.ObjectIDFromHex(claims["id"].(string))
}

func parseRequestBody(c *fiber.Ctx) (map[string]interface{}, error) {
	var body map[string]interface{}
	err := c.BodyParser(&body)
	return body, err
}

func hashPasswordIfProvided(body map[string]interface{}) error {
	if password, ok := body["password"].(string); ok {
		hashedPassword, err := auth.HashPassword(password)
		if err != nil {
			return err
		}
		body["password"] = hashedPassword
	}
	return nil
}

func updateUserInfos(client *mongo.Client, ctx context.Context, userID primitive.ObjectID, body map[string]interface{}) error {
	collection := client.Database("kedubak").Collection("users")
	update := bson.M{"$set": body}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	return err
}

func fetchUpdatedUser(client *mongo.Client, ctx context.Context, userID primitive.ObjectID) (models.User, error) {
	var user models.User
	collection := client.Database("kedubak").Collection("users")
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	return user, err
}

func EditUserInfos(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	token, err := parseAndValidateToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token"})
	}

	userID, err := extractUserIDFromToken(token)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	body, err := parseRequestBody(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"ok": false, "message": "Failed to parse request body"})
	}

	err = hashPasswordIfProvided(body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Failed to hash password"})
	}

	err = updateUserInfos(client, ctx, userID, body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	user, err := fetchUpdatedUser(client, ctx, userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"ok": true, "data": user})
}
