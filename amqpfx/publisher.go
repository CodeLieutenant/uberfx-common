package amqpfx

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	"github.com/nano-interactive/go-amqp/v3/connection"
	"github.com/nano-interactive/go-amqp/v3/publisher"
)

func GetPublisherName(connectionName, exchangeName string) string {
	return fmt.Sprintf("amqp-publisher-%s-%s", exchangeName, connectionName)
}

func GetPublisherParamName(connectionName, exchangeName string) string {
	return fmt.Sprintf(`name:"amqp-publisher-param-%s-%s"`, exchangeName, connectionName)
}

func PublisherModule[T any](
	connectionOptions connection.Config,
	exchangeName string,
	options ...publisher.Option[T],
) fx.Option {
	module := fmt.Sprintf("amqp-publisher-module-%s-%s", exchangeName, connectionOptions.ConnectionName)

	return fx.Module(module, fx.Provide(fx.Annotate(func(lc fx.Lifecycle) (*publisher.Publisher[T], error) {
		ctx, cancel := context.WithCancel(context.Background())

		pub, err := publisher.New(ctx, connectionOptions, exchangeName, options...)
		if err != nil {
			cancel()
			return nil, err
		}

		lc.Append(fx.StopHook(func(ctx context.Context) error {
			cancel()
			return pub.CloseWithContext(ctx)
		}))

		return pub, err
	},
		fx.ResultTags(
			GetPublisherParamName(connectionOptions.ConnectionName, exchangeName)),
			fx.As(new(publisher.Pub[T]),
		),
	)))
}
