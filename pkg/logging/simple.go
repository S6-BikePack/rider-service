package logging

import (
	"context"
	"fmt"
	"rider-service/config"
)

type SimpleLogger struct {
	Config *config.Config
}

func NewSimpleLogger(cfg *config.Config) (*SimpleLogger, error) {
	return &SimpleLogger{Config: cfg}, nil
}

func (l *SimpleLogger) Close() error {
	return nil
}

func (l *SimpleLogger) Panic(ctx context.Context, args ...interface{}) {
	panic(args)
}

func (l *SimpleLogger) Fatal(ctx context.Context, args ...interface{}) {
	fmt.Println("FATAL: ", args)
}

func (l *SimpleLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	fmt.Printf("INFO: "+msg+"\n", keysAndValues)
}

func (l *SimpleLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	fmt.Printf("DEBUG: "+msg+"\n", keysAndValues)
}

func (l *SimpleLogger) Warning(ctx context.Context, msg string, keysAndValues ...interface{}) {
	fmt.Printf("WARNING: "+msg+"\n", keysAndValues)
}

func (l *SimpleLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	fmt.Printf("ERROR: "+msg+"\n", keysAndValues)
}
