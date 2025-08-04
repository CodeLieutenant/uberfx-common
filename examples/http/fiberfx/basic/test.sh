#!/bin/bash

# Exit on any error
set -e

# Build the example
go build -o basic_example main.go

# Start the server in the background
./basic_example &
SERVER_PID=$!

# Wait for the server to start
sleep 2

# Initialize test status
TEST_STATUS=0

# Test the /hello endpoint
echo "Testing /hello endpoint:"
EXPECTED_HELLO="Hello, World!"
ACTUAL_HELLO=$(curl -s http://localhost:3001/hello)
echo "Expected: $EXPECTED_HELLO"
echo "Actual: $ACTUAL_HELLO"
if [ "$ACTUAL_HELLO" = "$EXPECTED_HELLO" ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
    TEST_STATUS=1
fi
echo ""

# Test the /users/123 endpoint
echo "Testing /users/123 endpoint:"
EXPECTED_USER="User ID: 123"
ACTUAL_USER=$(curl -s http://localhost:3001/users/123)
echo "Expected: $EXPECTED_USER"
echo "Actual: $ACTUAL_USER"
if [ "$ACTUAL_USER" = "$EXPECTED_USER" ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
    TEST_STATUS=1
fi
echo ""

# Kill the server
kill $SERVER_PID

# Clean up
rm basic_example

# Report final status
if [ $TEST_STATUS -eq 0 ]; then
    echo "🎉 All tests passed successfully!"
    exit 0
else
    echo "❌ Some tests failed!"
    exit 1
fi
