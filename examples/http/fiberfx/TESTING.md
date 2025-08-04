# Testing FiberFX Examples

This document explains how to test the FiberFX examples to ensure they're working correctly.

## Changes Made

The following changes were made to fix issues in the FiberFX library:

1. Fixed the `RouteWithMiddleware` function in `router.go` to correctly handle handler functions by creating a wrapper function that adapts the handler to a `fiber.Handler`.
2. Fixed the `RouteWithMiddlewareFx` function in `router.go` with the same approach.
3. Modified each example to use a different port to avoid conflicts:
   - Basic example: Port 3001
   - Middleware example: Port 3002
   - Middleware with Dependencies example: Port 3003
   - Route-Specific Middleware example: Port 3004

## Testing the Examples

Each example has a test script that can be run to verify it's working correctly. The test scripts will:

1. Build the example
2. Start the server in the background
3. Test the endpoints with curl and validate the output
4. Kill the server
5. Clean up
6. Report test results with proper exit codes

### Output Validation

The test scripts now validate the output from the HTTP server by:
1. Defining expected output for each endpoint
2. Capturing actual output from curl
3. Comparing expected and actual output
4. Reporting success or failure for each test
5. Providing an overall test status at the end

If all tests pass, the script will exit with code 0 and display:
```
🎉 All tests passed successfully!
```

If any test fails, the script will exit with code 1 and display:
```
❌ Some tests failed!
```

This makes the test scripts suitable for use in CI/CD pipelines, as they will fail the build if any test fails.

### Running the Tests

To run a test script, navigate to the example directory and run:

```bash
chmod +x test.sh
./test.sh
```

### Basic Example

The basic example demonstrates the basic usage of FiberFX with simple routes:
- `/hello` - Returns a simple "Hello, World!" message
- `/users/:id` - Returns a message with the user ID from the URL parameter

To test manually:
```bash
cd basic
go run main.go
```

In another terminal:
```bash
curl http://localhost:3001/hello
curl http://localhost:3001/users/123
```

### Middleware Example

The middleware example demonstrates how to use global middleware with FiberFX:
- `/public` - Accessible without authentication
- `/private/data` - Requires authentication

To test manually:
```bash
cd middleware
go run main.go
```

In another terminal:
```bash
curl http://localhost:3002/public
curl http://localhost:3002/private/data  # Should fail
curl -H "Authorization: valid-token" http://localhost:3002/private/data
```

### Middleware with Dependencies Example

The middleware with dependencies example demonstrates how to use middleware with dependencies:
- `/public` - Accessible without authentication
- `/private/data` - Requires authentication with a valid token

To test manually:
```bash
cd middleware_with_deps
go run main.go
```

In another terminal:
```bash
curl http://localhost:3003/public
curl http://localhost:3003/private/data  # Should fail
curl -H "Authorization: valid-token" http://localhost:3003/private/data
```

### Route-Specific Middleware Example

The route-specific middleware example demonstrates how to use route-specific middleware:
- `/public` - No middleware
- `/api/data` - With middleware (no dependencies)
- `/private` - With middleware with dependencies
- `/admin` - With multiple middleware with dependencies
- `/metrics/requests` - With router callback and middleware with dependencies

To test manually:
```bash
cd route_middleware
go run main.go
```

In another terminal:
```bash
curl http://localhost:3004/public
curl http://localhost:3004/api/data
curl http://localhost:3004/private  # Should fail
curl -H "Authorization: valid-token" http://localhost:3004/private
curl http://localhost:3004/admin  # Should fail
curl -H "Authorization: valid-token" http://localhost:3004/admin
curl http://localhost:3004/metrics/requests
```
