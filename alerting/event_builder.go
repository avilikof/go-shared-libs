package alerting

import (
	"github.com/alex/go-shared-libs/event"
	"time"
)

func eventBuilder(alertID string, action event.Action, eventType event.Type, additionalData map[string]any) event.Event {
	newEvent := event.NewEvent("alerting", eventType, time.Now(), action, map[string]any{})
	newEvent.TimestampNow()

	if newEvent.Type == event.TypeLog {
		newEvent.TypeLog()
	} else {
		newEvent.TypeEvent()
	}

	newEvent.Action = event.Action(action)
	newEvent.Message = map[string]any{"alert_id": alertID}

	for key, value := range additionalData {
		newEvent.Message[key] = value
	}

	return newEvent
}
func logEvent(errorMessage error, alertID string) event.Event {
	return eventBuilder(alertID, event.ActionError, event.TypeLog, map[string]any{"error_message": errorMessage.Error()})
}

func firingEvent(alertID string) event.Event {
	return eventBuilder(alertID, event.ActionFiring, event.TypeEvent, nil)
}

func resolvedEvent(alertID string) event.Event {
	return eventBuilder(alertID, event.ActionResolved, event.TypeEvent, nil)
}
