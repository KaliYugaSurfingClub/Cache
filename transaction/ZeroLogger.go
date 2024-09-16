package transaction

import (
	"cache/core"
	"context"
)

type ZeroLogger struct{}

func (tl *ZeroLogger) WriteEvent(core.EventType, string, string) {}

func (tl *ZeroLogger) Shutdown(context.Context) error { return nil }

func (tl *ZeroLogger) ReadEvents() (<-chan core.Event, <-chan error) {
	events := make(chan core.Event)
	errs := make(chan error)

	close(events)
	close(errs)

	return events, errs
}

func (tl *ZeroLogger) Start() <-chan error {
	errs := make(chan error)

	close(errs)

	return errs
}
