package errz

import (
	"errors"
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	e := fmt.Errorf("hello world")
	wrapErr := fmt.Errorf("bad error: %w", e)
	wrapWrapErr := ConflictErr.Wrap(wrapErr)
	err := BadRequestErr.Wrap(wrapWrapErr)

	fmt.Println(errors.Is(err, NotFoundErr))
}
