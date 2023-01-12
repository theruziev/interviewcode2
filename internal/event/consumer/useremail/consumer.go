package useremail

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-pkgz/repeater"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/theruziev/oson_auth/internal/email/template"
	"github.com/theruziev/oson_auth/internal/pkg/logging"
	v0 "github.com/theruziev/oson_auth/pkg/events/v0"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

const (
	repeats        = 3
	repeatsDelay   = 5 * time.Second
	resetSubject   = "Reset Password"
	welcomeSubject = "Welcome"
)

type ConsumerOpt struct {
	ActivationLinkTemplate string        `help:"kafka address" required:"" env:"ACTIVATION_LINK_FORMAT"`
	ResetPasswordFormat    string        `help:"kafka address" required:"" env:"RESET_PASSWORD_LINK_FORMAT"`
	Sender                 string        `help:"kafka address" required:"" env:"SENDER"`
	Timeout                time.Duration `help:"kafka address" required:"" env:"TIMEOUT"`
}

type ConsumerHandler struct {
	repeater      *repeater.Repeater
	opt           ConsumerOpt
	mailgunClient *mailgun.MailgunImpl
}

func NewConsumerHandler(opt ConsumerOpt, mailgunClient *mailgun.MailgunImpl) *ConsumerHandler {
	rp := repeater.NewDefault(repeats, repeatsDelay)
	return &ConsumerHandler{
		repeater:      rp,
		opt:           opt,
		mailgunClient: mailgunClient,
	}
}

func (c *ConsumerHandler) ResetPasswordEmail(ctx context.Context) func(message rabbitmq.Delivery) rabbitmq.Action {
	return func(message rabbitmq.Delivery) rabbitmq.Action {
		logger := logging.FromContext(ctx).With(
			zap.String("id", message.MessageId),
			zap.String("response", message.RoutingKey),
		)
		logger.Infof("new message")
		var resetPasswordEvent v0.UserResetPasswordEvent
		if err := json.Unmarshal(message.Body, &resetPasswordEvent); err != nil {
			logger.Error("failed to process json: %s", err)
			return rabbitmq.NackDiscard
		}

		link := fmt.Sprintf(c.opt.ResetPasswordFormat, resetPasswordEvent.ResetPasswordCode)
		resetPasswordBody, err := template.ResetPasswordEmail(resetPasswordEvent.Email, link)
		if err != nil {
			logger.Error("failed to create template: %s", err)
			return rabbitmq.NackDiscard
		}

		emailMsg := c.mailgunClient.NewMessage(c.opt.Sender, resetSubject, link, resetPasswordEvent.Email)
		emailMsg.SetHtml(resetPasswordBody)
		ctx, cancel := context.WithTimeout(ctx, c.opt.Timeout)
		defer cancel()
		err = c.repeater.Do(ctx, func() error {
			_, _, err = c.mailgunClient.Send(ctx, emailMsg)
			return err
		})
		if err != nil {
			logger.Error("failed to send email: %s", err)
			return rabbitmq.NackDiscard
		}

		logger.Debugf("email sended")
		return rabbitmq.Ack
	}
}

func (c *ConsumerHandler) WelcomeMessageEmail(ctx context.Context) func(message rabbitmq.Delivery) rabbitmq.Action {
	return func(message rabbitmq.Delivery) rabbitmq.Action {
		logger := logging.FromContext(ctx).With(
			zap.String("id", message.MessageId),
			zap.String("response", message.RoutingKey),
		)
		logger.Infof("new message")
		var userRegisteredMsg v0.UserRegisteredEvent
		if err := json.Unmarshal(message.Body, &userRegisteredMsg); err != nil {
			logger.Error("failed to process json: %s", err)
			return rabbitmq.NackDiscard
		}
		logger.Infof("%s", message.Headers)

		link := fmt.Sprintf(c.opt.ActivationLinkTemplate, userRegisteredMsg.ActivationCode)
		welcomeEmailBody, err := template.WelcomeEmail(userRegisteredMsg.Email, link)
		if err != nil {
			logger.Error("failed to create template: %s", err)
			return rabbitmq.NackRequeue
		}

		emailMsg := c.mailgunClient.NewMessage(c.opt.Sender, welcomeSubject, link, userRegisteredMsg.Email)
		emailMsg.SetHtml(welcomeEmailBody)

		ctx, cancel := context.WithTimeout(ctx, c.opt.Timeout)
		defer cancel()

		err = c.repeater.Do(ctx, func() error {
			_, _, err = c.mailgunClient.Send(ctx, emailMsg)
			return err
		})
		if err != nil {
			logger.Error("failed to send email: %s", err)
			return rabbitmq.NackDiscard
		}
		logger.Debugf("email sended")
		return rabbitmq.Ack
	}
}
