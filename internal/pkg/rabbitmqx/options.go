package rabbitmqx

import (
	"context"

	"github.com/theruziev/oson_auth/internal/pkg/logging"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

func WithConsumerOptionsQueueQuorum(options *rabbitmq.ConsumerOptions) {
	if options.QueueOptions.Args == nil {
		options.QueueOptions.Args = rabbitmq.Table{}
	}
	options.QueueOptions.Args["x-queue-type"] = "quorum"
}

func WithConsumerOptionsLogger(logger *zap.SugaredLogger) func(options *rabbitmq.ConsumerOptions) {
	return func(options *rabbitmq.ConsumerOptions) {
		options.Logger = newLoggerx(logger)
	}
}

func WithPublisherOptionsLogger(logger *zap.SugaredLogger) func(options *rabbitmq.PublisherOptions) {
	return func(options *rabbitmq.PublisherOptions) {
		options.Logger = newLoggerx(logger)
	}
}

func DefaultWithConsumerOptions(ctx context.Context, exchangeName, routingKey string) []func(options *rabbitmq.ConsumerOptions) {
	return []func(options *rabbitmq.ConsumerOptions){
		WithConsumerOptionsLogger(logging.FromContext(ctx)),
		rabbitmq.WithConsumerOptionsExchangeName(exchangeName),
		rabbitmq.WithConsumerOptionsRoutingKey(routingKey),
		rabbitmq.WithConsumerOptionsQueueDurable,
		WithConsumerOptionsQueueQuorum,
	}
}

func DefaultWithPublisherOptions(ctx context.Context, exchangeName string) []func(options *rabbitmq.PublisherOptions) {
	return []func(options *rabbitmq.PublisherOptions){
		WithPublisherOptionsLogger(logging.FromContext(ctx)),
		rabbitmq.WithPublisherOptionsExchangeName(exchangeName),
		rabbitmq.WithPublisherOptionsExchangeDurable,
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	}
}

func WithPublishJSONContentType() func(options *rabbitmq.PublishOptions) {
	return rabbitmq.WithPublishOptionsContentType("application/json")
}
