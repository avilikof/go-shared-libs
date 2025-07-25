package alerts

import (
	"encoding/json"
	"errors"
	"time"
)

type Action struct {
	Type             string `json:"type"`
	Target           string `json:"target"`
	AutoCreate       bool   `json:"auto_create,omitempty"`
	EscalationPolicy string `json:"escalation_policy,omitempty"`
}
type AlertState int

const (
	AlertStateActive AlertState = iota
	AlertStateResolved
)

func (s AlertState) String() string {
	switch s {
	case AlertStateActive:
		return "active"
	case AlertStateResolved:
		return "resolved"
	default:
		return "unknown"
	}
}

type AlertV2 struct {
	id               string
	source           string
	receivedAt       time.Time
	severity         string
	alertType        string
	message          string
	labels           map[string]string
	annotations      map[string]string
	deduplicationKey string
	correlationID    *string
	actions          []Action
	state            AlertState
}

func NewAlertV2(id, source, severity, alertType, message, dedupKey string, receivedAt time.Time, state AlertState) *AlertV2 {
	return &AlertV2{
		id:               id,
		source:           source,
		severity:         severity,
		alertType:        alertType,
		message:          message,
		deduplicationKey: dedupKey,
		receivedAt:       receivedAt,
		labels:           make(map[string]string),
		annotations:      make(map[string]string),
		state:            state,
	}
}
func (a *AlertV2) Labels() map[string]string      { return a.labels }
func (a *AlertV2) Annotations() map[string]string { return a.annotations }
func (a *AlertV2) DeduplicationKey() string       { return a.deduplicationKey }
func (a *AlertV2) CorrelationID() *string         { return a.correlationID }
func (a *AlertV2) Actions() []Action              { return a.actions }

// Controlled mutators
func (a *AlertV2) SetSeverity(sev string) error {
	switch sev {
	case "critical", "warning", "info":
		a.severity = sev
		return nil
	default:
		return errors.New("invalid severity level")
	}
}

func (a *AlertV2) AddLabel(key, value string) {
	a.labels[key] = value
}

func (a *AlertV2) AddAnnotation(key, value string) {
	a.annotations[key] = value
}

func (a *AlertV2) SetCorrelationID(id string) {
	a.correlationID = &id
}

func (a *AlertV2) AddAction(action Action) {
	a.actions = append(a.actions, action)
}

func (a *AlertV2) MarshalJSON() ([]byte, error) {
	type alias AlertV2 // prevent infinite loop
	return json.Marshal(&struct {
		ID               string            `json:"id"`
		Source           string            `json:"source"`
		ReceivedAt       time.Time         `json:"received_at"`
		Severity         string            `json:"severity"`
		Type             string            `json:"type"`
		Message          string            `json:"message"`
		Labels           map[string]string `json:"labels,omitempty"`
		Annotations      map[string]string `json:"annotations,omitempty"`
		DeduplicationKey string            `json:"deduplication_key"`
		CorrelationID    *string           `json:"correlation_id,omitempty"`
		Actions          []Action          `json:"actions,omitempty"`
	}{
		ID:               a.id,
		Source:           a.source,
		ReceivedAt:       a.receivedAt,
		Severity:         a.severity,
		Type:             a.alertType,
		Message:          a.message,
		Labels:           a.labels,
		Annotations:      a.annotations,
		DeduplicationKey: a.deduplicationKey,
		CorrelationID:    a.correlationID,
		Actions:          a.actions,
	})
}
