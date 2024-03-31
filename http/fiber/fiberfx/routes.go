package fiberfx

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

type (
	RouteOptions func(*routeOptions)

	routeOptions struct {
		cb     func(fiber.Router)
		prefix string
	}
)

func WithPrefix(prefix string) RouteOptions {
	return func(opt *routeOptions) {
		opt.prefix = prefix
	}
}

func WithRouterCallback(cb func(fiber.Router)) RouteOptions {
	return func(opt *routeOptions) {
		opt.cb = cb
	}
}

func Routes(routes []RouteFx, opts ...RouteOptions) RoutesFx {
	var opt routeOptions

	for _, o := range opts {
		o(&opt)
	}

	return func(appName string) fx.Option {
		options := lo.Map(routes, func(fn RouteFx, _ int) fx.Option {
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
