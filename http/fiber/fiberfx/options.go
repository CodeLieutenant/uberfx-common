package fiberfx

import "github.com/gofiber/fiber/v2"

type (
	appOptions struct {
		cfg         fiber.Config
		afterCreate func(app *fiber.App)
	}

	Option func(opts *appOptions)
)

func WithFiberConfig(cfg fiber.Config) Option {
	return func(opts *appOptions) {
		opts.cfg = cfg
	}
}

func WithAfterCreate(afterCreate func(app *fiber.App)) Option {
	return func(opts *appOptions) {
		opts.afterCreate = afterCreate
	}
}
