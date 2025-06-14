package loggerfx

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/nano-interactive/go-utils/v2"
	appLogger "github.com/nano-interactive/go-utils/v2/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"go.uber.org/fx"
)

type (
	SinkType int
	Sink     struct {
		Level       string
		Args        []any
		Type        SinkType
		PrettyPrint bool
	}
)

const (
	Stdout SinkType = iota << 1
	Stderr
	File
	BufferedIO
)

var (
	ErrArgForFileNotProvided    = errors.New("output file path must be provided for SINK FILE")
	ErrUnexpectedArgForFileType = errors.New("invalid type for the sink FILE => expected string or fmt.Stringer")
	ErrInvalidSinkType          = errors.New("invalid sink type")
)

func ZerologModule(sink Sink) fx.Option {
	return fx.Module("ZerologLogger", fx.Provide(
		func(lc fx.Lifecycle) (zerolog.Logger, error) {
			w, closer, err := getZerologWriter(sink)
			if err != nil {
				return zerolog.Logger{}, err
			}

			if closer != nil {
				lc.Append(fx.StopHook(closer))
			}

			return appLogger.New(sink.Level, false).
				Output(w).
				With().
				Stack().
				Logger(), nil
		}),
		fx.Invoke(func(lc fx.Lifecycle) error {
			w, closer, err := getZerologWriter(sink)
			if err != nil {
				return err
			}

			if closer != nil {
				lc.Append(fx.StopHook(closer))
			}

			//nolint:all
			zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
			appLogger.ConfigureDefaultLogger(sink.Level, false, w)

			return nil
		}),
	)
}

func getZerologWriter(sink Sink) (io.Writer, func() error, error) {
	switch sink.Type {
	case Stdout:
		buffer := bufio.NewWriterSize(os.Stdout, 32*1024)

		if sink.PrettyPrint {
			options := []func(w *zerolog.ConsoleWriter){
				func(w *zerolog.ConsoleWriter) {
					w.Out = bufio.NewWriterSize(os.Stdin, 32*1024)
				},
			}
			if len(sink.Args) > 0 {
				for _, arg := range sink.Args {
					switch vals := arg.(type) {
					case []func(w *zerolog.ConsoleWriter):
						options = append(options, vals...)
					case func(w *zerolog.ConsoleWriter):
						options = append(options, vals)
					default:
						return nil, nil, fmt.Errorf("invalid type for sink.Args: %T", vals)
					}
				}
			}

			return zerolog.NewConsoleWriter(options...), nil, nil
		}

		return buffer, buffer.Flush, nil
	case Stderr:
		buffer := bufio.NewWriterSize(os.Stderr, 32*1024)

		if sink.PrettyPrint {
			options := []func(w *zerolog.ConsoleWriter){
				func(w *zerolog.ConsoleWriter) {
					w.Out = buffer
				},
			}

			if len(sink.Args) > 0 {
				for _, arg := range sink.Args {
					switch vals := arg.(type) {
					case []func(w *zerolog.ConsoleWriter):
						options = append(options, vals...)
					case func(w *zerolog.ConsoleWriter):
						options = append(options, vals)
					default:
						return nil, nil, fmt.Errorf("invalid type for sink.Args: %T", vals)
					}
				}
			}

			return zerolog.NewConsoleWriter(options...), nil, nil
		}

		return buffer, buffer.Flush, nil
	case BufferedIO:
		buffer, err := newNonBlockingBufferedWriter(sink)
		if err != nil {
			return nil, nil, err
		}

		return buffer, buffer.Close, nil
	case File:
		f, err := extractFile(sink)
		if err != nil {
			return nil, nil, err
		}

		return f, f.Close, nil
	default:
		return nil, nil, ErrInvalidSinkType
	}
}

func extractFile(sink Sink) (*os.File, error) {
	if len(sink.Args) == 0 {
		return nil, ErrArgForFileNotProvided
	}

	var path string

	switch val := sink.Args[0].(type) {
	case string:
		path = val
	case fmt.Stringer:
		path = val.String()
	default:
		return nil, ErrUnexpectedArgForFileType
	}

	f, err := utils.CreateLogFile(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}
