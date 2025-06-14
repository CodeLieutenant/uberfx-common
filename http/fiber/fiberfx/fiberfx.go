package fiberfx

import (
	"context"
	"time"

	corehttp "github.com/CodeLieutenant/uberfx-common/v3/http/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.uber.org/fx"
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

	return fx.Module("fiber-"+appName,
		fx.Supply(fx.Annotate(
			&routerCallbacks{
				cbs: make([]routerCallback, 0),
			},
			fx.ResultTags(routerCallbacksName(appName)),
		)),
		routes(appName),
		fx.Provide(fx.Annotate(
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

					if rt := router.Add(r.Method, r.Path, r.Handler); r.CallBack != nil {
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
		)),
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
