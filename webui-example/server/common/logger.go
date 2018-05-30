package common

import (
	"context"

	"google.golang.org/appengine/log"
)

type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Critical(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type gaeLogger struct {
	ctx context.Context
}

func (g *gaeLogger) Debug(format string, v ...interface{}) {
	log.Debugf(g.ctx, format, v...)
}

func (g *gaeLogger) Info(format string, v ...interface{}) {
	log.Infof(g.ctx, format, v...)
}

func (g *gaeLogger) Error(format string, v ...interface{}) {
	log.Errorf(g.ctx, format, v...)
}

func (g *gaeLogger) Warn(format string, v ...interface{}) {
	log.Warningf(g.ctx, format, v...)
}

func (g *gaeLogger) Critical(format string, v ...interface{}) {
	log.Criticalf(g.ctx, format, v...)
}
