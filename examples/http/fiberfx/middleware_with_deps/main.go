package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
)

// LogService is a service for logging
type LogService struct {
	Prefix string
}

// NewLogService creates a new LogService
func NewLogService() *LogService {
	return &LogService{
		Prefix: "[LOG]",
	}
}

// AuthService is a service for authentication
type AuthService struct {
	// In a real application, this might contain database connections,
	// JWT verification keys, etc.
	AllowedTokens []string
}

// NewAuthService creates a new AuthService
func NewAuthService() *AuthService {
	return &AuthService{
		AllowedTokens: []string{"valid-token"},
	}
}

// IsValidToken checks if a token is valid
func (s *AuthService) IsValidToken(token string) bool {
	for _, t := range s.AllowedTokens {
		if t == token {
			return true
		}
	}
	return false
}

// LogMiddleware creates a middleware that depends on LogService
func LogMiddleware(logService *LogService) fiberfx.Middleware {
	return func(c *fiber.Ctx) error {
		fmt.Printf("%s %s - %s\n", logService.Prefix, c.Method(), c.Path())
		return c.Next()
	}
}

// AuthMiddleware creates a middleware that depends on AuthService
func AuthMiddleware(authService *AuthService) fiberfx.Middleware {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" || !authService.IsValidToken(token) {
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
		// Provide the services
		fx.Provide(
			NewLogService,
			NewAuthService,
		),

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

		// Register global middleware with dependency
		fiberfx.RegisterMiddleware("example", LogMiddleware),

		// Register middleware with a specific prefix and dependency
		fiberfx.RegisterMiddlewareWithPrefix("example", "/private", AuthMiddleware),

		// Run the app
		fiberfx.RunApp(":3003", "example", 5*time.Second),
	)

	// Start the application
	app.Run()
}
