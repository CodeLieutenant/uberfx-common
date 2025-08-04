package fiberfx_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
)

// TestWithPrefix tests the WithPrefix function
func TestWithPrefix(t *testing.T) {
	t.Parallel()

	// Create a test handler
	handler := fiberfx.RouteTestHandler

	// Create a list of routes
	routes := []fiberfx.RouteFx{
		fiberfx.Get("/test", handler),
		fiberfx.Post("/test", handler),
	}

	// Create a RoutesFx with a prefix
	routesFx := fiberfx.Routes(routes, fiberfx.WithPrefix("/api"))

	// Verify that routesFx is a function
	require.NotNil(t, routesFx)

	// Call routesFx with app name
	option := routesFx("testapp")

	// Verify that option is an fx.Option
	_, ok := option.(fx.Option)
	require.True(t, ok)
}

// TestWithRouterCallback tests the WithRouterCallback function
func TestWithRouterCallback(t *testing.T) {
	t.Parallel()

	// Create a test handler
	handler := fiberfx.RouteTestHandler

	// Create a list of routes
	routes := []fiberfx.RouteFx{
		fiberfx.Get("/test", handler),
		fiberfx.Post("/test", handler),
	}

	// Create a test router callback
	callback := func(router fiber.Router) {
		// Do nothing in the test
	}

	// Create a RoutesFx with a router callback
	routesFx := fiberfx.Routes(routes, fiberfx.WithRouterCallback(callback))

	// Verify that routesFx is a function
	require.NotNil(t, routesFx)

	// Call routesFx with app name
	option := routesFx("testapp")

	// Verify that option is an fx.Option
	_, ok := option.(fx.Option)
	require.True(t, ok)
}

// TestRoutes tests the Routes function
func TestRoutes(t *testing.T) {
	t.Parallel()

	// Create a test handler
	handler := fiberfx.RouteTestHandler

	// Test with no options
	t.Run("with no options", func(t *testing.T) {
		t.Parallel()

		// Create a list of routes
		routes := []fiberfx.RouteFx{
			fiberfx.Get("/test", handler),
			fiberfx.Post("/test", handler),
		}

		// Create a RoutesFx with no options
		routesFx := fiberfx.Routes(routes)

		// Verify that routesFx is a function
		require.NotNil(t, routesFx)

		// Call routesFx with app name
		option := routesFx("testapp")

		// Verify that option is an fx.Option
		_, ok := option.(fx.Option)
		require.True(t, ok)
	})

	// Test with prefix option
	t.Run("with prefix option", func(t *testing.T) {
		t.Parallel()

		// Create a list of routes
		routes := []fiberfx.RouteFx{
			fiberfx.Get("/test", handler),
			fiberfx.Post("/test", handler),
		}

		// Create a RoutesFx with a prefix
		routesFx := fiberfx.Routes(routes, fiberfx.WithPrefix("/api"))

		// Verify that routesFx is a function
		require.NotNil(t, routesFx)

		// Call routesFx with app name
		option := routesFx("testapp")

		// Verify that option is an fx.Option
		_, ok := option.(fx.Option)
		require.True(t, ok)
	})

	// Test with router callback option
	t.Run("with router callback option", func(t *testing.T) {
		t.Parallel()

		// Create a list of routes
		routes := []fiberfx.RouteFx{
			fiberfx.Get("/test", handler),
			fiberfx.Post("/test", handler),
		}

		// Create a test router callback
		callback := func(router fiber.Router) {
			// Do nothing in the test
		}

		// Create a RoutesFx with a router callback
		routesFx := fiberfx.Routes(routes, fiberfx.WithRouterCallback(callback))

		// Verify that routesFx is a function
		require.NotNil(t, routesFx)

		// Call routesFx with app name
		option := routesFx("testapp")

		// Verify that option is an fx.Option
		_, ok := option.(fx.Option)
		require.True(t, ok)
	})

	// Test with both prefix and router callback options
	t.Run("with both prefix and router callback options", func(t *testing.T) {
		t.Parallel()

		// Create a list of routes
		routes := []fiberfx.RouteFx{
			fiberfx.Get("/test", handler),
			fiberfx.Post("/test", handler),
		}

		// Create a test router callback
		callback := func(router fiber.Router) {
			// Do nothing in the test
		}

		// Create a RoutesFx with both prefix and router callback
		routesFx := fiberfx.Routes(
			routes,
			fiberfx.WithPrefix("/api"),
			fiberfx.WithRouterCallback(callback),
		)

		// Verify that routesFx is a function
		require.NotNil(t, routesFx)

		// Call routesFx with app name
		option := routesFx("testapp")

		// Verify that option is an fx.Option
		_, ok := option.(fx.Option)
		require.True(t, ok)
	})
}

// TestCombineRoutes tests the CombineRoutes function
func TestCombineRoutes(t *testing.T) {
	t.Parallel()

	// Create a test handler
	handler := fiberfx.RouteTestHandler

	// Create lists of routes
	routes1 := []fiberfx.RouteFx{
		fiberfx.Get("/test1", handler),
		fiberfx.Post("/test1", handler),
	}

	routes2 := []fiberfx.RouteFx{
		fiberfx.Get("/test2", handler),
		fiberfx.Post("/test2", handler),
	}

	// Create RoutesFx functions
	routesFx1 := fiberfx.Routes(routes1, fiberfx.WithPrefix("/api1"))
	routesFx2 := fiberfx.Routes(routes2, fiberfx.WithPrefix("/api2"))

	// Combine the RoutesFx functions
	combinedRoutesFx := fiberfx.CombineRoutes(routesFx1, routesFx2)

	// Verify that combinedRoutesFx is a function
	require.NotNil(t, combinedRoutesFx)

	// Call combinedRoutesFx with app name
	option := combinedRoutesFx("testapp")

	// Verify that option is an fx.Option
	_, ok := option.(fx.Option)
	require.True(t, ok)
}
