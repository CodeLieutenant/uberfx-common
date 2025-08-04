package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
)

// Simple middleware without dependencies

// LogMiddleware logs request information
func LogMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Printf("[LOG] %s - %s\n", c.Method(), c.Path())
		return c.Next()
	}
}

// RateLimitMiddleware implements a simple rate limiter
func RateLimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// In a real application, this would check rate limits
		fmt.Println("[RATE LIMIT] Checking rate limits")
		return c.Next()
	}
}

// Middleware with dependencies

// AuthService handles authentication
type AuthService struct {
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

// MetricsService handles metrics collection
type MetricsService struct {
	Prefix string
}

// NewMetricsService creates a new MetricsService
func NewMetricsService() *MetricsService {
	return &MetricsService{
		Prefix: "[METRICS]",
	}
}

// AuthMiddleware creates a middleware that depends on AuthService
func AuthMiddleware(authService *AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" || !authService.IsValidToken(token) {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}
		return c.Next()
	}
}

// MetricsMiddleware creates a middleware that depends on MetricsService
func MetricsMiddleware(metricsService *MetricsService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Printf("%s Recording metrics for %s\n", metricsService.Prefix, c.Path())
		return c.Next()
	}
}

// Route handlers

// PublicHandler handles public requests
func PublicHandler(c *fiber.Ctx) error {
	return c.SendString("This is public content")
}

// PrivateHandler handles private requests
func PrivateHandler(c *fiber.Ctx) error {
	return c.SendString("This is private content")
}

// AdminHandler handles admin requests
func AdminHandler(c *fiber.Ctx) error {
	return c.SendString("This is admin content")
}

// MetricsHandler handles metrics requests
func MetricsHandler(c *fiber.Ctx) error {
	return c.SendString("Metrics data")
}

func main() {
	app := fx.New(
		// Provide the services
		fx.Provide(
			NewAuthService,
			NewMetricsService,
		),

		// Create the Fiber app
		fiberfx.App(
			"example",
			fiberfx.Routes(
				[]fiberfx.RouteFx{
					// Route with no middleware
					fiberfx.Get("/public", PublicHandler),

					// Route with middleware (no dependencies)
					fiberfx.GetWithMiddleware("/api/data",
						[]fiber.Handler{LogMiddleware(), RateLimitMiddleware()},
						PublicHandler),

					// Route with middleware with dependencies
					fiberfx.GetWithMiddlewareFx("/private",
						[]fiberfx.RouteMiddlewareFunc{AuthMiddleware},
						PrivateHandler),

					// Route with multiple middleware with dependencies
					fiberfx.GetWithMiddlewareFx("/admin",
						[]fiberfx.RouteMiddlewareFunc{AuthMiddleware, MetricsMiddleware},
						AdminHandler),

					// Route with router callback and middleware with dependencies
					fiberfx.GetWithRouterCallbackAndMiddlewareFx("/metrics/:type",
						func(router fiber.Router) {
							// Configure the route
							router.Use(func(c *fiber.Ctx) error {
								fmt.Println("[ROUTER] Processing metrics request")
								return c.Next()
							})
						},
						[]fiberfx.RouteMiddlewareFunc{MetricsMiddleware},
						MetricsHandler),
				},
			),
		),

		// Run the app
		fiberfx.RunApp(":3004", "example", 5*time.Second),
	)

	// Start the application
	app.Run()
}
