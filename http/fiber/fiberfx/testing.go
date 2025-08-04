//go:build testing

package fiberfx

import (
	"github.com/gofiber/fiber/v2"
)

// ApplyOption applies an option to the appOptions and returns the modified appOptions
// This is used for testing option functions
func ApplyOption(opt Option) *appOptions {
	opts := &appOptions{
		cfg: fiber.Config{},
	}
	opt(opts)
	return opts
}

// ApplyOptions applies multiple options in sequence to the appOptions and returns the modified appOptions
// This is used for testing option functions
func ApplyOptions(opts ...Option) *appOptions {
	appOpts := &appOptions{
		cfg: fiber.Config{},
	}
	for _, opt := range opts {
		opt(appOpts)
	}
	return appOpts
}

// GetAfterCreate returns the afterCreate function from appOptions
// This is used for testing option functions
func GetAfterCreate(opts *appOptions) func(app *fiber.App) {
	return opts.afterCreate
}

// GetConfig returns the cfg from appOptions
// This is used for testing option functions
func GetConfig(opts *appOptions) fiber.Config {
	return opts.cfg
}

// GetUseMiddlewares returns the useMiddlewares flag from appOptions
// This is used for testing option functions
func GetUseMiddlewares(opts *appOptions) bool {
	return opts.useMiddlewares
}

// RouteTestHandler is a simple handler for testing routes
func RouteTestHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("Test")
	}
}

// ExportedRoute is an exported version of the route struct for testing
type ExportedRoute struct {
	Method      string
	Path        string
	Prefix      string
	HasCallback bool
	HasHandler  bool
	Middlewares int
}

// CreateExportedRoute creates an ExportedRoute from a route
func CreateExportedRoute(r route) ExportedRoute {
	return ExportedRoute{
		Method:      r.Method,
		Path:        r.Path,
		Prefix:      r.Prefix,
		HasCallback: r.CallBack != nil,
		HasHandler:  r.Handler != nil,
		Middlewares: len(r.Middlewares),
	}
}

// MiddlewareTracker holds a middleware and a function to check if it was called
type MiddlewareTracker struct {
	// Middleware is the fiber middleware handler
	Middleware fiber.Handler
	// WasCalled returns true if the middleware was called
	WasCalled func() bool
}

// CreateTestMiddleware creates a test middleware that tracks if it was called
// This helps avoid race conditions in parallel tests by giving each test its own tracker
func CreateTestMiddleware() MiddlewareTracker {
	var called bool
	middleware := func(c *fiber.Ctx) error {
		called = true
		return c.Next()
	}

	wasCalled := func() bool {
		return called
	}

	return MiddlewareTracker{
		Middleware: middleware,
		WasCalled:  wasCalled,
	}
}
