package fiberfx

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.uber.org/fx"

	corehttp "github.com/dmalusev/uberfx-common/fiber"
)

type (
	RoutesFx func(appName string) fx.Option
	RouteFx  func(appName, prefix string) fx.Option

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

func App(appName string, displayInfo bool, routes RoutesFx, errorHandler fiber.ErrorHandler) fx.Option {
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
				app := corehttp.CreateApplication(appName, displayInfo, logger, errorHandler)

				for _, r := range handlers {
					router := app.Group(r.Prefix)

					cb, exists := cb[r.Prefix]
					if exists {
						cb(router)
					}

					route := router.Add(r.Method, r.Path, r.Handler)
					if r.CallBack != nil {
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

func Get(path string, handler any) RouteFx {
	return Route(http.MethodGet, path, handler)
}

func GetWithRouterCallback(path string, cb func(fiber.Router), handler any) RouteFx {
	return RouteWithRouterCallback(http.MethodGet, path, cb, handler)
}

func Post(path string, handler any) RouteFx {
	return Route(http.MethodPost, path, handler)
}

func PostWithRouterCallback(path string, cb func(fiber.Router), handler any) RouteFx {
	return RouteWithRouterCallback(http.MethodPost, path, cb, handler)
}

func Put(path string, handler any) RouteFx {
	return Route(http.MethodPut, path, handler)
}

func PutWithRouterCallback(path string, cb func(fiber.Router), handler any) RouteFx {
	return RouteWithRouterCallback(http.MethodPut, path, cb, handler)
}

func Patch(path string, handler any) RouteFx {
	return Route(http.MethodPatch, path, handler)
}

func PatchWithRouterCallback(path string, cb func(fiber.Router), handler any) RouteFx {
	return RouteWithRouterCallback(http.MethodPatch, path, cb, handler)
}

func Delete(path string, handler any) RouteFx {
	return Route(http.MethodDelete, path, handler)
}

func DeleteWithRouterCallback(path string, cb func(fiber.Router), handler any) RouteFx {
	return RouteWithRouterCallback(http.MethodDelete, path, cb, handler)
}

func Route(method, path string, handler any) RouteFx {
	return RouteWithRouterCallback(method, path, nil, handler)
}

type route struct {
	Handler  fiber.Handler
	CallBack func(fiber.Router)
	Prefix   string
	Method   string
	Path     string
}

func RouteWithRouterCallback(method, path string, cb func(fiber.Router), handler any) RouteFx {
	return func(appName, prefix string) fx.Option {
		return fx.Provide(
			fx.Annotate(
				handler,
				fx.ResultTags(fiberHandlers(appName, method, prefix, path)),
			),
			fx.Annotate(
				func(handler fiber.Handler) route {
					return route{
						Prefix:   prefix,
						Method:   method,
						Path:     path,
						Handler:  handler,
						CallBack: cb,
					}
				},
				fx.ParamTags(fiberHandlers(appName, method, prefix, path)),
				fx.ResultTags(fiberHandlerRoutes(appName)),
			),
		)
	}
}

func fiberHandlers(appName, method, prefix, path string) string {
	return fmt.Sprintf(`name:"fiber-handler-%s-%s-%s-%s"`, appName, method, prefix, path)
}

func fiberHandlerRoutes(appName string) string {
	return fmt.Sprintf(`group:"fiber-handlers-%s"`, appName)
}

func GetFiberApp(appName string) string {
	return fmt.Sprintf(`name:"fiber-%s"`, appName)
}

type (
	OptionsFunc func(*options)

	options struct {
		prefix string
		cb     func(fiber.Router)
		routes []RouteFx
	}
)

func WithPrefix(prefix string) OptionsFunc {
	return func(opt *options) {
		opt.prefix = prefix
	}
}

func WithRoutes(routes ...RouteFx) OptionsFunc {
	return func(opt *options) {
		opt.routes = routes
	}
}

func WithRouterCallback(cb func(fiber.Router)) OptionsFunc {
	return func(opt *options) {
		opt.cb = cb
	}
}

func Routes(opts ...OptionsFunc) RoutesFx {
	opt := options{}

	for _, o := range opts {
		o(&opt)
	}

	return func(appName string) fx.Option {
		options := lo.Map(opt.routes, func(fn RouteFx, _ int) fx.Option {
			return fn(appName, opt.prefix)
		})

		if opt.cb == nil {
			return fx.Options(options...)
		}

		return fx.Options(append(options, fx.Invoke(fx.Annotate(
			func(callbacks routerCallbacks) {
				callbacks[opt.prefix] = opt.cb
			},
			fx.ParamTags(routerCallbacksName(appName)),
		)))...)
	}
}
