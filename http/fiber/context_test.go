package fiber_test

import (
	"testing"

	"github.com/CodeLieutenant/uberfx-common/v3/constants"
	"github.com/CodeLieutenant/uberfx-common/v3/http/fiber"
	gofiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestContextMiddleware(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	app := gofiber.New()
	app.Use(fiber.Context())
	app.Get("/", func(ctx *gofiber.Ctx) error {
		return ctx.SendStatus(gofiber.StatusOK)
	})

	h := app.Handler()
	ctx := &fasthttp.RequestCtx{}
	h(ctx)

	assert.NotNil(ctx.UserValue(constants.CancelFuncContextKey))
}
