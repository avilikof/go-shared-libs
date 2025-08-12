package alerting

import (
	"os"

	natsdriver "github.com/avilikof/go-shared-libs/nats"

	"github.com/nats-io/nats.go"
)

// StreamHandler provides methods for handling message streaming operations
// including publishing and subscribing to topics through a NATS connection.
type StreamHandler struct {
	natsDriver *natsdriver.NatsConnection
}

// NewStreamHandler creates and initializes a new StreamHandler instance.
// It establishes a connection to the NATS server using the provided URL environment variable.
// If the environment variable is not set, it falls back to the default NATS URL.
// Returns an error if the connection to NATS server fails.
func NewStreamHandler(url string) (*StreamHandler, error) {
	natsUrl := os.Getenv(url)
	if natsUrl == "" {
		natsUrl = nats.DefaultURL
	}

	natsConnection, err := natsdriver.NewNatsConnection(natsUrl)
	if err != nil {
		return nil, err
	}
	return &StreamHandler{
		natsDriver: natsConnection,
	}, nil
}

// Subscribe sets up a subscription to receive messages from a specified topic.
// Messages received from the topic are forwarded to the provided channel.
// Returns an error if the subscription setup fails.
func (sh *StreamHandler) Subscribe(topic string, alertChan chan []byte) error {
	pubSub := natsdriver.NewPubSub(sh.natsDriver)
	err := pubSub.Subscribe(topic, func(msg *nats.Msg) {
		alertChan <- msg.Data
	})
	if err != nil {
		return err
	}
	return nil
}

// Publish sends data to a specified topic through the message streaming system.
// Returns an error if the publishing operation fails.
func (sh *StreamHandler) Publish(topic string, data []byte) error {
	pubSub := natsdriver.NewPubSub(sh.natsDriver)
	err := pubSub.Publish(topic, data)
	if err != nil {
		return err
	}
	return nil
}

// Close terminates the connection to the NATS server and releases associated resources.
// This method should be called when the StreamHandler is no longer needed.
func (sh *StreamHandler) Close() {
	sh.natsDriver.Close()
}
