package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/theruziev/oson_auth/app/useremail"
	useremailconsumer "github.com/theruziev/oson_auth/internal/event/consumer/useremail"
	"github.com/theruziev/oson_auth/internal/pkg/logging"
	"github.com/theruziev/oson_auth/internal/pkg/mailgunx"
	"github.com/theruziev/oson_auth/internal/pkg/rabbitmqx"
)

type userEmail struct {
	RabbitMQOpt         rabbitmqx.RabbitMQOpts        `embed:"" prefix:"rabbitmq." envprefix:"RABBITMQ_" validate:"required,dive,required"`
	UserMailConsumerOpt useremailconsumer.ConsumerOpt `embed:"" prefix:"user-mail-consumer." envprefix:"USER_MAIL_CONSUMER_" validate:"required,dive,required"`
	MailgunOpt          mailgunx.MailgunOpt           `embed:"" prefix:"mailgun." envprefix:"MAILGUN_" validate:"required,dive,required"`
}

func (s *userEmail) Run(cliCtx *Ctx) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	logger := logging.NewLogger(cliCtx.LogLevel, cliCtx.IsDebug)
	ctx = logging.WithLogger(ctx, logger)
	defer cancel()
	go func() {
		sig := <-sigs
		logger.Warnf("interrupt signal: %s", sig)
		cancel()
	}()

	app := useremail.NewUserEmailApp(&useremail.Option{
		RabbitMQ:             s.RabbitMQOpt,
		UserEmailConsumerOpt: s.UserMailConsumerOpt,
		MailgunOpt:           s.MailgunOpt,
		IsDebug:              cliCtx.IsDebug,
	})

	return app.Run(ctx)
}
