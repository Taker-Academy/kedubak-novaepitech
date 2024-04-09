// router/add_groups.go

package router

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/Taker-Academy/kedubak-novaepitech/router/auth"
)

func AddAuthGroup(app *fiber.App, client *mongo.Client, ctx context.Context) {
	authGroup := app.Group("/auth")

	authGroup.Post("/register", func(c *fiber.Ctx) error {
		return auth.RegisterHandler(c, client, ctx)
	})
	authGroup.Post("/login", func(c *fiber.Ctx) error {
		return auth.LoginHandler(c, client, ctx)
	})
}
