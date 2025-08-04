# UberFX Common Bindings

[![Testing](https://github.com/CodeLieutenant/uberfx-common/actions/workflows/test.yml/badge.svg)](https://github.com/CodeLieutenant/uberfx-common/actions/workflows/test.yml)
[![Linting](https://github.com/CodeLieutenant/uberfx-common/actions/workflows/lint.yml/badge.svg)](https://github.com/CodeLieutenant/uberfx-common/actions/workflows/lint.yml)
[![Security](https://github.com/CodeLieutenant/uberfx-common/actions/workflows/security.yml/badge.svg)](https://github.com/CodeLieutenant/uberfx-common/actions/workflows/security.yml)
[![codecov](https://codecov.io/gh/CodeLieutenant/uberfx-common/graph/badge.svg?token=opB123xE25)](https://codecov.io/gh/CodeLieutenant/uberfx-common)
[![Go Report Card](https://goreportcard.com/badge/github.com/CodeLieutenant/uberfx-common)](https://goreportcard.com/report/github.com/CodeLieutenant/uberfx-common)

A collection of common modules and utilities for building Go applications with [Uber FX](https://github.com/uber-go/fx) dependency injection framework. This library provides ready-to-use FX modules for configuration, logging, database connections, HTTP servers, and AMQP messaging.

## Table of Contents

- [Installation](#installation)
- [Modules](#modules)
  - [configfx](#configfx)
  - [loggerfx](#loggerfx)
  - [databasesfx](#databasesfx)
  - [http/fiber/fiberfx](#httpfiberfiberfx)
  - [amqpfx](#amqpfx)
- [Examples](#examples)
- [License](#license)
- [Contributing](#contributing)

## Installation

This library requires Go 1.24 or later.

```bash
go get github.com/CodeLieutenant/uberfx-common/v3
```

## Modules

### configfx

The `configfx` module provides configuration management using [Viper](https://github.com/spf13/viper), integrated with Uber FX.

Features:

- Type-safe configuration with generics
- Automatic loading from standard locations
- Easy integration with FX dependency injection

Example:

```go
type AppConfig struct {
    Port int `mapstructure:"port" yaml:"port" default:"8080"`
    LogLevel string `mapstructure:"log_level" yaml:"log_level" default:"info"`
}

func main() {
    app := fx.New(
        // Load configuration
        fx.Provide(func() (AppConfig, error) {
            return configfx.New[AppConfig]("myapp")
        }),

        // Create a module with the loaded configuration
        fx.Invoke(func(cfg AppConfig) {
            // Use configuration
        }),
    )

    app.Run()
}
```

### loggerfx

The `loggerfx` module provides logging functionality using [zerolog](https://github.com/rs/zerolog), integrated with Uber FX.

Features:

- Multiple output sinks (stdout, stderr, file, buffered I/O)
- Pretty printing option for development
- Proper lifecycle management

Example:

```go
func main() {
    app := fx.New(
        // Configure logging
        loggerfx.ZerologModule(loggerfx.Sink{
            Level:       "info",
            Type:        loggerfx.Stdout,
            PrettyPrint: true,
        }),

        // Use the logger
        fx.Invoke(func(log zerolog.Logger) {
            log.Info().Msg("Application started")
        }),
    )

    app.Run()
}
```

### databasesfx

The `databasesfx` module provides PostgreSQL database integration using [pgx](https://github.com/jackc/pgx), with support for migrations via [golang-migrate](https://github.com/golang-migrate/migrate).

Features:

- Connection pooling
- Configuration via struct
- Database migrations
- Proper lifecycle management

Example:

```go
func main() {
    app := fx.New(
        // Configure and provide PostgreSQL connection
        databasesfx.PostgresModule(databasesfx.PostgresConfig{
            ApplicationName:    "myapp",
            DBName:             "mydatabase",
            Host:               "localhost",
            Port:               5432,
            Username:           "postgres",
            Password:           "password",
            MaxOpenConnections: 10,
            // ... other configuration options
        }),

        // Use the database connection
        fx.Invoke(func(pool *pgxpool.Pool) {
            // Use the connection pool
        }),
    )

    app.Run()
}
```

### http/fiber/fiberfx

The `http/fiber/fiberfx` module provides integration with the [Fiber](https://github.com/gofiber/fiber) web framework.

Features:

- Easy route definition
- Middleware support
- Proper lifecycle management
- Type-safe handler registration

Example:

```go
func HelloHandler(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
}

func main() {
    app := fx.New(
        // Create the Fiber app
        fiberfx.App(
            "myapp",
            fiberfx.Routes(
                []fiberfx.RouteFx{
                    fiberfx.Get("/hello", HelloHandler),
                    fiberfx.Get("/users/:id", UserHandler),
                },
            ),
        ),

        // Run the app
        fiberfx.RunApp(":3000", "myapp", 5*time.Second),
    )

    app.Run()
}
```

### amqpfx

The `amqpfx` module provides integration with AMQP (RabbitMQ) for messaging, using [go-amqp](https://github.com/nano-interactive/go-amqp).

Features:

- Consumer and publisher support
- Type-safe message handling with generics
- Proper lifecycle management
- Multiple consumer types (function-based, interface-based, raw)

Example:

```go
type Message struct {
    Content string `json:"content"`
}

func handleMessage(ctx context.Context, msg Message) error {
    // Process the message
    return nil
}

func main() {
    app := fx.New(
        // Configure and provide AMQP consumer
        amqpfx.ConsumerModuleFunc(
            handleMessage,
            consumer.QueueDeclare{
                QueueName: "my-queue",
                Durable:   true,
            },
            connection.Config{
                ConnectionName: "my-connection",
                URI:            "amqp://guest:guest@localhost:5672/",
            },
        ),

        // Configure and provide AMQP publisher
        amqpfx.PublisherModule[Message](
            connection.Config{
                ConnectionName: "my-connection",
                URI:            "amqp://guest:guest@localhost:5672/",
            },
            "my-exchange",
        ),

        // Use the publisher
        fx.Invoke(func(pub *publisher.Publisher[Message]) {
            // Publish messages
        }),
    )

    app.Run()
}
```

## Examples

The repository includes several examples demonstrating how to use the various modules:

- [Basic Fiber Example](examples/http/fiberfx/basic/main.go): A simple HTTP server using Fiber
- [Middleware Example](examples/http/fiberfx/middleware/main.go): Using middleware with Fiber
- [Route Middleware Example](examples/http/fiberfx/route_middleware/main.go): Applying middleware to specific routes
- [Middleware with Dependencies Example](examples/http/fiberfx/middleware_with_deps/main.go): Using middleware with FX dependencies

## License

This project is licensed under the [Apache License 2.0](LICENSE).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
