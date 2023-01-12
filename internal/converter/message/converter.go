package message

import (
	"github.com/theruziev/oson_auth/internal/model"
	v0 "github.com/theruziev/oson_auth/pkg/events/v0"
)

func ToUserEvent(user *model.User) *v0.UserEvent {
	var userStatus v0.UserStatus
	switch user.Status {
	case model.UserStatusActivate:
		userStatus = v0.UserStatusActivate
	case model.UserStatusRegistered:
		userStatus = v0.UserStatusRegistered
	}
	return &v0.UserEvent{
		PublicID:  user.PublicID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Status:    userStatus,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserRegisteredEvent(user *model.User) *v0.UserRegisteredEvent {
	return &v0.UserRegisteredEvent{
		PublicID:       user.PublicID,
		Email:          user.Email,
		ActivationCode: user.ActivationCode,
	}
}

func ToUserResetPasswordEvent(user *model.User, newResetCode string) *v0.UserResetPasswordEvent {
	return &v0.UserResetPasswordEvent{
		PublicID:          user.PublicID,
		Email:             user.Email,
		ResetPasswordCode: newResetCode,
	}
}
