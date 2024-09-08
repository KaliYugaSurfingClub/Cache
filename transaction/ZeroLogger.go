package transaction

type ZeroLogger struct{}

func (tl *ZeroLogger) WritePut(key string, value []byte)        {}
func (tl *ZeroLogger) WriteDelete(key string)                   {}
func (tl *ZeroLogger) ReadEvents() (<-chan Event, <-chan error) { return nil, nil }
func (tl *ZeroLogger) ErrCh() <-chan error                      { return nil }
func (tl *ZeroLogger) Start()                                   {}
func (tl *ZeroLogger) Close() error                             { return nil }
