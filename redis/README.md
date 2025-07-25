# Redis Driver for Go

A simple and lightweight Go package for interacting with Redis. This driver provides a clean API for basic Redis operations such as setting, getting, and deleting keys.

## Installation

```bash
go get github.com/go-github/pkg/redis
```

## Usage

### Connecting to Redis

```go
package main

import (
    "fmt"
    "log"

    "github.com/go-github/pkg/redis"
)

func main() {
    // Connect to Redis at localhost:6379, database 0
    driver, err := redisdriver.NewDriver("localhost:6379", 0)
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    defer driver.Close() // Don't forget to close the connection when done

    fmt.Println("Successfully connected to Redis")
}
```

### Setting Values

```go
// Set a key with a 10-minute expiration
err := driver.Set("user:123", []byte(`{"name":"John","age":30}`), 10*time.Minute)
if err != nil {
    log.Fatalf("Failed to set key: %v", err)
}
```

### Getting Values

```go
value, err := driver.Get("user:123")
if err != nil {
    if err.Error() == redisdriver.ErrKeyNotFound.Error() {
        fmt.Println("Key not found")
    } else {
        log.Fatalf("Error retrieving key: %v", err)
    }
} else {
    fmt.Printf("Value: %s\n", value)
}
```

### Retrieving All Keys

```go
keys, err := driver.GetAll()
if err != nil {
    log.Fatalf("Failed to get all keys: %v", err)
}

fmt.Println("Redis keys:")
for _, key := range keys {
    fmt.Println(key)
}
```

### Deleting All Keys

```go
err := driver.DeleteAll()
if err != nil {
    log.Fatalf("Failed to delete all keys: %v", err)
}
fmt.Println("All keys deleted successfully")
```

## Error Handling

The package provides predefined errors for common scenarios:

- `ErrKeyNotFound`: Returned when trying to get a key that doesn't exist
- `ErrConnectionFailed`: Returned when the connection to Redis fails

Example of proper error handling:

```go
value, err := driver.Get("nonexistent-key")
if err != nil {
    if err.Error() == redisdriver.ErrKeyNotFound.Error() {
        // Handle the case of key not found
    } else {
        // Handle other errors
    }
}
```

## Configuration

When creating a new driver instance, you can specify:

- Redis server address (e.g., "localhost:6379")
- Database number (an integer)

```go
driver, err := redisdriver.NewDriver("redis.example.com:6379", 1)
```

## Thread Safety

The Redis driver is safe for concurrent use by multiple goroutines.

## Limitations

- Currently only supports basic operations (Get, Set, GetAll, DeleteAll)
- Does not support Redis clusters or sentinel
- No support for Redis transactions

## License

This project is licensed under the MIT License - see your organization's licensing information.
