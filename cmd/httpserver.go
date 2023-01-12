package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	apppkg "github.com/theruziev/oson_auth/app/httpserver"
	"github.com/theruziev/oson_auth/internal/pkg/auth"
	"github.com/theruziev/oson_auth/internal/pkg/dbx"
	"github.com/theruziev/oson_auth/internal/pkg/httpx"
	"github.com/theruziev/oson_auth/internal/pkg/logging"
	"github.com/theruziev/oson_auth/internal/pkg/rabbitmqx"
)

type httpserver struct {
	ServerOpts   httpx.ServerOpts       `embed:"" prefix:"http." envprefix:"HTTP_" validate:"required,dive,required"`
	PostgresOpts dbx.PostgresOpt        `embed:"" prefix:"postgres." envprefix:"POSTGRES_" validate:"required,dive,required"`
	RabbitMQOpt  rabbitmqx.RabbitMQOpts `embed:"" prefix:"rabbitmq." envprefix:"RABBITMQ_" validate:"required,dive,required"`
	AuthOpts     auth.AuthOption        `embed:"" prefix:"auth." envprefix:"AUTH_" validate:"required,dive,required"`
}

func (s *httpserver) Run(cliCtx *Ctx) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := logging.NewLogger(cliCtx.LogLevel, cliCtx.IsDebug)
	ctx = logging.WithLogger(ctx, logger)
	go func() {
		sig := <-sigs
		logger.Warnf("interrupt signal: %s", sig)
		cancel()
	}()

	app := apppkg.NewApp(&apppkg.Option{
		Server:   s.ServerOpts,
		Postgres: s.PostgresOpts,
		RabbitMQ: s.RabbitMQOpt,
		Auth:     s.AuthOpts,
		IsDebug:  cliCtx.IsDebug,
	})

	return app.Run(ctx)
}
