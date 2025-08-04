#!/bin/bash

# Exit on any error
set -e

# Build the example
go build -o middleware_with_deps_example main.go

# Start the server in the background
./middleware_with_deps_example &
SERVER_PID=$!

# Wait for the server to start
sleep 2

# Initialize test status
TEST_STATUS=0

# Test the /public endpoint
echo "Testing /public endpoint:"
EXPECTED_PUBLIC="This is public content"
ACTUAL_PUBLIC=$(curl -s http://localhost:3003/public)
echo "Expected: $EXPECTED_PUBLIC"
echo "Actual: $ACTUAL_PUBLIC"
if [ "$ACTUAL_PUBLIC" = "$EXPECTED_PUBLIC" ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
    TEST_STATUS=1
fi
echo ""

# Test the /private/data endpoint without authorization
echo "Testing /private/data endpoint without authorization (should fail):"
EXPECTED_UNAUTH="Unauthorized"
ACTUAL_UNAUTH=$(curl -s http://localhost:3003/private/data)
echo "Expected: $EXPECTED_UNAUTH"
echo "Actual: $ACTUAL_UNAUTH"
if [ "$ACTUAL_UNAUTH" = "$EXPECTED_UNAUTH" ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
    TEST_STATUS=1
fi
echo ""

# Test the /private/data endpoint with authorization
echo "Testing /private/data endpoint with authorization:"
EXPECTED_PRIVATE="This is private content"
ACTUAL_PRIVATE=$(curl -s -H "Authorization: valid-token" http://localhost:3003/private/data)
echo "Expected: $EXPECTED_PRIVATE"
echo "Actual: $ACTUAL_PRIVATE"
if [ "$ACTUAL_PRIVATE" = "$EXPECTED_PRIVATE" ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
    TEST_STATUS=1
fi
echo ""

# Kill the server
kill $SERVER_PID

# Clean up
rm middleware_with_deps_example

# Report final status
if [ $TEST_STATUS -eq 0 ]; then
    echo "🎉 All tests passed successfully!"
    exit 0
else
    echo "❌ Some tests failed!"
    exit 1
fi
