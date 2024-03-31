package fiber

import (
	"fmt"
	"net"

	"github.com/goccy/go-json"
	gofiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

func RunServer(ip string, port int, app *gofiber.App) {
	addr := fmt.Sprintf("%s:%d", ip, port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.
			Fatal().
			Err(err).
			Msg("Error while creating net.Listener for HTTP Server")
	}

	err = app.Listener(listener)

	if err != nil {
		log.
			Fatal().
			Err(err).
			Msg("Cannot start Fiber HTTP Server")
	}
}

var DefaultFiberConfig = gofiber.Config{
	StrictRouting:                true,
	EnablePrintRoutes:            false,
	Prefork:                      false,
	DisableStartupMessage:        true,
	DisableDefaultDate:           true,
	DisableHeaderNormalizing:     false,
	DisablePreParseMultipartForm: true,
	DisableKeepalive:             false,
	EnableIPValidation:           true,
	EnableSplittingOnParsers:     true,
	ErrorHandler:                 ErrorHandler(),
	JSONEncoder:                  json.Marshal,
	JSONDecoder:                  json.Unmarshal,
}

type AfterCreate func(*gofiber.App)

func CreateApplication(after AfterCreate, cfg ...gofiber.Config) *gofiber.App {
	c := DefaultFiberConfig

	if len(cfg) > 0 {
		c = cfg[0]
	}

	app := gofiber.New(c)

	app.Use(recover.New())
	app.Use(Context())

	if after != nil {
		after(app)
	}

	return app
}
