package fiber

import (
	"errors"

	gofiber "github.com/gofiber/fiber/v2"
	"github.com/invopop/validation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrDefaultHandler = errors.New("default handler")

type ErrorResponse struct {
	Message any `json:"message,omitempty"`
}

func ErrorHandler() gofiber.ErrorHandler {
	return ErrorHandlerWithCustomHandler(log.Logger, nil)
}

func ErrorHandlerWithCustomHandler(logger zerolog.Logger, handler gofiber.ErrorHandler) gofiber.ErrorHandler {
	return func(c *gofiber.Ctx, err error) error {
		if handler != nil {
			if errH := handler(c, err); !errors.Is(errH, ErrDefaultHandler) {
				return errH
			}
		}

		c.Set(gofiber.HeaderContentType, gofiber.MIMEApplicationJSONCharsetUTF8)

		if errors.Is(err, primitive.ErrInvalidHex) {
			return c.Status(gofiber.StatusBadRequest).JSON(ErrorResponse{
				Message: "Invalid JSON Payload, check your input",
			})
		}

		{
			var fiberErr *gofiber.Error
			if errors.As(err, &fiberErr) {
				return c.Status(fiberErr.Code).JSON(ErrorResponse{
					Message: fiberErr.Message,
				})
			}
		}

		{
			var validationErr validation.Errors
			if errors.As(err, &validationErr) {
				return c.Status(gofiber.StatusUnprocessableEntity).JSON(validationErr)
			}
		}

		logger.Error().Err(err).
			Str("path", c.Route().Path).
			Msg("Failed to process request")

		return c.Status(gofiber.StatusInternalServerError).JSON(ErrorResponse{
			Message: "An error has occurred!",
		})
	}
}
