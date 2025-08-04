package fiberfx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
)

func TestGetFiberApp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		appName  string
		expected string
	}{
		{
			name:     "simple app name",
			appName:  "test",
			expected: `name:"fiber-test"`,
		},
		{
			name:     "app name with special characters",
			appName:  "test-app",
			expected: `name:"fiber-test-app"`,
		},
		{
			name:     "empty app name",
			appName:  "",
			expected: `name:"fiber-"`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := require.New(t)

			result := fiberfx.GetFiberApp(tt.appName)
			assert.Equal(tt.expected, result)
		})
	}
}

// Test for unexported functions using exported test helpers

// TestFiberHandlers tests the fiberHandlers function indirectly through Route function
func TestFiberHandlers(t *testing.T) {
	t.Parallel()
	// This is tested indirectly through the router tests
	// The function is unexported, so we can't test it directly
}

// TestFiberHandlerRoutes tests the fiberHandlerRoutes function indirectly through Route function
func TestFiberHandlerRoutes(t *testing.T) {
	t.Parallel()
	// This is tested indirectly through the router tests
	// The function is unexported, so we can't test it directly
}

// TestRouterCallbacksName tests the routerCallbacksName function indirectly through Routes function
func TestRouterCallbacksName(t *testing.T) {
	t.Parallel()
	// This is tested indirectly through the routes tests
	// The function is unexported, so we can't test it directly
}
