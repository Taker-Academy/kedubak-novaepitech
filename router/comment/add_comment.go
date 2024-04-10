// router/user/add_comment.go

package comment

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Taker-Academy/kedubak-novaepitech/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func parseToken(c *fiber.Ctx) (*jwt.Token, error) {
	authHeader := c.Get("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	return token, err
}

func getClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid JWT token")
	}

	return claims, nil
}

func getPostID(c *fiber.Ctx) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(c.Params("id"))
}

func getRequestBody(c *fiber.Ctx) (string, error) {
	var body struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&body); err != nil {
		return "", err
	}

	return body.Content, nil
}

func createComment(id string, postID primitive.ObjectID, content string) models.Comment {
	return models.Comment{
		ID:        postID.Hex(),
		FirstName: id,
		Content:   content,
		CreatedAt: time.Now(),
	}
}

func addCommentToPost(client *mongo.Client, ctx context.Context, postID primitive.ObjectID, comment models.Comment) error {
    collection := client.Database("kedubak").Collection("posts")
    filter := bson.M{"_id": postID}
    update := bson.M{"$push": bson.M{"comments": comment}}

    _, err := collection.UpdateOne(ctx, filter, update)

    return err
}

func AddComment(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	token, err := parseToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token"})
	}

	claims, err := getClaims(token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token"})
	}

	id, ok := claims["id"].(string)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token: missing id claim"})
	}

	postID, err := getPostID(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"ok": false, "message": "Invalid post ID"})
	}

	content, err := getRequestBody(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"ok": false, "message": "Invalid request body"})
	}

	comment := createComment(id, postID, content)

    err = addCommentToPost(client, ctx, postID, comment)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
    }

    return c.Status(http.StatusCreated).JSON(fiber.Map{"ok": true, "data": comment})
}
