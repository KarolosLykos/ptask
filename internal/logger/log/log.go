package log

import (
	"context"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/KarolosLykos/ptask/internal/logger"
)

type logruslog struct {
	logger *logrus.Logger
}

var defaultLogger = &logrus.Logger{
	Out:          os.Stderr,
	Hooks:        make(logrus.LevelHooks),
	ReportCaller: false,
	ExitFunc:     os.Exit,
	Level:        logrus.DebugLevel,
	Formatter:    &logrus.JSONFormatter{},
}

func New(l *logrus.Logger) logger.Logger {
	return &logruslog{logger: l}
}

func Default(debug bool, formatter string) logger.Logger {
	l := defaultLogger

	l.SetFormatter(setFormatter(formatter))

	if debug {
		l.SetLevel(setLevel("trace"))
	} else {
		l.SetLevel(setLevel("info"))
	}

	return New(l)
}

func setFormatter(format string) logrus.Formatter {
	if strings.ToLower(format) == "json" {
		return &logrus.JSONFormatter{}
	}

	return &logrus.TextFormatter{}
}

func setLevel(lvl string) logrus.Level {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return logrus.InfoLevel
	}

	return level
}

func (l *logruslog) SetLevel(lvl string) {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		l.logger.SetLevel(logrus.InfoLevel)
	}

	l.logger.SetLevel(level)
}

func (l *logruslog) Trace(ctx context.Context, msg ...interface{}) {
	le := l.parseMessages(ctx, nil)
	le.Trace(msg...)
}

func (l *logruslog) Debug(ctx context.Context, msg ...interface{}) {
	le := l.parseMessages(ctx, nil)
	le.Debug(msg...)
}

func (l *logruslog) Info(ctx context.Context, msg ...interface{}) {
	le := l.parseMessages(ctx, nil)
	le.Info(msg...)
}

func (l *logruslog) Warn(ctx context.Context, err error, msg ...interface{}) {
	le := l.parseMessages(ctx, err)
	le.Warn(msg...)
}

func (l *logruslog) Error(ctx context.Context, err error, msg ...interface{}) {
	le := l.parseMessages(ctx, err)
	le.Error(msg...)

	le.Trace()
}

func (l *logruslog) Panic(ctx context.Context, err error, msg ...interface{}) {
	le := l.parseMessages(ctx, err)
	le.Panic(msg...)
}

func (l *logruslog) parseMessages(_ context.Context, err error) *logrus.Entry {
	e := l.logger.WithFields(logrus.Fields{"service": "ptask"})

	if err != nil {
		e = e.WithField("err", err.Error())
	}

	return e
}
