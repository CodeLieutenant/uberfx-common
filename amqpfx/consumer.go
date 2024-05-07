package amqpfx

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"

	"github.com/nano-interactive/go-amqp/v3/connection"
	"github.com/nano-interactive/go-amqp/v3/consumer"
)

func ConsumerModuleFunc[T consumer.Message](
	handler func(context.Context, T) error,
	queueOptions consumer.QueueDeclare,
	connectionOptions connection.Config,
	options ...consumer.Option[T],
) fx.Option {
	create := func(opts ...consumer.Option[T]) (consumer.Consumer[T], error) {
		return consumer.NewFunc(handler, connectionOptions, queueOptions, opts...)
	}

	return c(queueOptions, connectionOptions, create, options...)
}

func ConsumerModuleRaw[T consumer.Message](
	handler consumer.RawHandler,
	queueOptions consumer.QueueDeclare,
	connectionOptions connection.Config,
	options ...consumer.Option[T],
) fx.Option {
	create := func(opts ...consumer.Option[T]) (consumer.Consumer[T], error) {
		return consumer.NewRaw(handler, connectionOptions, queueOptions, opts...)
	}

	return c(queueOptions, connectionOptions, create, options...)
}

func ConsumerModuleRawFunc[T consumer.Message](
	handler func(context.Context, *amqp091.Delivery) error,
	queueOptions consumer.QueueDeclare,
	connectionOptions connection.Config,
	options ...consumer.Option[T],
) fx.Option {
	create := func(opts ...consumer.Option[T]) (consumer.Consumer[T], error) {
		return consumer.NewRawFunc(handler, connectionOptions, queueOptions, opts...)
	}

	return c(queueOptions, connectionOptions, create, options...)
}

func ConsumerModule[T consumer.Message](
	handler consumer.Handler[T],
	queueOptions consumer.QueueDeclare,
	connectionOptions connection.Config,
	options ...consumer.Option[T],
) fx.Option {
	create := func(opts ...consumer.Option[T]) (consumer.Consumer[T], error) {
		return consumer.New(handler, connectionOptions, queueOptions, opts...)
	}

	return c(queueOptions, connectionOptions, create, options...)
}

func c[T consumer.Message](
	queueOptions consumer.QueueDeclare,
	connectionOptions connection.Config,
	createConsumer func(...consumer.Option[T]) (consumer.Consumer[T], error),
	options ...consumer.Option[T],
) fx.Option {
	module := fmt.Sprintf("amqp-consumer-module-%s-%s", queueOptions.QueueName, connectionOptions.ConnectionName)
	name := fmt.Sprintf("amqp-consumer-%s-%s", queueOptions.QueueName, connectionOptions.ConnectionName)

	return fx.Module(
		module,
		fx.Provide(fx.Annotate(func() (consumer.Consumer[T], error) {
			opts := make([]consumer.Option[T], 0, len(options)+1)
			opts = append(opts, options...)

			c, err := createConsumer(opts...)
			if err != nil {
				return consumer.Consumer[T]{}, err
			}

			return c, nil
		}, fx.ResultTags(`name:"`+name+`"`))),
		fx.Invoke(fx.Annotate(func(lc fx.Lifecycle, c consumer.Consumer[T]) {
			lc.Append(fx.StartStopHook(
				func(ctx context.Context) {
					go func() {
						if err := c.Start(ctx); err != nil {
							panic(err)
						}
					}()
				},
				c.CloseWithContext,
			))
		},
			fx.ParamTags(`name:"`+name+`"`)),
		),
	)
}
