// router/user/get_post_by_id.go

package post

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/Taker-Academy/kedubak-novaepitech/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func fetchPostFromDB(client *mongo.Client, ctx context.Context, postID string) (*models.Post, error) {
    collection := client.Database("kedubak").Collection("posts")
    var post models.Post
    objectID, err := primitive.ObjectIDFromHex(postID)
    if err != nil {
        return nil, err
    }
    err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
    return &post, err
}

func GetPostByID(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	_, err := ParseJWTToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token"})
	}

	postID := c.Params("id")
	if postID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"ok": false, "message": "Missing post ID"})
	}

	post, err := fetchPostFromDB(client, ctx, postID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"ok": false, "message": "Post not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"ok": true, "data": post})
}
