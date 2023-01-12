package rabbitmqx

import "go.uber.org/zap"

type loggerx struct {
	logger *zap.SugaredLogger
}

func newLoggerx(logger *zap.SugaredLogger) *loggerx {
	return &loggerx{
		logger: logger,
	}
}

func (r *loggerx) Fatalf(template string, args ...interface{}) {
	r.logger.Fatalf(template, args...)
}

func (r *loggerx) Errorf(template string, args ...interface{}) {
	r.logger.Errorf(template, args...)
}

func (r *loggerx) Warnf(template string, args ...interface{}) {
	r.logger.Warnf(template, args...)
}

func (r *loggerx) Infof(template string, args ...interface{}) {
	r.logger.Infof(template, args...)
}

func (r *loggerx) Debugf(template string, args ...interface{}) {
	r.logger.Debugf(template, args...)
}

func (r *loggerx) Tracef(template string, args ...interface{}) {
	r.logger.Infof(template, args...)
}
