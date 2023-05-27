package infrastructure

import (
	"github.com/go-stack/stack"
	"github.com/inconshreveable/log15"
	"os"
)

type Log15Logger struct {
	isDebugMode bool
	logger      log15.Logger
}

func (l *Log15Logger) Debug(msg string, ctx ...interface{}) {
	ctx = append(ctx, "caller", stack.Caller(1).String())
	l.logger.Debug(msg, ctx...)
}

func (l *Log15Logger) Info(msg string, ctx ...interface{}) {
	ctx = append(ctx, "caller", stack.Caller(1).String())
	l.logger.Info(msg, ctx...)
}

func (l *Log15Logger) Warn(msg string, ctx ...interface{}) {
	ctx = append(ctx, "caller", stack.Caller(1).String())
	l.logger.Warn(msg, ctx...)
}

func (l *Log15Logger) Error(msg string, ctx ...interface{}) {
	ctx = append(ctx, "caller", stack.Caller(1).String())
	l.logger.Error(msg, ctx...)
}

func (l *Log15Logger) Crit(msg string, ctx ...interface{}) {
	ctx = append(ctx, "caller", stack.Caller(1).String())
	l.logger.Crit(msg, ctx...)
}

func (l *Log15Logger) IsDebugMode() bool {
	return l.isDebugMode
}

func NewLog15Logger(isDebugMode bool, exitOnCrit bool, extraHandlers []log15.Handler) ILogger {
	logLevel := log15.LvlInfo
	if isDebugMode {
		logLevel = log15.LvlDebug
	}

	logger := log15.Root().New()

	handlers := []log15.Handler{
		log15.StdoutHandler,
	}
	handlers = append(handlers, extraHandlers...)

	if exitOnCrit {
		handlers = append(handlers, log15.FuncHandler(
			func(r *log15.Record) error {
				if r.Lvl == log15.LvlCrit {
					os.Exit(1)
				}
				return nil
			},
		))
	}
	logger.SetHandler(log15.LvlFilterHandler(logLevel, log15.MultiHandler(handlers...)))

	return &Log15Logger{
		isDebugMode: isDebugMode,
		logger:      logger,
	}
}
