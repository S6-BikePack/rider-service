package logging

import "context"

type MockLogger struct{}

func (MockLogger) Close() error {
	return nil
}

func (MockLogger) Panic(ctx context.Context, args ...interface{}) {
	return
}

func (MockLogger) Fatal(ctx context.Context, args ...interface{}) {
	return
}

func (MockLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	return
}

func (MockLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	return
}

func (MockLogger) Warning(ctx context.Context, msg string, keysAndValues ...interface{}) {
	return
}

func (MockLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	return
}
