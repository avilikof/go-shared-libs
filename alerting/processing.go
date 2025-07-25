// Package alerting provides alert processing functionality for managing
// alert lifecycle, including storage, firing, resolution, and change detection.
package alerting

import (
	"fmt"
	"reflect"
	"time"

	"github.com/alex/go-shared-libs/alerts"
)

const StorageTopic = "alert.store"
const EvetTopic = "alert.event"

type Processor struct {
	input   <-chan *alerts.Alert
	storage Storage
	stream  Stream
}

func NewProcessor(input <-chan *alerts.Alert, storage Storage, stream Stream) *Processor {
	return &Processor{
		input:   input,
		storage: storage,
		stream:  stream,
	}
}

func (p *Processor) Process() {
	for alert := range p.input {
		storedAlertBytes, err := p.storage.Get(alert.ID)
		if err != nil {
			err := p.storeNewAlert(alert)
			if err != nil {
				println(err.Error())
				logEvent := logEvent(err, alert.ID)
				err := p.stream.Publish(EvetTopic, logEvent.Bytes())
				if err != nil {
					panic(err)
				}
			}
			continue
		}
		storedAlert, err := alerts.AlertFromBytes(storedAlertBytes)
		if err != nil {
			fmt.Printf("Error decoding alert %s: %v\n", alert.ID, err)
			logEvent := logEvent(err, alert.ID)
			err := p.stream.Publish(EvetTopic, logEvent.Bytes())
			if err != nil {
				panic(err)
			}
			continue
		}

		if reflect.DeepEqual(storedAlert, alert) {
			fmt.Println("Alerts are same")
		} else {
			if alert.IsFiring() != storedAlert.IsFiring() {
				if !alert.IsFiring() {
					err := p.resolveAlert(alert)
					if err != nil {
						println(err.Error())
						logEvent := logEvent(err, alert.ID)
						err := p.stream.Publish(EvetTopic, logEvent.Bytes())
						if err != nil {
							panic(err)
						}
					}
					continue
				} else {
					err := p.fireAlert(alert)
					if err != nil {
						return
					}
					continue
				}
			} else {
				fmt.Printf("%s is different", diffAlerts(alert, storedAlert))
			}
		}
	}
}

func (p *Processor) resolveAlert(alert *alerts.Alert) error {
	alert.Resolve(time.Now())
	alert.Firing = false
	err := p.storeAlert(alert)
	if err != nil {
		return err
	}
	fmt.Printf("Alert :: %s resolved\n", alert.ID)

	resolveEvent := resolvedEvent(alert.ID)
	err = p.stream.Publish(EvetTopic, resolveEvent.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (p *Processor) fireAlert(alert *alerts.Alert) error {
	alert.Firing = true
	err := p.storeAlert(alert)
	if err != nil {
		println(err.Error())
		return err
	}
	firingEvent := firingEvent(alert.ID)
	err = p.stream.Publish(EvetTopic, firingEvent.Bytes())
	if err != nil {
		return err
	}
	fmt.Printf("Alert :: %s fired\n", alert.ID)
	return nil
}

func (p *Processor) storeNewAlert(alert *alerts.Alert) error {
	const alertNotStoredMsg = "alert not stored, new alert with Resolved status"
	if !alert.IsFiring() {
		fmt.Println(alertNotStoredMsg)
		logEvent := logEvent(fmt.Errorf(alertNotStoredMsg), alert.ID)
		err := p.stream.Publish(EvetTopic, logEvent.Bytes())
		if err != nil {
			return err
		}
		return nil
	}
	err := p.storeAlert(alert)
	if err != nil {
		return err
	}
	return nil
}

func (p *Processor) storeAlert(alert *alerts.Alert) error {
	return p.stream.Publish(StorageTopic, alert.Bytes())
}

func diffAlerts(a1, a2 *alerts.Alert) []string {
	var diffs []string

	if a1.Title != a2.Title {
		diffs = append(diffs, fmt.Sprintf("Title differs: %s vs %s\n", a1.Title, a2.Title))
	}
	if a1.Message != a2.Message {
		diffs = append(diffs, fmt.Sprintf("Message differs: %s vs %s\n", a1.Message, a2.Message))
	}
	if !a1.Timestamp.StartTime.Equal(a2.Timestamp.StartTime) {
		if startsLater(a1.Timestamp.StartTime, a2.Timestamp.StartTime) {
			fmt.Printf("Alert :: %s :: refiring\n", a1.Id())
			return []string{a1.Id(), a2.Id()}
		}
	}
	if !a1.Timestamp.EndTime.Equal(a2.Timestamp.EndTime) {
		diffs = append(diffs, fmt.Sprintf("EndTime differs: %v vs %v\n", a1.Timestamp.EndTime, a2.Timestamp.EndTime))
	}
	if a1.Firing != a2.Firing {
		diffs = append(diffs, fmt.Sprintf("Firing differs: %v vs %v\n", a1.Firing, a2.Firing))
	}
	if a1.Acknowledged != a2.Acknowledged {
		diffs = append(diffs, fmt.Sprintf("Acknowledged differs: %v vs %v\n", a1.Acknowledged, a2.Acknowledged))
	}

	return diffs
}

func startsLater(t1, t2 time.Time) bool {
	return t1.After(t2)
}
