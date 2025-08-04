# FiberFX

FiberFX is a package that integrates the [Fiber web framework](https://github.com/gofiber/fiber) with [Uber FX](https://github.com/uber-go/fx) for dependency injection.

## Features

- Create Fiber applications with dependency injection
- Register routes with DI
- Configure router callbacks
- Inject middleware using DI (new feature)

## Usage

### Basic Usage

```go
package main

import (
    "time"

    "github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
    "github.com/gofiber/fiber/v2"
    "go.uber.org/fx"
)

// Define a route handler
func HelloHandler(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
}

func main() {
    app := fx.New(
        // Create the Fiber app
        fiberfx.App(
            "example",
            fiberfx.Routes(
                []fiberfx.RouteFx{
                    fiberfx.Get("/hello", HelloHandler),
                },
            ),
        ),
        
        // Run the app
        fiberfx.RunApp(":3000", "example", 5*time.Second),
    )

    app.Run()
}
```

### Using Middleware with DI (New Feature)

You can now inject middleware using dependency injection. This allows you to create middleware that depends on other services.

```go
package main

import (
    "time"

    "github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
    "github.com/gofiber/fiber/v2"
    "go.uber.org/fx"
)

// Define a service that the middleware depends on
type LogService struct {
    Prefix string
}

// Constructor for the service
func NewLogService() *LogService {
    return &LogService{
        Prefix: "[LOG]",
    }
}

// Define a middleware that depends on the service
func NewLogMiddleware(logService *LogService) fiberfx.Middleware {
    return func(c *fiber.Ctx) error {
        fmt.Printf("%s Request: %s\n", logService.Prefix, c.Path())
        return c.Next()
    }
}

// Define a route handler
func HelloHandler(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
}

func main() {
    app := fx.New(
        // Provide the service
        fx.Provide(NewLogService),
        
        // Create the Fiber app with middleware support
        fiberfx.App(
            "example",
            fiberfx.Routes(
                []fiberfx.RouteFx{
                    fiberfx.Get("/hello", HelloHandler),
                },
            ),
            // Enable middleware injection
            fiberfx.WithMiddlewares(),
        ),
        
        // Register the middleware
        fiberfx.RegisterMiddleware("example", NewLogMiddleware),
        
        // Run the app
        fiberfx.RunApp(":3000", "example", 5*time.Second),
    )

    app.Run()
}
```

### Using Middleware with a Specific Prefix

You can also register middleware for specific route prefixes:

```go
package main

import (
    "time"

    "github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
    "github.com/gofiber/fiber/v2"
    "go.uber.org/fx"
)

// Define an auth service
type AuthService struct {
    // ...
}

// Constructor for the auth service
func NewAuthService() *AuthService {
    return &AuthService{}
}

// Define a middleware that depends on the auth service
func NewAuthMiddleware(authService *AuthService) fiberfx.Middleware {
    return func(c *fiber.Ctx) error {
        // Check for auth header
        if c.Get("Authorization") == "" {
            return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
        }
        return c.Next()
    }
}

// Define route handlers
func PublicHandler(c *fiber.Ctx) error {
    return c.SendString("Public content")
}

func PrivateHandler(c *fiber.Ctx) error {
    return c.SendString("Private content")
}

func main() {
    app := fx.New(
        // Provide the service
        fx.Provide(NewAuthService),
        
        // Create the Fiber app with middleware support
        fiberfx.App(
            "example",
            fiberfx.Routes(
                []fiberfx.RouteFx{
                    fiberfx.Get("/public", PublicHandler),
                    fiberfx.Get("/private/data", PrivateHandler),
                },
            ),
            // Enable middleware injection
            fiberfx.WithMiddlewares(),
        ),
        
        // Register the middleware with a specific prefix
        fiberfx.RegisterMiddlewareWithPrefix("example", "/private", NewAuthMiddleware),
        
        // Run the app
        fiberfx.RunApp(":3000", "example", 5*time.Second),
    )

    app.Run()
}
```

### Using Per-Route Middleware (New Feature)

You can now apply middleware to specific routes using the new `*WithMiddleware` functions:

```go
package main

import (
    "time"

    "github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
    "github.com/gofiber/fiber/v2"
    "go.uber.org/fx"
)

// Define middleware functions
func LogMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        fmt.Println("Request received:", c.Path())
        return c.Next()
    }
}

func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        if c.Get("Authorization") == "" {
            return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
        }
        return c.Next()
    }
}

// Define route handlers
func PublicHandler(c *fiber.Ctx) error {
    return c.SendString("Public content")
}

func ProtectedHandler(c *fiber.Ctx) error {
    return c.SendString("Protected content")
}

func AdminHandler(c *fiber.Ctx) error {
    return c.SendString("Admin content")
}

func main() {
    app := fx.New(
        // Create the Fiber app
        fiberfx.App(
            "example",
            fiberfx.Routes(
                []fiberfx.RouteFx{
                    // Route with no middleware
                    fiberfx.Get("/public", PublicHandler),
                    
                    // Route with a single middleware
                    fiberfx.GetWithMiddleware("/protected", 
                        []fiber.Handler{AuthMiddleware()}, 
                        ProtectedHandler),
                    
                    // Route with multiple middleware
                    fiberfx.GetWithMiddleware("/admin", 
                        []fiber.Handler{LogMiddleware(), AuthMiddleware()}, 
                        AdminHandler),
                },
            ),
        ),
        
        // Run the app
        fiberfx.RunApp(":3000", "example", 5*time.Second),
    )

    app.Run()
}
```

You can also combine router callbacks with middleware:

```go
// Route with both router callback and middleware
fiberfx.GetWithRouterCallbackAndMiddleware("/api/users/:id",
    func(router fiber.Router) {
        // Configure the route
        router.Use(limiter.New())
    },
    []fiber.Handler{AuthMiddleware()},
    GetUserHandler)
```

### Using Per-Route Middleware with Dependencies (New Feature)

You can now apply middleware with dependencies to specific routes using the new `*WithMiddlewareFx` functions:

```go
package main

import (
    "time"

    "github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
    "github.com/gofiber/fiber/v2"
    "go.uber.org/fx"
)

// Define a service that the middleware depends on
type AuthService struct {
    // Authentication logic
}

// Constructor for the service
func NewAuthService() *AuthService {
    return &AuthService{}
}

// Define middleware function with dependencies
func AuthMiddlewareWithDeps(authService *AuthService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Use authService to validate the request
        if c.Get("Authorization") == "" {
            return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
        }
        return c.Next()
    }
}

// Define a logging service
type LogService struct {
    Prefix string
}

// Constructor for the logging service
func NewLogService() *LogService {
    return &LogService{
        Prefix: "[LOG]",
    }
}

// Define another middleware function with dependencies
func LogMiddlewareWithDeps(logService *LogService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        fmt.Printf("%s Request: %s\n", logService.Prefix, c.Path())
        return c.Next()
    }
}

// Define route handlers
func PublicHandler(c *fiber.Ctx) error {
    return c.SendString("Public content")
}

func ProtectedHandler(c *fiber.Ctx) error {
    return c.SendString("Protected content")
}

func main() {
    app := fx.New(
        // Provide the services
        fx.Provide(
            NewAuthService,
            NewLogService,
        ),
        
        // Create the Fiber app
        fiberfx.App(
            "example",
            fiberfx.Routes(
                []fiberfx.RouteFx{
                    // Route with no middleware
                    fiberfx.Get("/public", PublicHandler),
                    
                    // Route with a single middleware with dependencies
                    fiberfx.GetWithMiddlewareFx("/protected", 
                        []fiberfx.RouteMiddlewareFunc{AuthMiddlewareWithDeps}, 
                        ProtectedHandler),
                    
                    // Route with multiple middleware with dependencies
                    fiberfx.GetWithMiddlewareFx("/admin", 
                        []fiberfx.RouteMiddlewareFunc{LogMiddlewareWithDeps, AuthMiddlewareWithDeps}, 
                        ProtectedHandler),
                    
                    // Route with router callback and middleware with dependencies
                    fiberfx.GetWithRouterCallbackAndMiddlewareFx("/api/users/:id",
                        func(router fiber.Router) {
                            // Configure the route
                            router.Use(func(c *fiber.Ctx) error {
                                // Additional middleware
                                return c.Next()
                            })
                        },
                        []fiberfx.RouteMiddlewareFunc{AuthMiddlewareWithDeps},
                        ProtectedHandler),
                },
            ),
        ),
        
        // Run the app
        fiberfx.RunApp(":3000", "example", 5*time.Second),
    )

    app.Run()
}
```

### Backward Compatibility

The middleware injection feature is opt-in, so existing code will continue to work without changes. If you want to use the traditional approach to adding middleware, you can use the `WithAfterCreate` option:

```go
package main

import (
    "time"

    "github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
    "github.com/gofiber/fiber/v2"
    "go.uber.org/fx"
)

// Define a route handler
func HelloHandler(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
}

func main() {
    app := fx.New(
        // Create the Fiber app without middleware support
        fiberfx.App(
            "example",
            fiberfx.Routes(
                []fiberfx.RouteFx{
                    fiberfx.Get("/hello", HelloHandler),
                },
            ),
            // Use traditional middleware approach
            fiberfx.WithAfterCreate(func(app *fiber.App) {
                app.Use(func(c *fiber.Ctx) error {
                    fmt.Println("Request received:", c.Path())
                    return c.Next()
                })
            }),
        ),
        
        // Run the app
        fiberfx.RunApp(":3000", "example", 5*time.Second),
    )

    app.Run()
}
```

## API Reference

### Middleware Functions

- `RegisterMiddleware(appName string, middleware any) fx.Option`: Registers a middleware to be used with a specific app.
- `RegisterMiddlewareWithPrefix(appName, prefix string, middleware any) fx.Option`: Registers a middleware to be used with a specific app and route prefix.
- `WithMiddlewares() Option`: Enables middleware injection for the app.

### Route Functions with Middleware

- `GetWithMiddleware(path string, middlewares []fiber.Handler, handler any) RouteFx`: Creates a GET route with specific middleware.
- `PostWithMiddleware(path string, middlewares []fiber.Handler, handler any) RouteFx`: Creates a POST route with specific middleware.
- `PutWithMiddleware(path string, middlewares []fiber.Handler, handler any) RouteFx`: Creates a PUT route with specific middleware.
- `PatchWithMiddleware(path string, middlewares []fiber.Handler, handler any) RouteFx`: Creates a PATCH route with specific middleware.
- `DeleteWithMiddleware(path string, middlewares []fiber.Handler, handler any) RouteFx`: Creates a DELETE route with specific middleware.

### Route Functions with Router Callback and Middleware

- `GetWithRouterCallbackAndMiddleware(path string, cb func(fiber.Router), middlewares []fiber.Handler, handler any) RouteFx`: Creates a GET route with router callback and specific middleware.
- `PostWithRouterCallbackAndMiddleware(path string, cb func(fiber.Router), middlewares []fiber.Handler, handler any) RouteFx`: Creates a POST route with router callback and specific middleware.
- `PutWithRouterCallbackAndMiddleware(path string, cb func(fiber.Router), middlewares []fiber.Handler, handler any) RouteFx`: Creates a PUT route with router callback and specific middleware.
- `PatchWithRouterCallbackAndMiddleware(path string, cb func(fiber.Router), middlewares []fiber.Handler, handler any) RouteFx`: Creates a PATCH route with router callback and specific middleware.
- `DeleteWithRouterCallbackAndMiddleware(path string, cb func(fiber.Router), middlewares []fiber.Handler, handler any) RouteFx`: Creates a DELETE route with router callback and specific middleware.

### Route Functions with Middleware with Dependencies

- `GetWithMiddlewareFx(path string, middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx`: Creates a GET route with middleware that can have dependencies.
- `PostWithMiddlewareFx(path string, middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx`: Creates a POST route with middleware that can have dependencies.
- `PutWithMiddlewareFx(path string, middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx`: Creates a PUT route with middleware that can have dependencies.
- `PatchWithMiddlewareFx(path string, middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx`: Creates a PATCH route with middleware that can have dependencies.
- `DeleteWithMiddlewareFx(path string, middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx`: Creates a DELETE route with middleware that can have dependencies.

### Route Functions with Router Callback and Middleware with Dependencies

- `GetWithRouterCallbackAndMiddlewareFx(path string, cb func(fiber.Router), middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx`: Creates a GET route with router callback and middleware that can have dependencies.
- `PostWithRouterCallbackAndMiddlewareFx(path string, cb func(fiber.Router), middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx`: Creates a POST route with router callback and middleware that can have dependencies.
- `PutWithRouterCallbackAndMiddlewareFx(path string, cb func(fiber.Router), middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx`: Creates a PUT route with router callback and middleware that can have dependencies.
- `PatchWithRouterCallbackAndMiddlewareFx(path string, cb func(fiber.Router), middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx`: Creates a PATCH route with router callback and middleware that can have dependencies.
- `DeleteWithRouterCallbackAndMiddlewareFx(path string, cb func(fiber.Router), middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx`: Creates a DELETE route with router callback and middleware that can have dependencies.

### Low-level Route Functions

- `RouteWithMiddleware(method, path string, cb func(fiber.Router), middlewares []fiber.Handler, handler any) RouteFx`: Creates a route with specific method, path, router callback, middleware, and handler.
- `RouteWithMiddlewareFx(method, path string, cb func(fiber.Router), middlewareFuncs []RouteMiddlewareFunc, handler any) RouteFx`: Creates a route with specific method, path, router callback, middleware with dependencies, and handler.

### Types

- `Middleware`: Represents a Fiber middleware function.
- `RouteMiddlewareFunc`: Represents a function that returns a Fiber middleware and can have dependencies injected by uberfx.
