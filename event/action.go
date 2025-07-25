package event

import (
	"encoding/json"
)

// Action represents the type of action performed in an event
type Action string

const (
	ActionError    Action = "error"
	ActionFiring   Action = "firing"
	ActionResolved Action = "resolved"
	ActionAlert    Action = "alert"
)

// String returns the string representation of the Action
func (a *Action) String() string {
	return string(*a)
}

// MarshalJSON implements json.Marshaler interface
func (a *Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*a))
}

// UnmarshalJSON implements json.Unmarshaler interface
func (a *Action) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*a = Action(s)
	return nil
}
