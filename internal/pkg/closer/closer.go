package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type closer func(ctx context.Context) error

type Closer struct {
	closers []closer
	mu      sync.Mutex
}

func NewCloser() *Closer {
	return &Closer{
		closers: make([]closer, 0),
	}
}

func (c *Closer) AddCloser(cl closer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closers = append(c.closers, cl)
}

func (c *Closer) Close(ctx context.Context) error {
	errors := make([]string, 0)
	for i := len(c.closers) - 1; i > 0; i-- {
		fn := c.closers[i]
		if err := fn(ctx); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errors, "\n"))

}
