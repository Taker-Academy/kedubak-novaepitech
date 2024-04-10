// router/user/delete_posts.go

package post

import (
	"context"
	"net/http"
	"errors"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/Taker-Academy/kedubak-novaepitech/models"
)

func GetUserIDFromTokenClaims(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid JWT token")
	}
	
	id, ok := claims["id"].(string)
	if !ok {
		return "", errors.New("invalid JWT token")
	}
	
	return id, nil
}

func getPostByID(ctx context.Context, collection *mongo.Collection, objectID primitive.ObjectID) (*models.Post, error) {
	var post models.Post
	err := collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		return nil, err
	}
	
	return &post, nil
}

func DeletePostFromDB(ctx context.Context, collection *mongo.Collection, objectID primitive.ObjectID) error {
	_, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}


func DeletePost(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	token, err := ParseJWTToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token"})
	}

	id, err := GetUserIDFromTokenClaims(token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token"})
	}

	postID := c.Params("id")
	if postID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"ok": false, "message": "Missing post ID"})
	}

	collection := client.Database("kedubak").Collection("posts")
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"ok": false, "message": "Invalid post ID"})
	}

	post, err := getPostByID(ctx, collection, objectID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"ok": false, "message": "Post not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	if post.UserID != id {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"ok": false, "message": "User is not the owner of the post"})
	}

	err = DeletePostFromDB(ctx, collection, objectID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"ok": true, "message": "Post deleted successfully"})
}
