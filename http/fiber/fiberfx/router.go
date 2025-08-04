package fiberfx

import (
	"fmt"
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

func GetWithMiddleware(path string, middlewares []fiber.Handler, handler any) RouteFx {
	return RouteWithMiddleware(http.MethodGet, path, nil, middlewares, handler)
}

func GetWithRouterCallbackAndMiddleware(path string, cb func(fiber.Router), middlewares []fiber.Handler, handler any) RouteFx {
	return RouteWithMiddleware(http.MethodGet, path, cb, middlewares, handler)
}

func GetWithMiddlewareFx(path string, middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx {
	return RouteWithMiddlewareFx(http.MethodGet, path, nil, middlewareFuncs, handler)
}

func GetWithRouterCallbackAndMiddlewareFx(path string, cb func(fiber.Router), middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx {
	return RouteWithMiddlewareFx(http.MethodGet, path, cb, middlewareFuncs, handler)
}

func Post(path string, handler any) RouteFx {
	return Route(http.MethodPost, path, handler)
}

func PostWithRouterCallback(path string, cb func(fiber.Router), handler any) RouteFx {
	return RouteWithRouterCallback(http.MethodPost, path, cb, handler)
}

func PostWithMiddleware(path string, middlewares []fiber.Handler, handler any) RouteFx {
	return RouteWithMiddleware(http.MethodPost, path, nil, middlewares, handler)
}

func PostWithRouterCallbackAndMiddleware(path string, cb func(fiber.Router), middlewares []fiber.Handler, handler any) RouteFx {
	return RouteWithMiddleware(http.MethodPost, path, cb, middlewares, handler)
}

func PostWithMiddlewareFx(path string, middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx {
	return RouteWithMiddlewareFx(http.MethodPost, path, nil, middlewareFuncs, handler)
}

func PostWithRouterCallbackAndMiddlewareFx(path string, cb func(fiber.Router), middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx {
	return RouteWithMiddlewareFx(http.MethodPost, path, cb, middlewareFuncs, handler)
}

func Put(path string, handler any) RouteFx {
	return Route(http.MethodPut, path, handler)
}

func PutWithRouterCallback(path string, cb func(fiber.Router), handler any) RouteFx {
	return RouteWithRouterCallback(http.MethodPut, path, cb, handler)
}

func PutWithMiddleware(path string, middlewares []fiber.Handler, handler any) RouteFx {
	return RouteWithMiddleware(http.MethodPut, path, nil, middlewares, handler)
}

func PutWithRouterCallbackAndMiddleware(path string, cb func(fiber.Router), middlewares []fiber.Handler, handler any) RouteFx {
	return RouteWithMiddleware(http.MethodPut, path, cb, middlewares, handler)
}

func PutWithMiddlewareFx(path string, middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx {
	return RouteWithMiddlewareFx(http.MethodPut, path, nil, middlewareFuncs, handler)
}

func PutWithRouterCallbackAndMiddlewareFx(path string, cb func(fiber.Router), middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx {
	return RouteWithMiddlewareFx(http.MethodPut, path, cb, middlewareFuncs, handler)
}

func Patch(path string, handler any) RouteFx {
	return Route(http.MethodPatch, path, handler)
}

func PatchWithRouterCallback(path string, cb func(fiber.Router), handler any) RouteFx {
	return RouteWithRouterCallback(http.MethodPatch, path, cb, handler)
}

func PatchWithMiddleware(path string, middlewares []fiber.Handler, handler any) RouteFx {
	return RouteWithMiddleware(http.MethodPatch, path, nil, middlewares, handler)
}

func PatchWithRouterCallbackAndMiddleware(path string, cb func(fiber.Router), middlewares []fiber.Handler, handler any) RouteFx {
	return RouteWithMiddleware(http.MethodPatch, path, cb, middlewares, handler)
}

func PatchWithMiddlewareFx(path string, middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx {
	return RouteWithMiddlewareFx(http.MethodPatch, path, nil, middlewareFuncs, handler)
}

func PatchWithRouterCallbackAndMiddlewareFx(path string, cb func(fiber.Router), middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx {
	return RouteWithMiddlewareFx(http.MethodPatch, path, cb, middlewareFuncs, handler)
}

func Delete(path string, handler any) RouteFx {
	return Route(http.MethodDelete, path, handler)
}

func DeleteWithRouterCallback(path string, cb func(fiber.Router), handler any) RouteFx {
	return RouteWithRouterCallback(http.MethodDelete, path, cb, handler)
}

func DeleteWithMiddleware(path string, middlewares []fiber.Handler, handler any) RouteFx {
	return RouteWithMiddleware(http.MethodDelete, path, nil, middlewares, handler)
}

func DeleteWithRouterCallbackAndMiddleware(path string, cb func(fiber.Router), middlewares []fiber.Handler, handler any) RouteFx {
	return RouteWithMiddleware(http.MethodDelete, path, cb, middlewares, handler)
}

func DeleteWithMiddlewareFx(path string, middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx {
	return RouteWithMiddlewareFx(http.MethodDelete, path, nil, middlewareFuncs, handler)
}

func DeleteWithRouterCallbackAndMiddlewareFx(path string, cb func(fiber.Router), middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx {
	return RouteWithMiddlewareFx(http.MethodDelete, path, cb, middlewareFuncs, handler)
}

func Route(method, path string, handler any) RouteFx {
	return RouteWithRouterCallback(method, path, nil, handler)
}

type route struct {
	Handler     fiber.Handler
	CallBack    func(fiber.Router)
	Prefix      string
	Method      string
	Path        string
	Middlewares []fiber.Handler
}

func RouteWithRouterCallback(method, path string, cb func(fiber.Router), handler any) RouteFx {
	return RouteWithMiddleware(method, path, cb, nil, handler)
}

func RouteWithMiddleware(method, path string, cb func(fiber.Router), middlewares []fiber.Handler, handler any) RouteFx {
	return func(appName, prefix string) fx.Option {
		// Create a wrapper function that adapts the handler to a fiber.Handler
		handlerWrapper := func() fiber.Handler {
			// Convert the handler to a fiber.Handler
			if h, ok := handler.(fiber.Handler); ok {
				return h
			}
			// If it's already a fiber.Handler function, return it directly
			if h, ok := handler.(func(*fiber.Ctx) error); ok {
				return h
			}
			// If it's a function that returns a fiber.Handler, call it to get the handler
			if h, ok := handler.(func() fiber.Handler); ok {
				return h()
			}
			// If it's a function that returns a fiber.Handler function, call it to get the handler
			if h, ok := handler.(func() func(*fiber.Ctx) error); ok {
				return h()
			}
			// Otherwise, panic as we can't handle this type
			panic(fmt.Sprintf("handler must be a fiber.Handler or func(*fiber.Ctx) error, got %T", handler))
		}

 	// Call handlerWrapper to get the actual handler
	handler := handlerWrapper()

	return fx.Provide(
			// Provide the handler directly
			fx.Annotate(
				func() fiber.Handler {
					return handler
				},
				fx.ResultTags(fiberHandlers(appName, method, prefix, path)),
			),
			// Create the route with the handler
			fx.Annotate(
				func(handler fiber.Handler) route {
					return route{
						Prefix:      prefix,
						Method:      method,
						Path:        path,
						Handler:     handler,
						CallBack:    cb,
						Middlewares: middlewares,
					}
				},
				fx.ParamTags(fiberHandlers(appName, method, prefix, path)),
				fx.ResultTags(fiberHandlerRoutes(appName)),
			),
		)
	}
}

// RouteWithMiddlewareFx is similar to RouteWithMiddleware but allows middleware functions
// to have dependencies injected by uberfx
func RouteWithMiddlewareFx(method, path string, cb func(fiber.Router), middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx {
	return func(appName, prefix string) fx.Option {
		// Create a wrapper function that adapts the handler to a fiber.Handler
		handlerWrapper := func() fiber.Handler {
			// Convert the handler to a fiber.Handler
			if h, ok := handler.(fiber.Handler); ok {
				return h
			}
			// If it's already a fiber.Handler function, return it directly
			if h, ok := handler.(func(*fiber.Ctx) error); ok {
				return h
			}
			// If it's a function that returns a fiber.Handler, call it to get the handler
			if h, ok := handler.(func() fiber.Handler); ok {
				return h()
			}
			// If it's a function that returns a fiber.Handler function, call it to get the handler
			if h, ok := handler.(func() func(*fiber.Ctx) error); ok {
				return h()
			}
			// Otherwise, panic as we can't handle this type
			panic(fmt.Sprintf("handler must be a fiber.Handler or func(*fiber.Ctx) error, got %T", handler))
		}

		// Call handlerWrapper to get the actual handler
		handler := handlerWrapper()

		// Create options for registering the handler
		handlerOption := fx.Provide(
			fx.Annotate(
				func() fiber.Handler {
					return handler
				},
				fx.ResultTags(fiberHandlers(appName, method, prefix, path)),
			),
		)

		// Create options for registering each middleware
		middlewareOptions := make([]fx.Option, 0, len(middlewareFuncs))
		for i, middlewareFunc := range middlewareFuncs {
			middlewareTag := fiberMiddlewareTag(appName, method+"-"+prefix+"-"+path+"-"+fmt.Sprintf("%d", i))
			middlewareOptions = append(middlewareOptions, fx.Provide(
				fx.Annotate(
					middlewareFunc,
					fx.ResultTags(middlewareTag),
				),
			))
		}

		// Create the route with placeholders for middleware handlers
		routeOption := fx.Provide(
			fx.Annotate(
				func(handler fiber.Handler, middlewares ...fiber.Handler) route {
					return route{
						Prefix:      prefix,
						Method:      method,
						Path:        path,
						Handler:     handler,
						CallBack:    cb,
						Middlewares: middlewares,
					}
				},
				fx.ParamTags(
					append(
						[]string{fiberHandlers(appName, method, prefix, path)},
						generateMiddlewareTags(appName, method, prefix, path, len(middlewareFuncs))...,
					)...,
				),
				fx.ResultTags(fiberHandlerRoutes(appName)),
			),
		)

		return fx.Options(append(append([]fx.Option{handlerOption}, middlewareOptions...), routeOption)...)
	}
}

// generateMiddlewareTags generates tags for middleware functions
func generateMiddlewareTags(appName, method, prefix, path string, count int) []string {
	tags := make([]string, count)
	for i := 0; i < count; i++ {
		tags[i] = fiberMiddlewareTag(appName, method+"-"+prefix+"-"+path+"-"+fmt.Sprintf("%d", i))
	}
	return tags
}
