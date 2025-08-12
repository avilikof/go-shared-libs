package alerting

import (
	"os"

	natsdriver "github.com/avilikof/go-shared-libs/nats"

	"github.com/nats-io/nats.go"
)

type StreamHandler struct {
	natsDriver *natsdriver.NatsConnection
}

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

func (sh *StreamHandler) Publish(topic string, data []byte) error {
	pubSub := natsdriver.NewPubSub(sh.natsDriver)
	err := pubSub.Publish(topic, data)
	if err != nil {
		return err
	}
	return nil
}

func (sh *StreamHandler) Close() {
	sh.natsDriver.Close()
}
