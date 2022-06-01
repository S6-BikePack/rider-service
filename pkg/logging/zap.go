package logging

import (
	"context"
	"rider-service/config"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type OtelzapSugaredLogger struct {
	Logger otelzap.SugaredLogger
	Config *config.Config
}

func NewSugaredOtelZap(cfg *config.Config) (*OtelzapSugaredLogger, error) {
	logger, err := zap.NewDevelopment()

	if err != nil {
		return nil, err
	}

	otelZap := otelzap.New(logger)
	sugar := otelZap.Sugar()

	return &OtelzapSugaredLogger{Logger: *sugar, Config: cfg}, nil
}

func (l *OtelzapSugaredLogger) Close() error {
	err := l.Logger.Sync()
	if err != nil {
		return err
	}

	return nil
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
