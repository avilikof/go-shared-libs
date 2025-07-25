package redpanda

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"sync"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

type RPCLIENT struct {
	client  *kgo.Client
	fetches *kgo.Fetches
}
type RPMessage struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       []byte
	Value     []byte
	Timestamp int64
}

func NewRPClient() (*RPCLIENT, error) {
	seeds := []string{os.Getenv("RP_URL")}
	opts := []kgo.Opt{}
	opts = append(opts,
		kgo.SeedBrokers(seeds...),
		kgo.ConsumeTopics("hello-world"),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)

	// Initialize public CAs for TLS
	opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))
	user := os.Getenv("RP_USER")
	password := os.Getenv("RP_PASSWORD")

	fmt.Printf("U:P -- %s:%s\n", user, password)

	// Initializes SASL/SCRAM 256
	opts = append(opts, kgo.SASL(scram.Auth{
		User: user,
		Pass: password,
	}.AsSha256Mechanism()))

	// Initializes SASL/SCRAM 512
	/*
		  opts = append(opts, kgo.SASL(scram.Auth{
				User: "<username>",
				Pass: "<password>",
			}.AsSha512Mechanism()))
	*/

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	fetches := client.PollFetches(context.Background())
	return &RPCLIENT{client: client, fetches: &fetches}, nil
}

func (r *RPCLIENT) Push(topic, key, message string) error {
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 1; i < 100; i++ {
		wg.Add(1)
		record := &kgo.Record{
			Topic: topic,
			Key:   []byte(key),
			Value: []byte(message),
		}
		r.client.Produce(ctx, record, func(record *kgo.Record, err error) {
			defer wg.Done()
			if err != nil {
				fmt.Printf("Error sending message: %v \n", err)
			} else {
				fmt.Printf("Message sent: topic: %s, offset: %d, value: %s \n",
					topic, record.Offset, record.Value)
			}
		})
	}
	wg.Wait()
	return nil
}

func (r *RPCLIENT) Pull(topic, key string) (RPMessage, error) {
	fetches := r.fetches

	if errs := fetches.Errors(); len(errs) > 0 {
		// All errors are retried internally when fetching, but non-retriable
		// errors are returned from polls so that users can notice and take
		// action.
		panic(fmt.Sprint(errs))
	}
	iter := fetches.RecordIter()

	record := iter.Next()
	topicInfo := fmt.Sprintf("topic: %s (%d|%d)",
		record.Topic, record.Partition, record.Offset)
	messageInfo := fmt.Sprintf("key: %s, Value: %s",
		record.Key, record.Value)
	fmt.Printf("Message consumed: %s, %s \n", topicInfo, messageInfo)

	msg := RPMessage{
		Topic:     record.Topic,
		Partition: record.Partition,
		Offset:    record.Offset,
		Key:       record.Key,
		Value:     record.Value,
	}
	return msg, nil
}

func (r *RPCLIENT) Ping() error {
	ctx := context.Background()
	err := r.client.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *RPCLIENT) Close() {
	r.client.Close()
}
