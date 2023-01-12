package validatorx

import (
	"context"
	"sync"

	"github.com/go-playground/validator/v10"
)

type contextType string

const validatorKey = contextType("validator")

var (
	defaultValidator     *validator.Validate
	defaultValidatorOnce sync.Once
)

func NewValidator() *validator.Validate {
	return validator.New()
}

func DefaultValidator() *validator.Validate {
	defaultValidatorOnce.Do(func() {
		defaultValidator = NewValidator()
	})

	return defaultValidator
}

func WithValidator(ctx context.Context, v *validator.Validate) context.Context {
	return context.WithValue(ctx, validatorKey, v)
}

func FromContext(ctx context.Context) *validator.Validate {
	if validate, ok := ctx.Value(validatorKey).(*validator.Validate); ok {
		return validate
	}

	return DefaultValidator()
}
