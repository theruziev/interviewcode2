package logging

import (
	"testing"

	"go.uber.org/zap"
)

func TestDefaultLogger(t *testing.T) {
	log := NewLogger(zap.DebugLevel.String(), false)
	log.Info("Hello World")
	log.With(zap.String("key", "value")).Info("Hello World")
}
