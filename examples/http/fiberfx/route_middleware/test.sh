#!/bin/bash

# Exit on any error
set -e

# Build the example
go build -o route_middleware_example main.go

# Start the server in the background
./route_middleware_example &
SERVER_PID=$!

# Wait for the server to start
sleep 2

# Initialize test status
TEST_STATUS=0

# Test the /public endpoint
echo "Testing /public endpoint:"
EXPECTED_PUBLIC="This is public content"
ACTUAL_PUBLIC=$(curl -s http://localhost:3004/public)
echo "Expected: $EXPECTED_PUBLIC"
echo "Actual: $ACTUAL_PUBLIC"
if [ "$ACTUAL_PUBLIC" = "$EXPECTED_PUBLIC" ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
    TEST_STATUS=1
fi
echo ""

# Test the /api/data endpoint
echo "Testing /api/data endpoint:"
EXPECTED_API="This is public content"
ACTUAL_API=$(curl -s http://localhost:3004/api/data)
echo "Expected: $EXPECTED_API"
echo "Actual: $ACTUAL_API"
if [ "$ACTUAL_API" = "$EXPECTED_API" ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
    TEST_STATUS=1
fi
echo ""

# Test the /private endpoint without authorization
echo "Testing /private endpoint without authorization (should fail):"
EXPECTED_PRIVATE_UNAUTH="Unauthorized"
ACTUAL_PRIVATE_UNAUTH=$(curl -s http://localhost:3004/private)
echo "Expected: $EXPECTED_PRIVATE_UNAUTH"
echo "Actual: $ACTUAL_PRIVATE_UNAUTH"
if [ "$ACTUAL_PRIVATE_UNAUTH" = "$EXPECTED_PRIVATE_UNAUTH" ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
    TEST_STATUS=1
fi
echo ""

# Test the /private endpoint with authorization
echo "Testing /private endpoint with authorization:"
EXPECTED_PRIVATE_AUTH="This is private content"
ACTUAL_PRIVATE_AUTH=$(curl -s -H "Authorization: valid-token" http://localhost:3004/private)
echo "Expected: $EXPECTED_PRIVATE_AUTH"
echo "Actual: $ACTUAL_PRIVATE_AUTH"
if [ "$ACTUAL_PRIVATE_AUTH" = "$EXPECTED_PRIVATE_AUTH" ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
    TEST_STATUS=1
fi
echo ""

# Test the /admin endpoint without authorization
echo "Testing /admin endpoint without authorization (should fail):"
EXPECTED_ADMIN_UNAUTH="Unauthorized"
ACTUAL_ADMIN_UNAUTH=$(curl -s http://localhost:3004/admin)
echo "Expected: $EXPECTED_ADMIN_UNAUTH"
echo "Actual: $ACTUAL_ADMIN_UNAUTH"
if [ "$ACTUAL_ADMIN_UNAUTH" = "$EXPECTED_ADMIN_UNAUTH" ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
    TEST_STATUS=1
fi
echo ""

# Test the /admin endpoint with authorization
echo "Testing /admin endpoint with authorization:"
EXPECTED_ADMIN_AUTH="This is admin content"
ACTUAL_ADMIN_AUTH=$(curl -s -H "Authorization: valid-token" http://localhost:3004/admin)
echo "Expected: $EXPECTED_ADMIN_AUTH"
echo "Actual: $ACTUAL_ADMIN_AUTH"
if [ "$ACTUAL_ADMIN_AUTH" = "$EXPECTED_ADMIN_AUTH" ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
    TEST_STATUS=1
fi
echo ""

# Test the /metrics/requests endpoint
echo "Testing /metrics/requests endpoint:"
EXPECTED_METRICS="Metrics data"
ACTUAL_METRICS=$(curl -s http://localhost:3004/metrics/requests)
echo "Expected: $EXPECTED_METRICS"
echo "Actual: $ACTUAL_METRICS"
if [ "$ACTUAL_METRICS" = "$EXPECTED_METRICS" ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
    TEST_STATUS=1
fi
echo ""

# Kill the server
kill $SERVER_PID

# Clean up
rm route_middleware_example

# Report final status
if [ $TEST_STATUS -eq 0 ]; then
    echo "🎉 All tests passed successfully!"
    exit 0
else
    echo "❌ Some tests failed!"
    exit 1
fi
