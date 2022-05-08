package logging

import (
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type Logger interface {
	Close()
	Panic(ctx context.Context, args ...interface{})
	Fatal(ctx context.Context, args ...interface{})
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
	Warning(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
}

type OtelzapSugaredLogger struct {
	Logger otelzap.SugaredLogger
}

func NewSugaredOtelZap() (*OtelzapSugaredLogger, error) {
	logger, err := zap.NewDevelopment()

	if err != nil {
		return nil, err
	}

	otelZap := otelzap.New(logger)
	sugar := otelZap.Sugar()

	return &OtelzapSugaredLogger{Logger: *sugar}, nil
}

func (l *OtelzapSugaredLogger) Close() {
	err := l.Logger.Sync()
	if err != nil {
		l.Logger.Error(err)
	}
}

func (l *OtelzapSugaredLogger) Panic(ctx context.Context, args ...interface{}) {
	l.Logger.Panic(args...)
}

func (l *OtelzapSugaredLogger) Fatal(ctx context.Context, args ...interface{}) {
	l.Logger.Fatal(args...)
}

func (l *OtelzapSugaredLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.Logger.InfowContext(ctx, msg, keysAndValues...)
}

func (l *OtelzapSugaredLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.Logger.Ctx(ctx).Debugw(msg, keysAndValues...)
}

func (l *OtelzapSugaredLogger) Warning(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.Logger.WarnwContext(ctx, msg, keysAndValues...)
}

func (l *OtelzapSugaredLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.Logger.ErrorwContext(ctx, msg, keysAndValues...)
}
