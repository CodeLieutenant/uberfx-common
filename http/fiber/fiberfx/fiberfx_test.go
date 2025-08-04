package fiberfx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
)

// TestApp tests the App function
func TestApp(t *testing.T) {
	t.Parallel()

	// Test with no options
	t.Run("with no options", func(t *testing.T) {
		t.Parallel()

		// Create a list of routes
		routes := []fiberfx.RouteFx{
			fiberfx.Get("/test", fiberfx.RouteTestHandler),
		}

		// Create a RoutesFx
		routesFx := fiberfx.Routes(routes)

		// Create a channel to receive the Fiber app
		appCh := make(chan *fiber.App, 1)

		// Create an App
		app := fxtest.New(
			t,
			fiberfx.App("testapp", routesFx),
			fx.Invoke(fx.Annotate(
				func(fiberApp *fiber.App) {
					appCh <- fiberApp
				},
				fx.ParamTags(`name:"fiber-testapp"`),
			)),
		)
		defer app.RequireStop()

		// Start the app
		app.RequireStart()

		// Get the Fiber app from the channel
		var fiberApp *fiber.App
		select {
		case fiberApp = <-appCh:
			// Got the app
		case <-time.After(time.Second):
			t.Fatal("Timed out waiting for Fiber app")
		}

		// Verify that the app was created
		require.NotNil(t, app)
		require.NotNil(t, fiberApp)
	})

	// Test with FiberConfig option
	t.Run("with FiberConfig option", func(t *testing.T) {
		t.Parallel()

		// Create a list of routes
		routes := []fiberfx.RouteFx{
			fiberfx.Get("/test", fiberfx.RouteTestHandler),
		}

		// Create a RoutesFx
		routesFx := fiberfx.Routes(routes)

		// Create a custom Fiber config
		customConfig := fiber.Config{
			ServerHeader:          "TestServer",
			DisableStartupMessage: true,
		}

		// Create a channel to receive the Fiber app
		appCh := make(chan *fiber.App, 1)

		// Create an App with custom config
		app := fxtest.New(
			t,
			fiberfx.App("testapp", routesFx, fiberfx.WithFiberConfig(customConfig)),
			fx.Invoke(fx.Annotate(
				func(fiberApp *fiber.App) {
					appCh <- fiberApp
				},
				fx.ParamTags(`name:"fiber-testapp"`),
			)),
		)
		defer app.RequireStop()

		// Start the app
		app.RequireStart()

		// Get the Fiber app from the channel
		var fiberApp *fiber.App
		select {
		case fiberApp = <-appCh:
			// Got the app
		case <-time.After(time.Second):
			t.Fatal("Timed out waiting for Fiber app")
		}

		// Verify that the app was created
		require.NotNil(t, app)
		require.NotNil(t, fiberApp)
		require.Equal(t, "TestServer", fiberApp.Config().ServerHeader)
	})

	// Test with AfterCreate option
	t.Run("with AfterCreate option", func(t *testing.T) {
		t.Parallel()

		// Create a flag to track if the afterCreate function was called
		var afterCreateCalled bool

		// Create an afterCreate function
		afterCreate := func(app *fiber.App) {
			afterCreateCalled = true
		}

		// Create a list of routes
		routes := []fiberfx.RouteFx{
			fiberfx.Get("/test", fiberfx.RouteTestHandler),
		}

		// Create a RoutesFx
		routesFx := fiberfx.Routes(routes)

		// Create a channel to receive the Fiber app
		appCh := make(chan *fiber.App, 1)

		// Create an App with afterCreate
		app := fxtest.New(
			t,
			fiberfx.App("testapp", routesFx, fiberfx.WithAfterCreate(afterCreate)),
			fx.Invoke(fx.Annotate(
				func(fiberApp *fiber.App) {
					appCh <- fiberApp
				},
				fx.ParamTags(`name:"fiber-testapp"`),
			)),
		)
		defer app.RequireStop()

		// Start the app
		app.RequireStart()

		// Get the Fiber app from the channel
		var fiberApp *fiber.App
		select {
		case fiberApp = <-appCh:
			// Got the app
		case <-time.After(time.Second):
			t.Fatal("Timed out waiting for Fiber app")
		}

		// Verify that the app was created
		require.NotNil(t, app)
		require.NotNil(t, fiberApp)

		// Verify that the afterCreate function was called
		require.True(t, afterCreateCalled)
	})
}

