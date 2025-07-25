package alerting

import (
	"fmt"
	"os"

	natsdriver "github.com/alex/go-shared-libs/nats"

	"github.com/nats-io/nats.go"
)

const (
	NATS_URL_JS  = "NATS_URL"
	ALERT_STREAM = "ALERTS"
	EVENT_STREAM = "EVENTS"
)

type JetStreamHandler struct {
	conn    *nats.Conn
	js      nats.JetStreamContext
	storage *natsdriver.JetStreamStorage
}

func NewJetStreamHandler() (*JetStreamHandler, error) {
	natsUrl := os.Getenv(NATS_URL_JS)
	if natsUrl == "" {
		natsUrl = "nats://localhost:4222"
	}

	// Connect to NATS with JetStream enabled
	conn, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create JetStream context
	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Create streams for alerts and events
	err = createStreams(js)
	if err != nil {
		return nil, fmt.Errorf("failed to create streams: %w", err)
	}

	// Create storage using JetStream KV
	storage, err := natsdriver.NewJetStreamStorage(conn, "alerts")
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream storage: %w", err)
	}

	return &JetStreamHandler{
		conn:    conn,
		js:      js,
		storage: storage,
	}, nil
}

func createStreams(js nats.JetStreamContext) error {
	// Create alert stream
	_, err := js.AddStream(&nats.StreamConfig{
		Name:        ALERT_STREAM,
		Description: "Alert processing stream",
		Subjects:    []string{"test.alert", "alert.store"},
		Storage:     nats.FileStorage,
		MaxAge:      24 * 60 * 60 * 1000000000, // 24 hours in nanoseconds
		MaxBytes:    100 * 1024 * 1024,         // 100MB
		Replicas:    1,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return fmt.Errorf("failed to create alert stream: %w", err)
	}

	// Create event stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:        EVENT_STREAM,
		Description: "Event logging stream",
		Subjects:    []string{"alert.event"},
		Storage:     nats.FileStorage,
		MaxAge:      7 * 24 * 60 * 60 * 1000000000, // 7 days in nanoseconds
		MaxBytes:    50 * 1024 * 1024,              // 50MB
		Replicas:    1,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return fmt.Errorf("failed to create event stream: %w", err)
	}

	return nil
}

// Subscribe to JetStream subjects
func (jsh *JetStreamHandler) Subscribe(subject string, alertChan chan []byte) error {
	// Create durable consumer
	consumerName := fmt.Sprintf("%s-consumer", subject)

	sub, err := jsh.js.Subscribe(subject, func(msg *nats.Msg) {
		alertChan <- msg.Data
		msg.Ack() // Acknowledge message processing
	}, nats.Durable(consumerName))

	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}

	// Keep subscription alive
	go func() {
		select {}
	}()

	_ = sub // Keep reference to prevent GC
	return nil
}

// Publish to JetStream
func (jsh *JetStreamHandler) Publish(subject string, data []byte) error {
	_, err := jsh.js.Publish(subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish to %s: %w", subject, err)
	}
	return nil
}

// Get storage interface
func (jsh *JetStreamHandler) Storage() Storage {
	return jsh.storage
}

// Close connections
func (jsh *JetStreamHandler) Close() {
	if jsh.storage != nil {
		jsh.storage.Close()
	}
	if jsh.conn != nil {
		jsh.conn.Close()
	}
}
