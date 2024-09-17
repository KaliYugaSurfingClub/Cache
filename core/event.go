package core

type EventType = byte

const (
	EventDelete EventType = iota
	EventPut
	EventClear
)

type Event struct {
	ID    uint64
	Type  EventType
	Key   string
	Value string
}
