package alerting

import "time"

// Storage defines the interface for persisting and retrieving alert data.
// Implementations should provide key-value storage with expiration support.
type Storage interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, expires time.Duration) error
}

// Stream defines the interface for publish-subscribe messaging operations.
// Implementations should support publishing alerts to topics and subscribing to receive them.
type Stream interface {
	Publish(topic string, data []byte) error
	Subscribe(topic string, channel chan []byte) error
}
