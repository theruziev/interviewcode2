package rabbitmqx

import "github.com/wagslane/go-rabbitmq"

type RabbitMQOpts struct {
	URL string `help:"rabbitmq address" required:"" env:"URL"`
}

func Connect(opts RabbitMQOpts) (*rabbitmq.Conn, error) {
	conn, err := rabbitmq.NewConn(opts.URL)
	return conn, err
}
