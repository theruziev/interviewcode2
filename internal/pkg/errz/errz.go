package errz

import (
	"fmt"
)

type CustomError struct {
	msg  string
	code int
	err  error
}

func NewCustomError(msg string, code int) *CustomError {
	return &CustomError{
		msg:  msg,
		code: code,
	}
}

func (c *CustomError) Error() string {
	if c.err != nil {
		return fmt.Sprintf("%s: %s", c.msg, c.err)
	}
	return c.msg
}

func (c *CustomError) Unwrap() error {
	return c.err
}

func (c *CustomError) Wrap(err error) *CustomError {
	if err == nil {
		return nil
	}

	return &CustomError{
		msg:  c.msg,
		err:  err,
		code: c.code,
	}
}

func (c *CustomError) Is(err error) bool {
	if v, ok := err.(*CustomError); ok {
		if v.code == c.code {
			return true
		}
	}
	return false
}

func (c *CustomError) New(format string, args ...any) *CustomError {
	msg := fmt.Sprintf(format, args...)
	return NewCustomError(msg, c.code)
}

var (
	BadRequestErr = NewCustomError("bad request", 100)
	ConflictErr   = NewCustomError("conflict error", 101)
	NotFoundErr   = NewCustomError("not found", 102)
	InternalErr   = NewCustomError("internal error", 103)
)
