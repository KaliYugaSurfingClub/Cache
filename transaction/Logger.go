package transaction

type EventType byte

const (
	EventDelete EventType = iota
	EventPut
)

type Event struct {
	Sequence uint64
	Type     EventType
	Key      string
	Value    []byte
}

type Logger interface {
	WritePut(key string, value []byte)
	WriteDelete(key string)
	ReadEvents() (<-chan Event, <-chan error)
	ErrCh() <-chan error
	Start()
	Close() error
}
