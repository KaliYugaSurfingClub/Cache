package transaction

import (
	"cache/core"
	"context"
)

type ZeroLogger struct{}

func (tl *ZeroLogger) WriteEvent(core.EventType, string, string)     {}
func (tl *ZeroLogger) ReadEvents() (<-chan core.Event, <-chan error) { return nil, nil }
func (tl *ZeroLogger) Start() <-chan error                           { return nil }
func (tl *ZeroLogger) Shutdown(context.Context) error                { return nil }
