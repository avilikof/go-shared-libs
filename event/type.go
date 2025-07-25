package event

import "encoding/json"

// Type represents the type of event
type Type string

const (
	TypeLog   Type = "LOG"
	TypeEvent Type = "EVENT"
)

// String returns the string representation of the Type
func (t *Type) String() string {
	return string(*t)
}

// MarshalJSON implements json.Marshaler interface
func (t *Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*t))
}

// UnmarshalJSON implements json.Unmarshaler interface
func (t *Type) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = Type(s)
	return nil
}