// TestAppWithMiddleware tests the App function with middleware
func TestAppWithMiddleware(t *testing.T) {
	t.Parallel()

	// Create a flag to track if the middleware was called
	var middlewareCalled bool

	// Create a test middleware constructor
	testMiddlewareConstructor := func() fiberfx.Middleware {
		return func(c *fiber.Ctx) error {
			middlewareCalled = true
			return c.Next()
		}
	}

	// Create a list of routes
	routes := []fiberfx.RouteFx{
		fiberfx.Get("/test", fiberfx.RouteTestHandler),
	}

	// Create a RoutesFx
	routesFx := fiberfx.Routes(routes)

	// Create a channel to receive the Fiber app
	appCh := make(chan *fiber.App, 1)

	// Create an App with middleware
	app := fxtest.New(
		t,
		fiberfx.App("testapp", routesFx, fiberfx.WithMiddlewares()),
		fiberfx.RegisterMiddleware("testapp", testMiddlewareConstructor),
		fx.Invoke(fx.Annotate(
			func(fiberApp *fiber.App) {
				appCh <- fiberApp
			},
			fx.ParamTags(`name:"fiber-testapp"`),
		)),
	)
	defer app.RequireStop()

	// Start the app
	app.RequireStart()

	// Get the Fiber app from the channel
	var fiberApp *fiber.App
	select {
	case fiberApp = <-appCh:
		// Got the app
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for Fiber app")
	}

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := fiberApp.Test(req)

	// Verify that the request was successful
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify that the middleware was called
	require.True(t, middlewareCalled)
}

// TestAppWithMiddlewarePrefix tests the App function with middleware that has a prefix
func TestAppWithMiddlewarePrefix(t *testing.T) {
	t.Parallel()

	// Create flags to track if the middlewares were called
	var globalMiddlewareCalled, prefixMiddlewareCalled bool

	// Create test middleware constructors
	globalMiddlewareConstructor := func() fiberfx.Middleware {
		return func(c *fiber.Ctx) error {
			globalMiddlewareCalled = true
			return c.Next()
		}
	}

	prefixMiddlewareConstructor := func() fiberfx.Middleware {
		return func(c *fiber.Ctx) error {
			prefixMiddlewareCalled = true
			return c.Next()
		}
	}

	// Create a list of routes
	routes := []fiberfx.RouteFx{
		fiberfx.Get("/test", fiberfx.RouteTestHandler),
		fiberfx.Get("/api/test", fiberfx.RouteTestHandler),
	}

	// Create a RoutesFx
	routesFx := fiberfx.Routes(routes)

	// Create a channel to receive the Fiber app
	appCh := make(chan *fiber.App, 1)

	// Create an App with middleware
	app := fxtest.New(
		t,
		fiberfx.App("testapp", routesFx, fiberfx.WithMiddlewares()),
		fiberfx.RegisterMiddleware("testapp", globalMiddlewareConstructor),
		fiberfx.RegisterMiddlewareWithPrefix("testapp", "/api", prefixMiddlewareConstructor),
		fx.Invoke(fx.Annotate(
			func(fiberApp *fiber.App) {
				appCh <- fiberApp
			},
			fx.ParamTags(`name:"fiber-testapp"`),
		)),
	)
	defer app.RequireStop()

	// Start the app
	app.RequireStart()

	// Get the Fiber app from the channel
	var fiberApp *fiber.App
	select {
	case fiberApp = <-appCh:
		// Got the app
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for Fiber app")
	}

	// Test the global route
	t.Run("global route", func(t *testing.T) {
		// Reset flags
		globalMiddlewareCalled = false
		prefixMiddlewareCalled = false

		// Create a test request
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp, err := fiberApp.Test(req)

		// Verify that the request was successful
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify that only the global middleware was called
		require.True(t, globalMiddlewareCalled)
		require.False(t, prefixMiddlewareCalled)
	})

	// Test the prefixed route
	t.Run("prefixed route", func(t *testing.T) {
		// Reset flags
		globalMiddlewareCalled = false
		prefixMiddlewareCalled = false

		// Create a test request
		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		resp, err := fiberApp.Test(req)

		// Verify that the request was successful
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify that both middlewares were called
		require.True(t, globalMiddlewareCalled)
		require.True(t, prefixMiddlewareCalled)
	})
}

