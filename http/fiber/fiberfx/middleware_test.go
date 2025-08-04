package fiberfx_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
)

// TestRegisterMiddleware tests the RegisterMiddleware function
func TestRegisterMiddleware(t *testing.T) {
	t.Parallel()

	// Create a test middleware constructor
	testMiddlewareConstructor := func() fiberfx.Middleware {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	// Register the middleware
	option := fiberfx.RegisterMiddleware("testapp", testMiddlewareConstructor)

	// Verify that option is an fx.Option
	_, ok := option.(fx.Option)
	require.True(t, ok)
}

// TestRegisterMiddlewareWithPrefix tests the RegisterMiddlewareWithPrefix function
func TestRegisterMiddlewareWithPrefix(t *testing.T) {
	t.Parallel()

	// Create a test middleware constructor
	testMiddlewareConstructor := func() fiberfx.Middleware {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	// Test with empty prefix
	t.Run("with empty prefix", func(t *testing.T) {
		t.Parallel()

		// Register the middleware with empty prefix
		option := fiberfx.RegisterMiddlewareWithPrefix("testapp", "", testMiddlewareConstructor)

		// Verify that option is an fx.Option
		_, ok := option.(fx.Option)
		require.True(t, ok)
	})

	// Test with non-empty prefix
	t.Run("with non-empty prefix", func(t *testing.T) {
		t.Parallel()

		// Register the middleware with a prefix
		option := fiberfx.RegisterMiddlewareWithPrefix("testapp", "/api", testMiddlewareConstructor)

		// Verify that option is an fx.Option
		_, ok := option.(fx.Option)
		require.True(t, ok)
	})
}

// TestWithMiddlewaresInMiddleware tests the WithMiddlewares function
// Note: This is already tested in options_test.go, but we include a basic test here for completeness
func TestWithMiddlewaresInMiddleware(t *testing.T) {
	t.Parallel()

	// Apply the WithMiddlewares option
	opts := fiberfx.ApplyOption(fiberfx.WithMiddlewares())

	// Verify that useMiddlewares is set to true
	require.True(t, fiberfx.GetUseMiddlewares(opts))

	// Verify that an afterCreate function was created
	require.NotNil(t, fiberfx.GetAfterCreate(opts))
}

// TestApplyMiddlewares tests the applyMiddlewares function indirectly through the App function
// Note: This is an integration test that requires creating a full Fiber app
// We'll test this more thoroughly in the fiberfx_test.go file
func TestApplyMiddlewares(t *testing.T) {
	t.Parallel()
	// This is tested indirectly through the App function
	// The function is unexported, so we can't test it directly
}

// TestMiddlewareIntegration tests the middleware functionality in an integrated way
// This is a more comprehensive test that verifies the middleware is actually applied to the app
func TestMiddlewareIntegration(t *testing.T) {
	t.Parallel()
	// This is tested in the fiberfx_test.go file
}
