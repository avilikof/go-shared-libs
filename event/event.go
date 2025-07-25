package event

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	Service   string
	Type      Type
	Timestamp time.Time
	Action    Action
	Message   map[string]any
}

// Action represents the type of action performed in an event

func NewEvent(service string, eventType Type, timestamp time.Time, action Action, message map[string]any) Event {
	return Event{
		Service:   service,
		Type:      eventType,
		Timestamp: timestamp,
		Action:    action,
		Message:   message,
	}
}

func (e *Event) Bytes() []byte {
	jsonStr, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return jsonStr
}

// FromBytes creates an Event from a byte slice
func FromBytes(data []byte) (*Event, error) {
	var event Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal event: %w", err)
	}
	return &event, nil
}

func (e *Event) TypeLog() {
	e.Type = TypeLog
}
func (e *Event) TypeEvent() {
	e.Type = TypeEvent
}
func (e *Event) TimestampNow() {
	e.Timestamp = time.Now()
}
