// router/user/create_post.go

package post

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Taker-Academy/kedubak-novaepitech/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func parseRequestBody(c *fiber.Ctx) (CreatePostRequest, error) {
	var body CreatePostRequest
	err := c.BodyParser(&body)
	return body, err
}

func createPostModel(body CreatePostRequest, claims jwt.MapClaims) models.Post {
	return models.Post {
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UserID:    claims["id"].(string),
		Title:     body.Title,
		Content:   body.Content,
		Comments:  []models.Comment{},
		UpVotes:   []string{},
	}
}

func insertPostIntoDB(client *mongo.Client, ctx context.Context, post models.Post) error {
	collection := client.Database("kedubak").Collection("posts")
	_, err := collection.InsertOne(ctx, post)
	return err
}

func CreatePost(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	token, err := ParseJWTToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token"})
	}

	body, err := parseRequestBody(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"ok": false, "message": "Invalid or missing parameters"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil {
		fmt.Printf("claims: %v\n", claims)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token"})
	}

	post := createPostModel(body, claims)

	err = insertPostIntoDB(client, ctx, post)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"ok": true, "data": post})
}
