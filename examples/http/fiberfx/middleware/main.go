package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
)

// LogMiddleware logs all requests
func LogMiddleware() fiberfx.Middleware {
	return func(c *fiber.Ctx) error {
		fmt.Printf("[LOG] %s - %s\n", c.Method(), c.Path())
		return c.Next()
	}
}

// AuthMiddleware checks for authorization header
func AuthMiddleware() fiberfx.Middleware {
	return func(c *fiber.Ctx) error {
		if c.Get("Authorization") == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}
		return c.Next()
	}
}

// PublicHandler handles public requests
func PublicHandler(c *fiber.Ctx) error {
	return c.SendString("This is public content")
}

// PrivateHandler handles private requests
func PrivateHandler(c *fiber.Ctx) error {
	return c.SendString("This is private content")
}

func main() {
	app := fx.New(
		// Create the Fiber app with middleware support
		fiberfx.App(
			"example",
			fiberfx.Routes(
				[]fiberfx.RouteFx{
					// Public route
					fiberfx.Get("/public", PublicHandler),

					// Private route
					fiberfx.Get("/private/data", PrivateHandler),
				},
			),
			// Enable middleware injection
			fiberfx.WithMiddlewares(),
		),

		// Register global middleware
		fiberfx.RegisterMiddleware("example", LogMiddleware),

		// Register middleware with a specific prefix
		fiberfx.RegisterMiddlewareWithPrefix("example", "/private", AuthMiddleware),

		// Run the app
		fiberfx.RunApp(":3002", "example", 5*time.Second),
	)

	// Start the application
	app.Run()
}
