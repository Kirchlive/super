# Event Handler Example

This example demonstrates event-driven architecture using NATS JetStream for the OpenCode-SuperClaude integration.

## 🎯 Overview

This example shows:
- NATS JetStream setup for event streaming
- Event publishing and subscription patterns
- Configuration hot-reload using events
- Plugin lifecycle event handling
- Error handling and retry mechanisms

## 📁 Structure

```
event-handler/
├── README.md              # This file
├── Makefile              # Build automation
├── docker-compose.yml    # NATS server setup
├── publisher/            # Event publisher
│   ├── main.go          # Publisher entry point
│   └── events.go        # Event definitions
├── subscriber/           # Event subscriber
│   ├── main.go          # Subscriber entry point
│   └── handlers.go      # Event handlers
└── shared/               # Shared code
    ├── events.go        # Event type definitions
    └── nats.go          # NATS helper functions
```

## 🚀 Quick Start

### Prerequisites
```bash
# Install Docker for NATS server
# Or install NATS server locally:
brew install nats-server
```

### Build and Run
```bash
# Start NATS server
docker-compose up -d

# Build everything
make build

# Run subscriber (in terminal 1)
make run-subscriber

# Run publisher (in terminal 2)
make run-publisher

# Run full demo
make demo
```

### Expected Output
```
[SUBSCRIBER] Connected to NATS
[SUBSCRIBER] Subscribed to: plugin.>
[SUBSCRIBER] Subscribed to: config.>
[PUBLISHER] Publishing plugin.loaded event
[SUBSCRIBER] Received: plugin.loaded - Plugin hello v1.0.0 loaded
[PUBLISHER] Publishing config.changed event
[SUBSCRIBER] Received: config.changed - Hot reloading configuration
```

## 💻 Code Walkthrough

### 1. Event Types (`shared/events.go`)
```go
type PluginEvent struct {
    Type      string    `json:"type"`
    Plugin    string    `json:"plugin"`
    Version   string    `json:"version"`
    Timestamp time.Time `json:"timestamp"`
}

type ConfigEvent struct {
    Type      string                 `json:"type"`
    Changes   map[string]interface{} `json:"changes"`
    Timestamp time.Time             `json:"timestamp"`
}
```

### 2. Event Publisher (`publisher/main.go`)
- Connects to NATS JetStream
- Publishes various event types
- Demonstrates event patterns

### 3. Event Subscriber (`subscriber/handlers.go`)
- Subscribes to event streams
- Handles different event types
- Implements retry logic

## 🔧 NATS Configuration

### JetStream Setup
```go
// Create stream
js.AddStream(&nats.StreamConfig{
    Name:     "OPENCODE",
    Subjects: []string{"plugin.>", "config.>", "command.>"},
    Retention: nats.WorkQueuePolicy,
    MaxAge:    24 * time.Hour,
})
```

### Consumer Configuration
```go
// Create durable consumer
js.AddConsumer("OPENCODE", &nats.ConsumerConfig{
    Durable:       "opencode-worker",
    DeliverPolicy: nats.DeliverNewPolicy,
    AckPolicy:     nats.AckExplicitPolicy,
    MaxDeliver:    3,
})
```

## 📊 Event Patterns

### Plugin Lifecycle Events
- `plugin.loaded` - Plugin successfully loaded
- `plugin.unloaded` - Plugin unloaded
- `plugin.error` - Plugin error occurred
- `plugin.reloaded` - Plugin hot-reloaded

### Configuration Events
- `config.changed` - Configuration updated
- `config.reloaded` - Full configuration reload
- `config.error` - Configuration error

### Command Events
- `command.requested` - Command execution requested
- `command.completed` - Command completed successfully
- `command.failed` - Command execution failed

## 🧪 Testing

```bash
# Unit tests
make test

# Integration tests
make test-integration

# Load test
make load-test
```

## 🔐 Security

- TLS connection to NATS
- Authentication with credentials
- Authorization per subject
- Message encryption option

## 📈 Monitoring

The example includes:
- Event metrics (published/consumed)
- Error rates and retry counts
- Latency measurements
- Queue depth monitoring

## 📚 Next Steps

1. Implement event sourcing
2. Add event replay capability
3. Implement saga patterns
4. Add distributed tracing
5. Integrate with plugin system

## 🔗 References

- [NATS JetStream Documentation](https://docs.nats.io/jetstream)
- [Event API Specification](../../API-SPECIFICATIONS/event-api.md)
- [EDA.md](../../EDA.md) - Event-driven architecture plan