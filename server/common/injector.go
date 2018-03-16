package common

import (
	"context"
)

type Injector interface {
	GAELogger(context.Context) Logger
}

type injector struct{}

func NewInjector() Injector {
	return new(injector)
}

func (i *injector) GAELogger(ctx context.Context) Logger {
	return &gaeLogger{ctx: ctx}
}