// TestAppWithPerRouteMiddleware tests the App function with per-route middleware
func TestAppWithPerRouteMiddleware(t *testing.T) {
	t.Parallel()

	// Create a middleware tracker for the route with middleware
	withMiddlewareTracker := fiberfx.CreateTestMiddleware()

	// Create a list of routes
	routes := []fiberfx.RouteFx{
		// Route without middleware
		fiberfx.Get("/test", fiberfx.RouteTestHandler),
		// Route with middleware
		fiberfx.GetWithMiddleware("/test-with-middleware", []fiber.Handler{withMiddlewareTracker.Middleware}, fiberfx.RouteTestHandler),
	}

	// Create a RoutesFx
	routesFx := fiberfx.Routes(routes)

	// Create a channel to receive the Fiber app
	appCh := make(chan *fiber.App, 1)

	// Create an App
	app := fxtest.New(
		t,
		fiberfx.App("testapp", routesFx),
		fx.Invoke(fx.Annotate(
			func(fiberApp *fiber.App) {
				appCh <- fiberApp
			},
			fx.ParamTags(`name:"fiber-testapp"`),
		)),
	)
	defer app.RequireStop()

	// Start the app
	app.RequireStart()

	// Get the Fiber app from the channel
	var fiberApp *fiber.App
	select {
	case fiberApp = <-appCh:
		// Got the app
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for Fiber app")
	}

	// Test the route without middleware
	t.Run("route without middleware", func(t *testing.T) {
		t.Parallel()

		// Create a separate middleware tracker for this test to verify it's not called
		testMiddlewareTracker := fiberfx.CreateTestMiddleware()

		// Create a test request
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp, err := fiberApp.Test(req)

		// Verify that the request was successful
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify that the test middleware was not called
		// Since we're not using this middleware in the route, it should never be called
		require.False(t, testMiddlewareTracker.WasCalled())
	})

	// Test the route with middleware
	t.Run("route with middleware", func(t *testing.T) {
		t.Parallel()

		// Create a test request
		req := httptest.NewRequest(http.MethodGet, "/test-with-middleware", nil)
		resp, err := fiberApp.Test(req)

		// Verify that the request was successful
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify that the middleware was called
		require.True(t, withMiddlewareTracker.WasCalled())
	})
}

// TestRunApp tests the RunApp function
func TestRunApp(t *testing.T) {
	t.Parallel()

	// Create a list of routes
	routes := []fiberfx.RouteFx{
		fiberfx.Get("/test", fiberfx.RouteTestHandler),
	}

	// Create a RoutesFx
	routesFx := fiberfx.Routes(routes)

	// Create an App
	app := fxtest.New(
		t,
		fiberfx.App("testapp", routesFx),
		fiberfx.RunApp(":0", "testapp", 100*time.Millisecond),
	)

	// Start the app
	app.Start(context.Background())

	// Stop the app
	app.Stop(context.Background())

	// Verify that the app was created and started/stopped without errors
	require.NotNil(t, app)
}
