// router/auth/login.go

package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Taker-Academy/kedubak-novaepitech/common"
	"github.com/Taker-Academy/kedubak-novaepitech/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/golang-jwt/jwt"
)

func parseLoginUser(c *fiber.Ctx) (*models.User, error) {
	user := new(models.User)
	err := c.BodyParser(user)
	return user, err
}

func validateUser(client *mongo.Client, ctx context.Context, user *models.User) (*models.User, error) {
	dbUser, err := common.GetUserByEmail(client, ctx, user.Email)
	if err != nil || !common.ComparePasswords(dbUser.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}
	return dbUser, nil
}

func generateLoginToken(dbUser *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  dbUser.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte("your-secret-key"))
}

func LoginHandler(c *fiber.Ctx, client *mongo.Client, ctx context.Context) error {
	user, err := parseLoginUser(c)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	dbUser, err := validateUser(client, ctx, user)
	if err != nil {
		return c.Status(401).SendString(err.Error())
	}
	tokenString, err := generateLoginToken(dbUser)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	response := models.Response{
		Ok: true,
		Data: models.Data{
			Token: tokenString,
			User:  *dbUser,
		},
	}
	fmt.Printf("User %s logged in\n", dbUser.Email)
	return c.JSON(response)
}
