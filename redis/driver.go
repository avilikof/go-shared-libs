package redisdriver

import (
	"context"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

var (
	ErrKeyNotFound      = fmt.Errorf("key not found")
	ErrConnectionFailed = fmt.Errorf("failed to connect to redis")
)

type Driver struct {
	client *redis.Client
}

func NewDriver(addr string, db int) (*Driver, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrConnectionFailed, err)
	}

	return &Driver{
		client: client,
	}, nil
}

func (d *Driver) Close() error {
	return d.client.Close()
}

func (d *Driver) Get(key string) ([]byte, error) {
	return d.client.Get(context.Background(), key).Bytes()
}

func (d *Driver) Set(key string, value []byte, expiration time.Duration) error {
	return d.client.Set(context.Background(), key, value, expiration).Err()
}

func (d *Driver) GetAll() ([]string, error) {
	keys, err := d.client.Keys(context.Background(), "*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get all keys: %w", err)
	}

	return keys, nil
}

func (d *Driver) DeleteAll() error {
	keys, err := d.client.Keys(context.Background(), "*").Result()
	if err != nil {
		return fmt.Errorf("failed to get all keys: %w", err)
	}

	for _, key := range keys {
		if err := d.client.Del(context.Background(), key).Err(); err != nil {
			return fmt.Errorf("failed to delete key %s: %w", key, err)
		}
	}

	return nil
}

func (d *Driver) Ping() error {
	return d.client.Ping(context.Background()).Err()
}
