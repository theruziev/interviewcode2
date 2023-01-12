package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/theruziev/oson_auth/internal/db"
	"github.com/theruziev/oson_auth/internal/event/constants"
	outboxevent "github.com/theruziev/oson_auth/internal/event/outbox"
	apphttp "github.com/theruziev/oson_auth/internal/http"
	"github.com/theruziev/oson_auth/internal/pkg/auth"
	"github.com/theruziev/oson_auth/internal/pkg/closer"
	"github.com/theruziev/oson_auth/internal/pkg/dbx"
	"github.com/theruziev/oson_auth/internal/pkg/httpx"
	"github.com/theruziev/oson_auth/internal/pkg/logging"
	"github.com/theruziev/oson_auth/internal/pkg/rabbitmqx"
	"github.com/theruziev/oson_auth/internal/pkg/validatorx"
	"github.com/theruziev/oson_auth/internal/service"
	"github.com/wagslane/go-rabbitmq"
	"golang.org/x/sync/errgroup"
)

type Option struct {
	Server   httpx.ServerOpts
	Postgres dbx.PostgresOpt
	RabbitMQ rabbitmqx.RabbitMQOpts
	Auth     auth.AuthOption
	IsDebug  bool
}

type HTTPServer struct {
	opt     *Option
	dbxPool *dbx.Dbx

	router chi.Router

	userService    *service.UserService
	contentService *service.ContentService

	userStore    *db.UserStore
	outboxStore  *db.OutBoxStore
	contentStore *db.ContentStore

	userHandler *apphttp.UserHandler

	rabbitmqConn *rabbitmq.Conn

	closer *closer.Closer
}

func NewApp(opt *Option) *HTTPServer {
	return &HTTPServer{
		opt:    opt,
		closer: closer.NewCloser(),
	}
}

func (s *HTTPServer) initPostgres(ctx context.Context) error {
	opt := s.opt
	dbxPool := dbx.NewDbx()
	err := dbxPool.Connect(ctx, opt.Postgres.DSN)
	if err != nil {
		return err
	}
	s.dbxPool = dbxPool
	s.closer.AddCloser(func(ctx context.Context) error {
		return s.dbxPool.Close(ctx)
	})
	return nil
}

func (s *HTTPServer) InitStore(_ context.Context) error {
	s.userStore = db.NewUserStore(s.dbxPool)
	s.outboxStore = db.NewOutBoxStore(s.dbxPool)
	s.contentStore = db.NewContentStore(s.dbxPool)
	return nil
}

func (s *HTTPServer) initService(_ context.Context) error {
	otp := auth.NewOtpConfig(&s.opt.Auth.Otp)
	s.userService = service.NewUserStore(&s.opt.Auth, s.outboxStore, s.userStore, otp)
	s.contentService = service.NewContentService(s.contentStore, s.dbxPool)
	return nil
}

func (s *HTTPServer) initHandler(_ context.Context) error {
	s.userHandler = apphttp.NewUserHandler(s.userService)
	return nil
}

func (s *HTTPServer) initRabbitMQ(_ context.Context) error {
	var err error
	s.rabbitmqConn, err = rabbitmqx.Connect(s.opt.RabbitMQ)
	return err
}

func (s *HTTPServer) serveOutbox(ctx context.Context) error {
	publisher, err := rabbitmq.NewPublisher(s.rabbitmqConn,
		rabbitmqx.DefaultWithPublisherOptions(ctx, constants.ExchangeUser)...,
	)
	if err != nil {
		return err
	}

	s.closer.AddCloser(func(ctx context.Context) error {
		publisher.Close()
		return nil
	})

	o := outboxevent.NewOutBox(publisher, s.outboxStore)
	o.Serve(ctx)
	return nil
}

func (s *HTTPServer) Init(ctx context.Context) error {
	if err := s.initRabbitMQ(ctx); err != nil {
		return fmt.Errorf("failed to init rabbitmq: %w", err)
	}

	if err := s.initPostgres(ctx); err != nil {
		return fmt.Errorf("failed to init mongo: %w", err)
	}

	if err := s.InitStore(ctx); err != nil {
		return fmt.Errorf("failed to init store: %w", err)
	}

	if err := s.initService(ctx); err != nil {
		return fmt.Errorf("failed to init service: %w", err)
	}

	if err := s.initHandler(ctx); err != nil {
		return fmt.Errorf("failed to init handler: %w", err)
	}

	s.initRouter(ctx)
	return nil
}

func (s *HTTPServer) initRouter(ctx context.Context) {
	logger := logging.FromContext(ctx)
	validator := validatorx.FromContext(ctx)
	authMiddleware := auth.Middleware(s.opt.Auth.JWTSecret)
	r := chi.NewRouter()
	r.Use(httpx.Recoverer(logger))
	r.Use(httpx.PopulateLogger(logger))
	r.Use(httpx.PopulateValidator(validator))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSONOKResponse(w)
	})

	tfaCheckMiddleware := chi.Middlewares{
		authMiddleware,
		auth.CheckScope(auth.TwoFACheckScope),
	}

	userMiddleware := chi.Middlewares{
		authMiddleware,
		auth.CheckScope(auth.UserScope),
	}

	r.Route("/user", func(r chi.Router) {
		r.Post("/activate/{aid}", s.userHandler.Activate)
		r.Post("/auth", s.userHandler.Auth)
		r.With(tfaCheckMiddleware...).Post("/auth-2fa", s.userHandler.AuthTwoFA)
		r.Post("/register", s.userHandler.Register)

		r.Get("/reset-password", s.userHandler.GetByResetPassword)
		r.Post("/reset-password", s.userHandler.ResetPasswordRequest)
		r.Put("/reset-password", s.userHandler.ResetPassword)
		r.Group(func(r chi.Router) {
			r.Use(userMiddleware...)
			r.Get("/me", s.userHandler.Me)
			r.Post("/change-password", s.userHandler.ChangePassword)
		})
		r.Route("/otp", func(r chi.Router) {
			r.Use(userMiddleware...)
			r.Post("/step1", s.userHandler.RequestEnableOTPStep1)
			r.Post("/step2", s.userHandler.RequestEnableOTPStep2)
			r.Post("/disable", s.userHandler.DisableOTP)
		})
	})

	s.router = r
}

func (s *HTTPServer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	logger := logging.FromContext(ctx)
	if err := s.Init(ctx); err != nil {
		return err
	}
	srv := http.Server{
		Addr:              s.opt.Server.Listen,
		ReadHeaderTimeout: s.opt.Server.ReadHeaderTimeout,
		Handler:           s.router,
	}
	g, childCtx := errgroup.WithContext(ctx)
	s.closer.AddCloser(srv.Shutdown)
	logger.Infof("run server on %s", s.opt.Server.Listen)
	g.Go(func() error {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	g.Go(func() error {
		if err := s.serveOutbox(childCtx); err != nil {
			return err
		}

		return nil
	})

	go func() {
		defer cancel()
		if err := g.Wait(); err != nil {
			logger.Errorf("failed to run: %s", err)
		}

	}()
	<-ctx.Done()

	closeCtx, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()
	if err := s.Close(closeCtx); err != nil {
		logger.Errorf("failed to shutdown app: %s", err)
	}
	return nil
}

func (s *HTTPServer) Close(ctx context.Context) error {
	return s.closer.Close(ctx)
}
