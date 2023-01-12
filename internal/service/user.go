package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/theruziev/oson_auth/internal/converter/message"
	"github.com/theruziev/oson_auth/internal/db"
	"github.com/theruziev/oson_auth/internal/event/constants"
	"github.com/theruziev/oson_auth/internal/model"
	"github.com/theruziev/oson_auth/internal/pkg/auth"
	"github.com/theruziev/oson_auth/internal/pkg/dbx"
	"github.com/theruziev/oson_auth/internal/pkg/errz"
)

type UserService struct {
	userStore   *db.UserStore
	authOpt     *auth.AuthOption
	outboxStore *db.OutBoxStore
	otp         *auth.Otp
}

func NewUserStore(authOpt *auth.AuthOption, outboxStore *db.OutBoxStore, userStore *db.UserStore, otp *auth.Otp) *UserService {
	return &UserService{
		authOpt:     authOpt,
		userStore:   userStore,
		outboxStore: outboxStore,
		otp:         otp,
	}
}

func (s *UserService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	user := &model.User{
		PublicID:       uuid.New().String(),
		ActivationCode: uuid.New().String(),
		Email:          req.Email,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Status:         model.UserStatusRegistered,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userStore.Insert(ctx, user); err != nil {
		if dbx.IsDuplicateErr(err) {
			return nil, errz.ConflictErr.Wrap(err)
		}
		return nil, err
	}

	if err := s.sendEventNewUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) sendEventNewUser(ctx context.Context, user *model.User) error {
	userEvent := message.ToUserEvent(user)

	messages := make([]*model.OutBox, 0)
	messages = append(messages, &model.OutBox{
		Topic:     constants.TopicUserChanged,
		Data:      userEvent,
		Status:    model.CreatedStatus,
		CreatedAt: time.Now(),
	})

	userRegisteredEvent := message.ToUserRegisteredEvent(user)
	messages = append(messages, &model.OutBox{
		Topic:     constants.TopicRegisteredUser,
		Data:      userRegisteredEvent,
		Status:    model.CreatedStatus,
		CreatedAt: time.Now(),
	})

	if err := s.outboxStore.Add(ctx, messages...); err != nil {
		return err
	}
	return nil
}

func (s *UserService) Activate(ctx context.Context, activationCode string) error {
	user, err := s.userStore.GetByActivationCode(ctx, activationCode)
	if err != nil {
		return err
	}
	if err := s.userStore.Activate(ctx, user.PublicID); err != nil {
		return err
	}

	user, err = s.GetByID(ctx, user.PublicID)
	if err != nil {
		return err
	}
	userEvent := message.ToUserEvent(user)
	if err := s.outboxStore.Add(ctx, &model.OutBox{
		Topic:     constants.TopicUserChanged,
		Data:      userEvent,
		Status:    model.CreatedStatus,
		CreatedAt: time.Now(),
	}); err != nil {
		return err
	}
	return nil
}

func (s *UserService) ResetPasswordRequest(ctx context.Context, email string) error {
	user, err := s.userStore.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	resetCode := uuid.New().String()
	if err := s.userStore.ResetPasswordRequest(ctx, user.PublicID, resetCode); err != nil {
		return err
	}

	resetEvent := message.ToUserResetPasswordEvent(user, resetCode)
	msg := &model.OutBox{
		Topic:     constants.TopicUserResetPassword,
		Data:      resetEvent,
		Status:    model.CreatedStatus,
		CreatedAt: time.Now(),
	}

	if err := s.outboxStore.Add(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetByResetPassword(ctx context.Context, resetCode string) (*model.User, error) {
	user, err := s.userStore.GetByResetPassword(ctx, resetCode)
	if err != nil {
		if dbx.IsErrNoRows(err) {
			return nil, errz.NotFoundErr.Wrap(err)
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) ResetPassword(ctx context.Context, resetCode, password string) error {
	user, err := s.userStore.GetByResetPassword(ctx, resetCode)
	if err != nil {
		return err
	}
	if err := s.userStore.ChangePassword(ctx, user.PublicID, password); err != nil {
		return err
	}
	return nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID, password string) error {
	err := s.userStore.ChangePassword(ctx, userID, password)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.userStore.GetByEmail(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, publicID string) (*model.User, error) {
	user, err := s.userStore.Get(ctx, publicID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
