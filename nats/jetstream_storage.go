package natsdriver

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type JetStreamStorage struct {
	js nats.JetStreamContext
	kv nats.KeyValue
}

func NewJetStreamStorage(nc *nats.Conn, bucketName string) (*JetStreamStorage, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Create or get KV bucket for alert storage
	kv, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:      bucketName,
		Description: "Alert storage bucket",
		TTL:         24 * time.Hour,    // Default TTL for alerts
		MaxBytes:    100 * 1024 * 1024, // 100MB max
		Storage:     nats.FileStorage,
		Replicas:    1,
	})
	if err != nil {
		// Try to get existing bucket if creation fails
		kv, err = js.KeyValue(bucketName)
		if err != nil {
			return nil, fmt.Errorf("failed to create/get KV bucket: %w", err)
		}
	}

	return &JetStreamStorage{
		js: js,
		kv: kv,
	}, nil
}

// Get retrieves a value by key
func (jss *JetStreamStorage) Get(key string) ([]byte, error) {
	entry, err := jss.kv.Get(key)
	if err != nil {
		if err == nats.ErrKeyNotFound {
			return nil, fmt.Errorf("key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return entry.Value(), nil
}

// Set stores a value with optional expiration
func (jss *JetStreamStorage) Set(key string, value []byte, expires time.Duration) error {
	// JetStream KV doesn't support per-key TTL, but we can use bucket TTL
	// For custom TTL, we'd need to implement cleanup logic
	_, err := jss.kv.Put(key, value)
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// Delete removes a key
func (jss *JetStreamStorage) Delete(key string) error {
	err := jss.kv.Delete(key)
	if err != nil && err != nats.ErrKeyNotFound {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// GetAll returns all keys (for debugging/admin)
func (jss *JetStreamStorage) GetAll() ([]string, error) {
	keys := make([]string, 0)

	// Watch all keys to get current state
	watcher, err := jss.kv.WatchAll()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Stop()

	// Collect existing keys with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		select {
		case entry := <-watcher.Updates():
			if entry == nil {
				return keys, nil
			}
			if entry.Operation() != nats.KeyValueDelete {
				keys = append(keys, entry.Key())
			}
		case <-ctx.Done():
			return keys, nil
		}
	}
}

// Close cleans up resources
func (jss *JetStreamStorage) Close() error {
	// KV buckets don't need explicit closing
	return nil
}

// Health check
func (jss *JetStreamStorage) Ping() error {
	// Try a simple operation to verify connectivity
	_, err := jss.kv.Status()
	return err
}
