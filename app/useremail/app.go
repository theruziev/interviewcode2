package useremail

import (
	"context"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/theruziev/oson_auth/internal/event/constants"
	"github.com/theruziev/oson_auth/internal/event/consumer/useremail"
	"github.com/theruziev/oson_auth/internal/pkg/closer"
	"github.com/theruziev/oson_auth/internal/pkg/logging"
	"github.com/theruziev/oson_auth/internal/pkg/mailgunx"
	"github.com/theruziev/oson_auth/internal/pkg/rabbitmqx"
	"github.com/wagslane/go-rabbitmq"
)

type Option struct {
	RabbitMQ             rabbitmqx.RabbitMQOpts
	UserEmailConsumerOpt useremail.ConsumerOpt
	MailgunOpt           mailgunx.MailgunOpt
	IsDebug              bool
}

type UserEmailApp struct {
	opt *Option

	mailgunClient *mailgun.MailgunImpl

	userEmailConsumer *useremail.ConsumerHandler

	closer *closer.Closer

	rabbitmqConn *rabbitmq.Conn
}

func NewUserEmailApp(opt *Option) *UserEmailApp {
	return &UserEmailApp{
		opt:    opt,
		closer: closer.NewCloser(),
	}
}

func (a *UserEmailApp) init(ctx context.Context) error {
	if err := a.initServices(ctx); err != nil {
		return err
	}

	if err := a.initRabbitMQ(ctx); err != nil {
		return err
	}
	return nil
}

func (a *UserEmailApp) initServices(_ context.Context) error {
	a.mailgunClient = mailgunx.NewMailgun(a.opt.MailgunOpt)
	a.userEmailConsumer = useremail.NewConsumerHandler(a.opt.UserEmailConsumerOpt, a.mailgunClient)

	return nil
}

func (a *UserEmailApp) initRabbitMQ(_ context.Context) error {
	conn, err := rabbitmqx.Connect(a.opt.RabbitMQ)
	if err != nil {
		return err
	}
	a.rabbitmqConn = conn
	return nil
}

func (a *UserEmailApp) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	logger := logging.FromContext(ctx)
	if err := a.init(ctx); err != nil {
		return err
	}

	consumerWelcome, err := rabbitmq.NewConsumer(
		a.rabbitmqConn,
		a.userEmailConsumer.WelcomeMessageEmail(ctx),
		constants.QueueWelcomeEmail,
		rabbitmqx.DefaultWithConsumerOptions(ctx,
			constants.ExchangeUser,
			constants.TopicRegisteredUser,
		)...,
	)
	if err != nil {
		return err
	}
	a.closer.AddCloser(func(ctx context.Context) error {
		consumerWelcome.Close()
		return nil
	})
	consumerResetPassword, err := rabbitmq.NewConsumer(
		a.rabbitmqConn,
		a.userEmailConsumer.ResetPasswordEmail(ctx),
		constants.QueueResetPasswordEmail,
		rabbitmqx.DefaultWithConsumerOptions(ctx,
			constants.ExchangeUser,
			constants.TopicUserResetPassword,
		)...,
	)
	if err != nil {
		return err
	}
	a.closer.AddCloser(func(ctx context.Context) error {
		consumerResetPassword.Close()
		return nil
	})

	<-ctx.Done()
	closeCtx, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()
	if err := a.Close(closeCtx); err != nil {
		logger.Errorf("failed to shutdown app: %s", err)
	}
	return nil
}

func (a *UserEmailApp) Close(ctx context.Context) error {
	return a.closer.Close(ctx)
}
