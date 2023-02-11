package logger

import (
	"context"
)

type Logger interface {
	SetLevel(lvl string)
	Trace(ctx context.Context, msg ...interface{})
	Debug(ctx context.Context, msg ...interface{})
	Info(ctx context.Context, msg ...interface{})
	Warn(ctx context.Context, err error, msg ...interface{})
	Error(ctx context.Context, err error, msg ...interface{})
	Panic(ctx context.Context, err error, msg ...interface{})
}
