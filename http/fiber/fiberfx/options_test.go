package fiberfx_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"github.com/CodeLieutenant/uberfx-common/v3/http/fiber/fiberfx"
)

func TestWithFiberConfig(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Create a custom Fiber config
	customConfig := fiber.Config{
		Prefork:               true,
		ServerHeader:          "TestServer",
		StrictRouting:         true,
		CaseSensitive:         true,
		DisableStartupMessage: true,
	}

	// Apply the WithFiberConfig option
	opts := fiberfx.ApplyOption(fiberfx.WithFiberConfig(customConfig))

	// Verify that the config was set correctly
	config := fiberfx.GetConfig(opts)
	assert.Equal(customConfig.Prefork, config.Prefork)
	assert.Equal(customConfig.ServerHeader, config.ServerHeader)
	assert.Equal(customConfig.StrictRouting, config.StrictRouting)
	assert.Equal(customConfig.CaseSensitive, config.CaseSensitive)
	assert.Equal(customConfig.DisableStartupMessage, config.DisableStartupMessage)
}

func TestWithAfterCreate(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Create a flag to track if the afterCreate function was called
	var afterCreateCalled bool

	// Create an afterCreate function
	afterCreate := func(app *fiber.App) {
		afterCreateCalled = true
	}

	// Apply the WithAfterCreate option
	opts := fiberfx.ApplyOption(fiberfx.WithAfterCreate(afterCreate))

	// Get the afterCreate function
	resultAfterCreate := fiberfx.GetAfterCreate(opts)

	// Verify that the afterCreate function was set
	assert.NotNil(resultAfterCreate)

	// Call the afterCreate function
	resultAfterCreate(fiber.New())

	// Verify that our original afterCreate function was called
	assert.True(afterCreateCalled)
}

func TestWithMiddlewares(t *testing.T) {
	t.Parallel()

	// Test with no existing afterCreate function
	t.Run("with no existing afterCreate", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		// Apply the WithMiddlewares option
		opts := fiberfx.ApplyOption(fiberfx.WithMiddlewares())

		// Verify that useMiddlewares is set to true
		assert.True(fiberfx.GetUseMiddlewares(opts))

		// Verify that an afterCreate function was created
		assert.NotNil(fiberfx.GetAfterCreate(opts))

		// Call the afterCreate function (should not panic)
		afterCreate := fiberfx.GetAfterCreate(opts)
		afterCreate(fiber.New())
	})

	// Test with an existing afterCreate function
	t.Run("with existing afterCreate", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)

		// Create a flag to track if the original afterCreate function was called
		var originalAfterCreateCalled bool

		// Create an original afterCreate function
		originalAfterCreate := func(app *fiber.App) {
			originalAfterCreateCalled = true
		}

		// Apply both options in sequence
		opts := fiberfx.ApplyOptions(
			fiberfx.WithAfterCreate(originalAfterCreate),
			fiberfx.WithMiddlewares(),
		)

		// Verify that useMiddlewares is set to true
		assert.True(fiberfx.GetUseMiddlewares(opts))

		// Verify that an afterCreate function exists
		assert.NotNil(fiberfx.GetAfterCreate(opts))

		// Call the afterCreate function
		afterCreate := fiberfx.GetAfterCreate(opts)
		afterCreate(fiber.New())

		// Verify that the original afterCreate function was called
		assert.True(originalAfterCreateCalled)
	})
}
