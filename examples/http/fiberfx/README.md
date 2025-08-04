# FiberFX Examples

This directory contains examples demonstrating how to use the FiberFX package, which integrates the [Fiber web framework](https://github.com/gofiber/fiber) with [Uber FX](https://github.com/uber-go/fx) for dependency injection.

## Examples

### 1. Basic Example

**File**: [basic/main.go](basic/main.go)

This example demonstrates the basic usage of FiberFX:
- Creating a Fiber application with FX
- Registering simple routes
- Running the application

### 2. Middleware Example

**File**: [middleware/main.go](middleware/main.go)

This example demonstrates how to use global middleware with FiberFX:
- Creating middleware functions
- Registering global middleware
- Registering middleware for specific route prefixes

### 3. Middleware with Dependencies Example

**File**: [middleware_with_deps/main.go](middleware_with_deps/main.go)

This example demonstrates how to use middleware with dependencies:
- Creating services that middleware depends on
- Creating middleware functions that accept these services as dependencies
- Registering middleware with the application
- Using both global middleware and prefix-specific middleware

### 4. Route-Specific Middleware Example

**File**: [route_middleware/main.go](route_middleware/main.go)

This example demonstrates how to use route-specific middleware:
- Creating middleware without dependencies
- Creating middleware with dependencies
- Applying middleware to specific routes
- Combining router callbacks with middleware

## Running the Examples

To run any of the examples, navigate to the example directory and run:

```bash
go run main.go
```

For example:

```bash
cd basic
go run main.go
```

Each example will start a server on port 3000. You can test the endpoints using curl or a web browser.

## Testing the Examples

You can test the examples using curl. For example:

### Basic Example

```bash
curl http://localhost:3000/hello
```

### Middleware Example

```bash
curl http://localhost:3000/public
curl -H "Authorization: valid-token" http://localhost:3000/private/data
```

### Middleware with Dependencies Example

```bash
curl http://localhost:3000/public
curl -H "Authorization: valid-token" http://localhost:3000/private/data
```

### Route-Specific Middleware Example

```bash
curl http://localhost:3000/public
curl http://localhost:3000/api/data
curl -H "Authorization: valid-token" http://localhost:3000/private
curl -H "Authorization: valid-token" http://localhost:3000/admin
curl http://localhost:3000/metrics/requests
```
