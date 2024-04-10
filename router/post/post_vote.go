// router/user/vote_post.go

package post

import (
	"context"
	"net/http"
	"time"

	"github.com/Taker-Academy/kedubak-novaepitech/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gofiber/fiber/v2"
)

func GetPostByIDForVote(c *fiber.Ctx, client *mongo.Client, ctx context.Context) (*models.Post, error) {
    postID := c.Params("id")
    if postID == "" {
        return nil, c.Status(http.StatusBadRequest).JSON(fiber.Map{"ok": false, "message": "Missing post ID"})
    }

    collection := client.Database("kedubak").Collection("posts")
    var post models.Post
    objectID, err := primitive.ObjectIDFromHex(postID)
    if err != nil {
        return nil, err
    }
    err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, c.Status(http.StatusNotFound).JSON(fiber.Map{"ok": false, "message": "Post not found"})
        }
        return nil, c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
    }

    return &post, nil
}

func VotePost(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	token, err := ParseJWTToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token"})
	}

	id, err := GetUserIDFromTokenClaims(token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"ok": false, "message": "Invalid JWT token"})
	}

	// Get the user from the database
	userCollection := client.Database("kedubak").Collection("users")
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	// Check if the user can vote
	if user.LastUpVote.Add(time.Minute).After(time.Now()) {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"ok": false, "message": "You can only vote every minute"})
	}

	post, err := GetPostByIDForVote(c, client, ctx)
	if err != nil {
		return err
	}

	for _, voter := range post.UpVotes {
		if voter == id {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"ok": false, "message": "You have already voted for this post"})
		}
	}

	post.UpVotes = append(post.UpVotes, id)
	user.LastUpVote = time.Now()

	// Update the post and the user in the database
	postCollection := client.Database("kedubak").Collection("posts")
	_, err = postCollection.UpdateOne(ctx, bson.M{"_id": post.ID}, bson.M{
		"$set": bson.M{
			"upVotes": post.UpVotes,
		},
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"lastUpVote": user.LastUpVote,
		},
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"ok": false, "message": "Internal server error"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"ok": true, "message": "Post upvoted"})
}
