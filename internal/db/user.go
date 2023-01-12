package db

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/theruziev/oson_auth/internal/model"
	"github.com/theruziev/oson_auth/internal/pkg/dbx"
)

const usersTable = "users"

var defaultUserFields = []string{
	"id",
	"public_id",
	"name",
	"email",
	"password",
	"status",
	"created_at",
	"updated_at",
	"concat(activation_code, '') as activation_code",
	"concat(reset_password_code, '') as reset_password_code",
	"otp_secret",
	"otp_recovery_codes",
	"otp_enabled",
}

type UserStore struct {
	db dbx.Querier
}

func NewUserStore(db dbx.Querier) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) Insert(ctx context.Context, user *model.User) error {
	builder := pgsql.Insert(usersTable).SetMap(map[string]interface{}{
		"public_id":           user.PublicID,
		"first_name":          user.FirstName,
		"last_name":           user.LastName,
		"email":               user.Email,
		"password":            user.Password,
		"status":              user.Status,
		"created_at":          user.CreatedAt,
		"updated_at":          user.UpdatedAt,
		"activation_code":     user.ActivationCode,
		"reset_password_code": user.ResetPasswordCode,
		"otp_secret":          user.OtpSecret,
		"otp_recovery_codes":  user.OtpRecoveryCodes,
		"otp_enabled":         user.OtpEnabled,
	}).Suffix("returning id")

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	err = pgxscan.Get(ctx, s.db, user, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Get(ctx context.Context, publicID string) (*model.User, error) {
	builder := pgsql.Select(
		defaultUserFields...,
	).From(usersTable).Where(squirrel.Eq{"public_id": publicID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	var user model.User
	if err := pgxscan.Get(ctx, s.db, &user, query, args...); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	builder := pgsql.Select(
		defaultUserFields...,
	).From(usersTable).Where(squirrel.Eq{"email": email})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	var user model.User
	if err := pgxscan.Get(ctx, s.db, &user, query, args...); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) GetByActivationCode(ctx context.Context, activationCode string) (*model.User, error) {
	builder := pgsql.Select(
		defaultUserFields...,
	).From(usersTable).Where(squirrel.Eq{"activation_code": activationCode})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	var user model.User
	if err := pgxscan.Get(ctx, s.db, &user, query, args...); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) Activate(ctx context.Context, publicID string) error {
	builder := pgsql.Update(usersTable).SetMap(map[string]interface{}{
		"activation_code": nil,
		"status":          model.UserStatusActivate,
	}).Where(squirrel.Eq{"public_id": publicID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	conn, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if conn.RowsAffected() == 0 {
		return fmt.Errorf("failed to update")
	}
	return nil
}

func (s *UserStore) GetByResetPassword(ctx context.Context, resetCode string) (*model.User, error) {
	builder := pgsql.Select(
		defaultUserFields...,
	).From(usersTable).Where(squirrel.Eq{"reset_code": resetCode})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	var user model.User
	if err := pgxscan.Get(ctx, s.db, &user, query, args...); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) ChangePassword(ctx context.Context, publicID, newPassword string) error {
	builder := pgsql.Update(usersTable).SetMap(map[string]interface{}{
		"password": newPassword,
	}).Where(squirrel.Eq{"public_id": publicID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	conn, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if conn.RowsAffected() == 0 {
		return fmt.Errorf("failed to update")
	}
	return nil
}

func (s *UserStore) ResetPasswordRequest(ctx context.Context, publicID, resetCode string) error {
	builder := pgsql.Update(usersTable).SetMap(map[string]interface{}{
		"reset_password_code": resetCode,
		"updated_at":          time.Now(),
	}).Where(squirrel.Eq{"public_id": publicID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	conn, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if conn.RowsAffected() == 0 {
		return fmt.Errorf("failed to update")
	}
	return nil
}

func (s *UserStore) RemoveCodeFromRecoveryCode(ctx context.Context, publicID, code string) error {
	builder := pgsql.Update(usersTable).SetMap(map[string]interface{}{
		"otp_recovery_codes": squirrel.Expr("array_remove(otp_recovery_codes, ?)", code),
		"updated_at":         time.Now(),
	}).Where(squirrel.Eq{"public_id": publicID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	conn, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if conn.RowsAffected() == 0 {
		return fmt.Errorf("failed to update")
	}
	return nil
}

func (s *UserStore) SetOtpSecret(ctx context.Context, publicID, otpSecret string) error {
	builder := pgsql.Update(usersTable).SetMap(map[string]interface{}{
		"otp_secret": otpSecret,
		"updated_at": time.Now(),
	}).Where(squirrel.Eq{"public_id": publicID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	conn, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if conn.RowsAffected() == 0 {
		return fmt.Errorf("failed to update")
	}
	return nil
}

func (s *UserStore) SetOtpRecoveryCodes(ctx context.Context, publicID string, codes []string) error {
	builder := pgsql.Update(usersTable).SetMap(map[string]interface{}{
		"otp_recovery_codes": codes,
		"updated_at":         time.Now(),
	}).Where(squirrel.Eq{"public_id": publicID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	conn, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if conn.RowsAffected() == 0 {
		return fmt.Errorf("failed to update")
	}
	return nil
}

func (s *UserStore) SetOtpEnabled(ctx context.Context, publicID string, isEnabled bool) error {
	builder := pgsql.Update(usersTable).SetMap(map[string]interface{}{
		"otp_enabled": isEnabled,
		"updated_at":  time.Now(),
	}).Where(squirrel.Eq{"public_id": publicID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	conn, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if conn.RowsAffected() == 0 {
		return fmt.Errorf("failed to update")
	}
	return nil
}
