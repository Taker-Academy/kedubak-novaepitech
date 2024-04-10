// router/user/get_posts.go

package post

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/Taker-Academy/kedubak-novaepitech/models"
)

func ParseJWTToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := strings.Split(c.Get("Authorization"), " ")[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	return token, err
}

func fetchPostsFromDB(client *mongo.Client, ctx context.Context) (*mongo.Cursor, error) {
	collection := client.Database("kedubak").Collection("posts")
	cursor, err := collection.Find(ctx, bson.M{})
	return cursor, err
}

func decodePosts(cursor *mongo.Cursor, ctx context.Context) ([]models.Post, error) {
	var posts []models.Post
	err := cursor.All(ctx, &posts)
	return posts, err
}

func GetPosts(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	_, err := ParseJWTToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token"})
	}

	cursor, err := fetchPostsFromDB(client, ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	posts, err := decodePosts(cursor, ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	if posts == nil {
		posts = []models.Post{}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"ok": true, "data": posts})
}
