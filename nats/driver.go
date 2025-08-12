package natsdriver

import (
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

var ErrInvalidUrl = errors.New("invalid url")

type NatsConnection struct {
	conn *nats.Conn
}

func NewNatsConnection(url string) (*NatsConnection, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("%w: %q: %v", ErrInvalidUrl, url, err)
	}
	return &NatsConnection{conn: conn}, nil
}

func (nc *NatsConnection) Close() {
	nc.conn.Close()
}

func setDefaultOptions() []nats.Option {
	return []nats.Option{
		nats.MaxReconnects(-1),
		nats.ReconnectWait(100 * time.Millisecond),
		nats.Timeout(10 * time.Second),
	}
}
