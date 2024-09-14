package transaction

import (
	"cache/core"
	"context"
)

type BackgroundLogger struct{}

func (tl *BackgroundLogger) WriteEvent(t core.EventType, key string, value string) {}
func (tl *BackgroundLogger) ReadEvents() (<-chan core.Event, <-chan error)         { return nil, nil }
func (tl *BackgroundLogger) Start() <-chan error                                   { return nil }
func (tl *BackgroundLogger) Shutdown(ctx context.Context) error                    { return nil }
