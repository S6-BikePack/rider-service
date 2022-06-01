package logging

import "context"

type Logger interface {
	Close() error
	Panic(ctx context.Context, args ...interface{})
	Fatal(ctx context.Context, args ...interface{})
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
	Warning(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
}
