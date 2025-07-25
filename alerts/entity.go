package alerts

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Alert struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Message      string    `json:"message"`
	Timestamp    Timestamp `json:"timestamp"`
	Firing       bool      `json:"firing"`
	Acknowledged bool      `json:"acknowledged"`
}

type Timestamp struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

func NewAlert(id, title, message string, startTime time.Time, firing bool) *Alert {
	return &Alert{
		ID:           id,
		Title:        title,
		Message:      message,
		Timestamp:    Timestamp{StartTime: startTime},
		Firing:       firing,
		Acknowledged: false,
	}
}

func AlertFromString(rawAlert string) (*Alert, error) {
	var alert Alert
	err := json.Unmarshal([]byte(rawAlert), &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

func AlertFromBytes(data []byte) (*Alert, error) {
	var alert Alert
	err := json.Unmarshal(data, &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

func (a *Alert) IsFiring() bool {
	return a.Firing
}

func (a *Alert) Resolve(endTime time.Time) {
	a.Firing = false
	a.Timestamp.EndTime = endTime
}

func (a *Alert) Acknowledge() {
	a.Acknowledged = true
}

func (a *Alert) Id() string {
	return a.ID
}

func (a *Alert) String() string {
	return fmt.Sprintf("Alert %s: %s: %s", a.ID, a.Message, a.Timestamp.StartTime.Format(time.RFC3339))
}

func (a *Alert) JSON() string {
	return fmt.Sprintf(`{"id": "%s", "title": "%s", "message": "%s", "firing": "%t", "acknowledged": "%t", "timestamp": {"startTime": "%s", "endTime": "%s"}}`, a.ID, a.Title, a.Message, a.Firing, a.Acknowledged, a.Timestamp.StartTime.Format(time.RFC3339), a.Timestamp.EndTime.Format(time.RFC3339))
}

func (a *Alert) Bytes() []byte {
	jsonStr, err := json.Marshal(a)
	if err != nil {
		return nil
	}
	return jsonStr
}

func (a *Alert) Hash() string {
	hash := sha256.Sum256([]byte(a.JSON()))
	return hex.EncodeToString(hash[:])
}
