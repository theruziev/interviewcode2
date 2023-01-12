package httpx

import (
	"time"
)

type ServerOpts struct {
	Listen            string        `help:"listen string" default:":3000" env:"LISTEN"`
	ReadHeaderTimeout time.Duration `help:"listen string" default:"10s" env:"READ_HEADER_TIMEOUT"`
}
