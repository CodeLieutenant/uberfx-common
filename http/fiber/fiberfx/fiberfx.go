package fiberfx

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.uber.org/fx"

	corehttp "github.com/CodeLieutenant/uberfx-common/v3/http/fiber"
)

type (
	RoutesFx        func(string) fx.Option
	RouteFx         func(string, string) fx.Option
	routerCallbacks map[string]func(fiber.Router)
)

func routerCallbacksName(appName string) string {
	return fmt.Sprintf(`name:"fiber-%s-router-callbacks"`, appName)
}

func RunApp(addr, appName string, shutdownTimeout time.Duration) fx.Option {
	return fx.Invoke(fx.Annotate(func(app *fiber.App, logger zerolog.Logger, lc fx.Lifecycle) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Info().Str("app", appName).Msg("Starting Fiber Application")
				go func() { _ = app.Listen(addr) }()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				newCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
				defer cancel()
				logger.Info().Str("app", appName).Msg("Stopping Fiber Application")
				return app.ShutdownWithContext(newCtx)
			},
		})
	}, fx.ParamTags(
		GetFiberApp(appName),
		`optional:"true"`,
	)))
}

func GetFiberApp(appName string) string {
	return fmt.Sprintf(`name:"fiber-%s"`, appName)
}

func App(appName string, routes RoutesFx, options ...Option) fx.Option {
	opts := appOptions{
		cfg: corehttp.DefaultFiberConfig,
	}

	for _, opt := range options {
		opt(&opts)
	}

	return fx.Module(fmt.Sprintf("fiber-%s", appName),
		fx.Provide(fx.Annotate(
			func() routerCallbacks {
				return make(routerCallbacks)
			},
			fx.ResultTags(routerCallbacksName(appName)),
		)),
		routes(appName),
		fx.Provide(fx.Annotate(
			func(logger zerolog.Logger, handlers []route, cb routerCallbacks, lc fx.Lifecycle) *fiber.App {
				app := corehttp.CreateApplication(opts.afterCreate, opts.cfg)

				for _, r := range handlers {
					var router fiber.Router

					if r.Prefix != "" {
						router = app.Group(r.Prefix)
					} else {
						router = app
					}

					if cb, exists := cb[r.Prefix]; exists {
						cb(router)
					}

					if route := router.Add(r.Method, r.Path, r.Handler); r.CallBack != nil {
						r.CallBack(route)
					}
				}

				return app
			},
			fx.ParamTags(
				`optional:"true"`,
				fiberHandlerRoutes(appName),
				routerCallbacksName(appName),
			),
			fx.ResultTags(GetFiberApp(appName)),
		)),
	)
}
