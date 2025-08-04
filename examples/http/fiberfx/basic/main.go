package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
)

// HelloHandler is a simple handler that returns a greeting
func HelloHandler(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

// UserHandler returns information about a user
func UserHandler(c *fiber.Ctx) error {
	userId := c.Params("id")
	return c.SendString(fmt.Sprintf("User ID: %s", userId))
}

func main() {
	app := fx.New(
		// Create the Fiber app
		fiberfx.App(
			"example",
			fiberfx.Routes(
				[]fiberfx.RouteFx{
					// Simple GET route
					fiberfx.Get("/hello", HelloHandler),

					// Route with parameter
					fiberfx.Get("/users/:id", UserHandler),
				},
			),
		),

		// Run the app
		fiberfx.RunApp(":3001", "example", 5*time.Second),
	)

	// Start the application
	app.Run()
}
