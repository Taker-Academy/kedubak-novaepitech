// router/add_groups.go

package router

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/Taker-Academy/kedubak-novaepitech/router/auth"
	"github.com/Taker-Academy/kedubak-novaepitech/router/user"
	"github.com/Taker-Academy/kedubak-novaepitech/router/post"
)

func AddGroups(app *fiber.App, client *mongo.Client, ctx context.Context) {
	AddAuthGroup(app, client, ctx)
	AddUserGroup(app, client, ctx)
	AddPostGroup(app, client, ctx)
}

func AddAuthGroup(app *fiber.App, client *mongo.Client, ctx context.Context) {
	authGroup := app.Group("/auth")

	authGroup.Post("/register", func(c *fiber.Ctx) error {
		return auth.RegisterHandler(c, client, ctx)
	})
	authGroup.Post("/login", func(c *fiber.Ctx) error {
		return auth.LoginHandler(c, client, ctx)
	})
}

func AddUserGroup(app *fiber.App, client *mongo.Client, ctx context.Context) {
	userGroup := app.Group("/user")

	userGroup.Get("/me", func(c *fiber.Ctx) error {
		return user.GetUserInfos(c, client, ctx)
	})
	userGroup.Put("/edit", func(c *fiber.Ctx) error {
		return user.EditUserInfos(c, client, ctx)
	})
	userGroup.Delete("/remove", func(c *fiber.Ctx) error {
		return user.RemoveUser(c, client, ctx)
	})
}

func AddPostGroup(app *fiber.App, client *mongo.Client, ctx context.Context) {
	postGroup := app.Group("/post")

	postGroup.Get("/", func(c *fiber.Ctx) error {
		return post.GetPosts(c, client, ctx)
	})
	postGroup.Post("/", func(c *fiber.Ctx) error {
		return post.CreatePost(c, client, ctx)
	})
	postGroup.Get("/me", func(c *fiber.Ctx) error {
		return post.GetMyPosts(c, client, ctx)
	})
	postGroup.Get("/:id", func(c *fiber.Ctx) error {
        return post.GetPostByID(c, client, ctx)
    })
}
