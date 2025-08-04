package fiberfx_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
)

// TestRouteWithMiddlewareFx tests the RouteWithMiddlewareFx function
func TestRouteWithMiddlewareFx(t *testing.T) {
	t.Parallel()

	// Create a test handler
	testHandler := func() fiber.Handler {
		return func(c *fiber.Ctx) error {
			return c.SendString("Hello, World!")
		}
	}

	// Create a test middleware with dependencies
	type TestDep struct {
		Value string
	}

	testMiddleware := func(dep TestDep) fiber.Handler {
		return func(c *fiber.Ctx) error {
			c.Locals("test_dep", dep.Value)
			return c.Next()
		}
	}

	// Test with GET method
	t.Run("GET with middleware fx", func(t *testing.T) {
		t.Parallel()

		// Create a route with middleware fx
		route := fiberfx.GetWithMiddlewareFx("/test", []fiberfx.RouteMiddlewareFunc{testMiddleware}, testHandler)

		// Verify that route is a RouteFx
		require.NotNil(t, route)

		// Create an fx app to test the route
		app := fxtest.New(
			t,
			fx.Supply(TestDep{Value: "test_value"}),
			route("testapp", ""),
		)
		require.NoError(t, app.Err())
	})

	// Test with POST method
	t.Run("POST with middleware fx", func(t *testing.T) {
		t.Parallel()

		// Create a route with middleware fx
		route := fiberfx.PostWithMiddlewareFx("/test", []fiberfx.RouteMiddlewareFunc{testMiddleware}, testHandler)

		// Verify that route is a RouteFx
		require.NotNil(t, route)

		// Create an fx app to test the route
		app := fxtest.New(
			t,
			fx.Supply(TestDep{Value: "test_value"}),
			route("testapp", ""),
		)
		require.NoError(t, app.Err())
	})

	// Test with router callback
	t.Run("with router callback", func(t *testing.T) {
		t.Parallel()

		// Create a router callback
		cb := func(router fiber.Router) {
			router.Use(func(c *fiber.Ctx) error {
				return c.Next()
			})
		}

		// Create a route with middleware fx and router callback
		route := fiberfx.GetWithRouterCallbackAndMiddlewareFx("/test", cb, []fiberfx.RouteMiddlewareFunc{testMiddleware}, testHandler)

		// Verify that route is a RouteFx
		require.NotNil(t, route)

		// Create an fx app to test the route
		app := fxtest.New(
			t,
			fx.Supply(TestDep{Value: "test_value"}),
			route("testapp", ""),
		)
		require.NoError(t, app.Err())
	})

	// Test with multiple middleware functions
	t.Run("with multiple middleware", func(t *testing.T) {
		t.Parallel()

		// Create another test middleware with dependencies
		type AnotherTestDep struct {
			Value string
		}

		anotherTestMiddleware := func(dep AnotherTestDep) fiber.Handler {
			return func(c *fiber.Ctx) error {
				c.Locals("another_test_dep", dep.Value)
				return c.Next()
			}
		}

		// Create a route with multiple middleware fx
		route := fiberfx.GetWithMiddlewareFx("/test", []fiberfx.RouteMiddlewareFunc{testMiddleware, anotherTestMiddleware}, testHandler)

		// Verify that route is a RouteFx
		require.NotNil(t, route)

		// Create an fx app to test the route
		app := fxtest.New(
			t,
			fx.Supply(TestDep{Value: "test_value"}, AnotherTestDep{Value: "another_test_value"}),
			route("testapp", ""),
		)
		require.NoError(t, app.Err())
	})
}
