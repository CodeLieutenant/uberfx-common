package fiberfx_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
)

// TestServiceMiddleware tests middleware that expects a service in the wrapper function
func TestServiceMiddleware(t *testing.T) {
	t.Parallel()

	// Define a simple service
	type UserService struct {
		Username string
	}

	// Create a middleware that expects a service
	middlewareWithService := func(service UserService) fiber.Handler {
		return func(c *fiber.Ctx) error {
			// Add the username from the service to the context locals
			c.Locals("username", service.Username)
			return c.Next()
		}
	}

	// Create a test handler that reads the username from context locals
	testHandler := func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			username := c.Locals("username").(string)
			return c.SendString("Hello, " + username)
		}
	}

	// Create a route with middleware fx
	route := fiberfx.GetWithMiddlewareFx("/test", []fiberfx.RouteMiddlewareFunc{middlewareWithService}, testHandler)

	// Verify that route is a RouteFx
	require.NotNil(t, route)

	// Create an fx app to test the route
	app := fxtest.New(
		t,
		fx.Supply(UserService{Username: "testuser"}),
		route("testapp", ""),
	)
	require.NoError(t, app.Err())
}

// TestServiceRoute tests a route that expects a service in the wrapper function
func TestServiceRoute(t *testing.T) {
	t.Parallel()

	// Define a simple service
	type ProductService struct {
		Products []string
	}

	// Create a middleware that adds the products to context locals
	productsMiddleware := func(service ProductService) fiber.Handler {
		return func(c *fiber.Ctx) error {
			// Add the products from the service to the context locals
			c.Locals("products", service.Products)
			return c.Next()
		}
	}

	// Create a handler that gets the products from context locals
	productsHandler := func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			// Get the products from context locals
			products := c.Locals("products").([]string)
			return c.JSON(products)
		}
	}

	// Create a route with middleware fx
	route := fiberfx.GetWithMiddlewareFx("/products", []fiberfx.RouteMiddlewareFunc{productsMiddleware}, productsHandler)

	// Verify that route is a RouteFx
	require.NotNil(t, route)

	// Create an fx app to test the route
	app := fxtest.New(
		t,
		fx.Supply(ProductService{Products: []string{"Product1", "Product2", "Product3"}}),
		route("testapp", ""),
	)
	require.NoError(t, app.Err())
}

// TestCombinedServiceMiddlewareAndRoute tests both middleware and route that expect services
func TestCombinedServiceMiddlewareAndRoute(t *testing.T) {
	t.Parallel()

	// Define services
	type AuthService struct {
		Role string
	}

	type DataService struct {
		Data map[string]string
	}

	// Create a middleware that expects a service
	authMiddleware := func(service AuthService) fiber.Handler {
		return func(c *fiber.Ctx) error {
			// Add the role from the service to the context locals
			c.Locals("role", service.Role)
			return c.Next()
		}
	}

	// Create a middleware that adds data to context locals
	dataMiddleware := func(service DataService) fiber.Handler {
		return func(c *fiber.Ctx) error {
			// Add the data from the service to the context locals
			c.Locals("data", service.Data)
			return c.Next()
		}
	}

	// Create a handler that uses context locals from middleware
	dataHandler := func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			// Get the role and data from context locals
			role := c.Locals("role").(string)
			data := c.Locals("data").(map[string]string)

			// Return data based on role
			if role == "admin" {
				return c.JSON(data)
			}

			return c.Status(fiber.StatusForbidden).SendString("Access denied")
		}
	}

	// Create a route with middleware fx
	route := fiberfx.GetWithMiddlewareFx("/data", []fiberfx.RouteMiddlewareFunc{authMiddleware, dataMiddleware}, dataHandler)

	// Verify that route is a RouteFx
	require.NotNil(t, route)

	// Create an fx app to test the route
	app := fxtest.New(
		t,
		fx.Supply(
			AuthService{Role: "admin"},
			DataService{Data: map[string]string{"key1": "value1", "key2": "value2"}},
		),
		route("testapp", ""),
	)
	require.NoError(t, app.Err())
}
