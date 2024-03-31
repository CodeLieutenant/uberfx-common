package fiberfx

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

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
