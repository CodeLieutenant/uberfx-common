package fiberfx_test

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
)

// TestRoute tests the basic Route function
func TestRoute(t *testing.T) {
	t.Parallel()

	// Create a test handler
	handler := fiberfx.RouteTestHandler

	// Test with GET method
	t.Run("GET method", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		// Create a route using the Route function
		routeFx := fiberfx.Route(http.MethodGet, "/test", handler)

		// Verify that routeFx is a function
		assert.NotNil(routeFx)

		// Call routeFx with app name and prefix
		option := routeFx("testapp", "")

		// Verify that option is an fx.Option
		_, ok := option.(fx.Option)
		assert.True(ok)
	})

	// Test with POST method
	t.Run("POST method", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		// Create a route using the Route function
		routeFx := fiberfx.Route(http.MethodPost, "/test", handler)

		// Verify that routeFx is a function
		assert.NotNil(routeFx)

		// Call routeFx with app name and prefix
		option := routeFx("testapp", "")

		// Verify that option is an fx.Option
		_, ok := option.(fx.Option)
		assert.True(ok)
	})
}

// TestRouteWithRouterCallback tests the RouteWithRouterCallback function
func TestRouteWithRouterCallback(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Create a test handler
	handler := fiberfx.RouteTestHandler

	// Create a test router callback
	callback := func(router fiber.Router) {
		// Do nothing in the test
	}

	// Create a route using the RouteWithRouterCallback function
	routeFx := fiberfx.RouteWithRouterCallback(http.MethodGet, "/test", callback, handler)

	// Verify that routeFx is a function
	assert.NotNil(routeFx)

	// Call routeFx with app name and prefix
	option := routeFx("testapp", "")

	// Verify that option is an fx.Option
	_, ok := option.(fx.Option)
	assert.True(ok)
}

// TestRouteWithMiddleware tests the RouteWithMiddleware function
func TestRouteWithMiddleware(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Create a test handler
	handler := fiberfx.RouteTestHandler

	// Create test middleware
	middleware := func(c *fiber.Ctx) error {
		return c.Next()
	}

	// Create a route using the RouteWithMiddleware function
	routeFx := fiberfx.RouteWithMiddleware(http.MethodGet, "/test", nil, []fiber.Handler{middleware}, handler)

	// Verify that routeFx is a function
	assert.NotNil(routeFx)

	// Call routeFx with app name and prefix
	option := routeFx("testapp", "")

	// Verify that option is an fx.Option
	_, ok := option.(fx.Option)
	assert.True(ok)
}

// TestHTTPMethodFunctions tests the HTTP method-specific functions
func TestHTTPMethodFunctions(t *testing.T) {
	t.Parallel()

	// Create a test handler
	handler := fiberfx.RouteTestHandler

	// Test Get function
	t.Run("Get", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		routeFx := fiberfx.Get("/test", handler)
		assert.NotNil(routeFx)

		option := routeFx("testapp", "")
		_, ok := option.(fx.Option)
		assert.True(ok)
	})

	// Test Post function
	t.Run("Post", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		routeFx := fiberfx.Post("/test", handler)
		assert.NotNil(routeFx)

		option := routeFx("testapp", "")
		_, ok := option.(fx.Option)
		assert.True(ok)
	})

	// Test Put function
	t.Run("Put", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		routeFx := fiberfx.Put("/test", handler)
		assert.NotNil(routeFx)

		option := routeFx("testapp", "")
		_, ok := option.(fx.Option)
		assert.True(ok)
	})

	// Test Patch function
	t.Run("Patch", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		routeFx := fiberfx.Patch("/test", handler)
		assert.NotNil(routeFx)

		option := routeFx("testapp", "")
		_, ok := option.(fx.Option)
		assert.True(ok)
	})

	// Test Delete function
	t.Run("Delete", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		routeFx := fiberfx.Delete("/test", handler)
		assert.NotNil(routeFx)

		option := routeFx("testapp", "")
		_, ok := option.(fx.Option)
		assert.True(ok)
	})
}

// TestHTTPMethodWithCallbackFunctions tests the HTTP method-specific functions with router callbacks
func TestHTTPMethodWithCallbackFunctions(t *testing.T) {
	t.Parallel()

	// Create a test handler
	handler := fiberfx.RouteTestHandler

	// Create a test router callback
	callback := func(router fiber.Router) {
		// Do nothing in the test
	}

	// Test GetWithRouterCallback function
	t.Run("GetWithRouterCallback", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		routeFx := fiberfx.GetWithRouterCallback("/test", callback, handler)
		assert.NotNil(routeFx)

		option := routeFx("testapp", "")
		_, ok := option.(fx.Option)
		assert.True(ok)
	})

	// Test PostWithRouterCallback function
	t.Run("PostWithRouterCallback", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		routeFx := fiberfx.PostWithRouterCallback("/test", callback, handler)
		assert.NotNil(routeFx)

		option := routeFx("testapp", "")
		_, ok := option.(fx.Option)
		assert.True(ok)
	})
}

// TestHTTPMethodWithMiddlewareFunctions tests the HTTP method-specific functions with middleware
func TestHTTPMethodWithMiddlewareFunctions(t *testing.T) {
	t.Parallel()

	// Create a test handler
	handler := fiberfx.RouteTestHandler

	// Create test middleware
	middleware := func(c *fiber.Ctx) error {
		return c.Next()
	}

	// Test GetWithMiddleware function
	t.Run("GetWithMiddleware", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		routeFx := fiberfx.GetWithMiddleware("/test", []fiber.Handler{middleware}, handler)
		assert.NotNil(routeFx)

		option := routeFx("testapp", "")
		_, ok := option.(fx.Option)
		assert.True(ok)
	})

	// Test PostWithMiddleware function
	t.Run("PostWithMiddleware", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		routeFx := fiberfx.PostWithMiddleware("/test", []fiber.Handler{middleware}, handler)
		assert.NotNil(routeFx)

		option := routeFx("testapp", "")
		_, ok := option.(fx.Option)
		assert.True(ok)
	})
}

// TestHTTPMethodWithCallbackAndMiddlewareFunctions tests the HTTP method-specific functions with router callbacks and middleware
func TestHTTPMethodWithCallbackAndMiddlewareFunctions(t *testing.T) {
	t.Parallel()

	// Create a test handler
	handler := fiberfx.RouteTestHandler

	// Create a test router callback
	callback := func(router fiber.Router) {
		// Do nothing in the test
	}

	// Create test middleware
	middleware := func(c *fiber.Ctx) error {
		return c.Next()
	}

	// Test GetWithRouterCallbackAndMiddleware function
	t.Run("GetWithRouterCallbackAndMiddleware", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		routeFx := fiberfx.GetWithRouterCallbackAndMiddleware("/test", callback, []fiber.Handler{middleware}, handler)
		assert.NotNil(routeFx)

		option := routeFx("testapp", "")
		_, ok := option.(fx.Option)
		assert.True(ok)
	})

	// Test PostWithRouterCallbackAndMiddleware function
	t.Run("PostWithRouterCallbackAndMiddleware", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		routeFx := fiberfx.PostWithRouterCallbackAndMiddleware("/test", callback, []fiber.Handler{middleware}, handler)
		assert.NotNil(routeFx)

		option := routeFx("testapp", "")
		_, ok := option.(fx.Option)
		assert.True(ok)
	})
}
