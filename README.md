# Go Shared Libraries

Shared Go libraries extracted from monorepo for distributed microservices architecture.

## Libraries:
- `alerts/` - Alert data structures and utilities
- `alerting/` - Alert processing and stream handling utilities
- `cfg_manager/` - Configuration management utilities
- `event/` - Event handling structures and utilities
- `games/` - Game-related utilities and logic
- `logger/` - Logging utilities and drivers
- `nats/` - NATS messaging client and utilities
- `redis/` - Redis client driver and utilities
- `redpanda/` - Redpanda producer utilities

## Usage

```go
import "github.com/avilikof/go-shared-libs/alerts"
import "github.com/avilikof/go-shared-libs/logger"
```

## Testing

```bash
make test              # Run all tests
make test-alerts       # Test specific package
make lint              # Run linter
# Check if ready for release
make release-check

# Create and push your first version tag
make tag VERSION=v1.0.0
```
