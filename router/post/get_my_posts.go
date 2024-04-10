// router/user/get_my_posts.go

package post

import (
	"context"
	"net/http"

	"github.com/Taker-Academy/kedubak-novaepitech/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func fetchUserPostsFromDB(client *mongo.Client, ctx context.Context, userID string) (*mongo.Cursor, error) {
	collection := client.Database("kedubak").Collection("posts")
	cursor, err := collection.Find(ctx, bson.M{"userId": userID})
	return cursor, err
}

func GetMyPosts(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	token, err := ParseJWTToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token"})
	}

	cursor, err := fetchUserPostsFromDB(client, ctx, claims["id"].(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	posts, err := DecodePosts(cursor, ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	if posts == nil {
		posts = []models.Post{}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"ok": true, "data": posts})
}
