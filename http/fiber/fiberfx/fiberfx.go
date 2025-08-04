package fiberfx

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.uber.org/fx"

	corehttp "github.com/CodeLieutenant/uberfx-common/v3/http/fiber"
)

type (
	RoutesFx func(string) fx.Option
	RouteFx  func(string, string) fx.Option

	routerCallback struct {
		Callback func(fiber.Router)
		Prefix   string
	}
	routerCallbacks struct {
		cbs []routerCallback
	}
)

func RunApp(addr, appName string, shutdownTimeout time.Duration) fx.Option {
	return fx.Invoke(fx.Annotate(func(app *fiber.App, lc fx.Lifecycle) {
		lc.Append(fx.StartStopHook(
			func() {
				go func() {
					if err := app.Listen(addr); err != nil {
						panic(err)
					}
				}()
			},
			func(ctx context.Context) error {
				newCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
				defer cancel()
				return app.ShutdownWithContext(newCtx)
			},
		))
	}, fx.ParamTags(
		GetFiberApp(appName),
		`optional:"true"`,
	)))
}

func App(appName string, routes RoutesFx, options ...Option) fx.Option {
	opts := appOptions{
		cfg: corehttp.DefaultFiberConfig,
	}

	for _, opt := range options {
		opt(&opts)
	}

	var appProvide fx.Option

	// Define a struct to inject the middlewares
	type middlewareIn struct {
		fx.In

		Handlers    []route                `group:"fiber-handlers"`
		Callbacks   *routerCallbacks       `name:"fiber-router-callbacks"`
		Middlewares []middlewareWithPrefix `group:"fiber-middlewares"`
	}

	if opts.useMiddlewares {
		// If middleware injection is enabled, include middlewares in the app creation
		appProvide = fx.Provide(fx.Annotate(
			func(handlers []route, cbs *routerCallbacks, middlewares []middlewareWithPrefix) *fiber.App {
				app := corehttp.CreateApplication(opts.afterCreate, opts.cfg)

				// Apply middlewares first
				applyMiddlewares(app, middlewares)

				for _, r := range handlers {
					var router fiber.Router

					if r.Prefix != "" {
						router = app.Group(r.Prefix)
					} else {
						router = app
					}

					cb, exists := cbs.Find(r.Prefix)

					if exists {
						cb.Callback(router)
					}

					// Create a handler chain with route-specific middleware
					var handlers []fiber.Handler
					if len(r.Middlewares) > 0 {
						handlers = append(handlers, r.Middlewares...)
					}
					handlers = append(handlers, r.Handler)

					// Add the route with middleware
					var rt fiber.Router
					if len(handlers) == 1 {
						// No middleware, just add the handler directly
						rt = router.Add(r.Method, r.Path, r.Handler)
					} else {
						// Apply middleware to this specific route
						rt = router.Add(r.Method, r.Path, handlers...)
					}

					// Apply callback if provided
					if r.CallBack != nil {
						r.CallBack(rt)
					}
				}

				return app
			},
			fx.ParamTags(
				fiberHandlerRoutes(appName),
				routerCallbacksName(appName),
				`group:"fiber-middlewares"`,
			),
			fx.ResultTags(GetFiberApp(appName)),
		))
	} else {
		// Backward compatibility: don't include middlewares
		appProvide = fx.Provide(fx.Annotate(
			func(handlers []route, cbs *routerCallbacks) *fiber.App {
				app := corehttp.CreateApplication(opts.afterCreate, opts.cfg)

				for _, r := range handlers {
					var router fiber.Router

					if r.Prefix != "" {
						router = app.Group(r.Prefix)
					} else {
						router = app
					}

					cb, exists := cbs.Find(r.Prefix)

					if exists {
						cb.Callback(router)
					}

					// Create a handler chain with route-specific middleware
					var handlers []fiber.Handler
					if len(r.Middlewares) > 0 {
						handlers = append(handlers, r.Middlewares...)
					}
					handlers = append(handlers, r.Handler)

					// Add the route with middleware
					var rt fiber.Router
					if len(handlers) == 1 {
						// No middleware, just add the handler directly
						rt = router.Add(r.Method, r.Path, r.Handler)
					} else {
						// Apply middleware to this specific route
						rt = router.Add(r.Method, r.Path, handlers...)
					}

					// Apply callback if provided
					if r.CallBack != nil {
						r.CallBack(rt)
					}
				}

				return app
			},
			fx.ParamTags(
				fiberHandlerRoutes(appName),
				routerCallbacksName(appName),
			),
			fx.ResultTags(GetFiberApp(appName)),
		))
	}

	return fx.Module("fiber-"+appName,
		fx.Supply(fx.Annotate(
			&routerCallbacks{
				cbs: make([]routerCallback, 0),
			},
			fx.ResultTags(routerCallbacksName(appName)),
		)),
		routes(appName),
		appProvide,
	)
}

func (c *routerCallbacks) Add(prefix string, cb func(fiber.Router)) {
	c.cbs = append(c.cbs, routerCallback{
		Prefix:   prefix,
		Callback: cb,
	})
}

func (c *routerCallbacks) Get() []routerCallback {
	return c.cbs
}

func (c *routerCallbacks) Find(prefix string) (routerCallback, bool) {
	return lo.Find(c.cbs, func(item routerCallback) bool {
		return item.Prefix == prefix
	})
}
