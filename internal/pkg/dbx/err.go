package dbx

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func IsDuplicateErr(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgerrcode.IsIntegrityConstraintViolation(pgErr.Code)
	}
	return false
}

func IsErrNoRows(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, pgx.ErrNoRows)
}
