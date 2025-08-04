package fiberfx

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type (
	// Middleware represents a Fiber middleware function
	Middleware fiber.Handler

	// MiddlewareWithPrefix represents a middleware with a specific prefix
	middlewareWithPrefix struct {
		Middleware fiber.Handler
		Prefix     string
	}

	// middlewareGroup is a collection of middlewares for a specific app
	middlewareGroup struct {
		middlewares []middlewareWithPrefix
	}

	// middlewareOut is used to register a middleware as part of a group
	middlewareOut struct {
		fx.Out

		Middleware middlewareWithPrefix `group:"fiber-middlewares"`
	}

	// RouteMiddlewareFunc represents a function that returns a Fiber middleware
	// This allows middleware to have dependencies injected by uberfx
	RouteMiddlewareFunc any
)

// RegisterMiddleware registers a middleware to be used with a specific app
func RegisterMiddleware(appName string, middleware any) fx.Option {
	return RegisterMiddlewareWithPrefix(appName, "", middleware)
}

// RegisterMiddlewareWithPrefix registers a middleware to be used with a specific app and route prefix
func RegisterMiddlewareWithPrefix(appName, prefix string, middleware any) fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				middleware,
				fx.ResultTags(middlewareTag(appName, prefix)),
			),
		),
		fx.Provide(
			fx.Annotate(
				func(m Middleware) middlewareOut {
					return middlewareOut{
						Middleware: middlewareWithPrefix{
							Middleware: m,
							Prefix:     prefix,
						},
					}
				},
				fx.ParamTags(middlewareTag(appName, prefix)),
			),
		),
	)
}

// middlewareTag generates a tag for a middleware
func middlewareTag(appName, prefix string) string {
	if prefix == "" {
		return fiberMiddlewareTag(appName, "global")
	}
	return fiberMiddlewareTag(appName, prefix)
}

// fiberMiddlewareTag generates a tag for a middleware
func fiberMiddlewareTag(appName, prefix string) string {
	return "name:\"fiber-middleware-" + appName + "-" + prefix + "\""
}

// middlewareGroupTag generates a tag for a middleware group
func middlewareGroupTag(appName string) string {
	return "group:\"fiber-middlewares-" + appName + "\""
}

// WithMiddlewares is an option that enables middleware injection for the app
func WithMiddlewares() Option {
	return func(opts *appOptions) {
		originalAfterCreate := opts.afterCreate

		opts.afterCreate = func(app *fiber.App) {
			// Call the original afterCreate function if it exists
			if originalAfterCreate != nil {
				originalAfterCreate(app)
			}
		}

		// Set the useMiddlewares flag to true
		opts.useMiddlewares = true
	}
}

// applyMiddlewares applies all registered middlewares to the app
func applyMiddlewares(app *fiber.App, middlewares []middlewareWithPrefix) {
	// Group middlewares by prefix
	prefixMap := make(map[string][]fiber.Handler)

	// Add global middlewares (empty prefix) first
	for _, m := range middlewares {
		if m.Prefix == "" {
			app.Use(m.Middleware)
		} else {
			if _, exists := prefixMap[m.Prefix]; !exists {
				prefixMap[m.Prefix] = make([]fiber.Handler, 0)
			}
			prefixMap[m.Prefix] = append(prefixMap[m.Prefix], m.Middleware)
		}
	}

	// Apply middlewares to specific prefixes
	for prefix, handlers := range prefixMap {
		for _, handler := range handlers {
			app.Use(prefix, handler)
		}
	}
}
