# Event System API Specification

Event-driven architecture specification for OpenCode SuperClaude integration using NATS JetStream/Kafka for asynchronous communication and reactive updates.

## Overview

The Event System provides asynchronous, decoupled communication between OpenCode components, plugins, and external services. It enables real-time updates, configuration changes, monitoring, and workflow orchestration through a publish-subscribe messaging pattern.

## Architecture

```
┌─────────────┐    Events    ┌─────────────┐    Topics    ┌─────────────┐
│ Publishers  │─────────────►│ Event Bus   │◄─────────────│ Subscribers │
│ (Plugins,   │              │ (NATS/Kafka)│              │ (Components)│
│  MCP Server)│              └─────────────┘              └─────────────┘
└─────────────┘                     │                            
                                    ▼                            
                              ┌─────────────┐                    
                              │ Persistence │                    
                              │ (JetStream) │                    
                              └─────────────┘                    
```

## Event Bus Technologies

### NATS JetStream (Primary)
- **Use Case**: Development workflows, real-time updates, low-latency (<5ms)
- **Throughput**: Up to 10k messages/second
- **Features**: At-least-once delivery, stream persistence, consumer groups

### Apache Kafka (Scale Option)
- **Use Case**: High-throughput scenarios (>10k messages/second)
- **Features**: Partitioned topics, horizontal scaling, long-term retention

## Event Schema

### Base Event Structure

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "type": "sc.plugin.executed",
  "source": "superclaude-analyzer",
  "timestamp": "2025-03-26T10:30:00Z",
  "version": "1.0",
  "data": {
    "plugin": "sc-analyzer",
    "command": "analyze",
    "status": "completed",
    "duration_ms": 2500,
    "artifacts": ["analysis-report.md"]
  },
  "metadata": {
    "correlation_id": "req-123",
    "user_id": "user-456",
    "session_id": "session-789"
  }
}
```

### Event Types and Schemas

#### Plugin Lifecycle Events

**Plugin Started**
```json
{
  "type": "sc.plugin.started",
  "source": "plugin-manager",
  "data": {
    "plugin_name": "sc-analyzer",
    "plugin_version": "1.2.0",
    "pid": 12345,
    "capabilities": ["filesystem", "events"],
    "memory_limit_mb": 512
  }
}
```

**Plugin Stopped**
```json
{
  "type": "sc.plugin.stopped",
  "source": "plugin-manager", 
  "data": {
    "plugin_name": "sc-analyzer",
    "reason": "requested",
    "exit_code": 0,
    "uptime_seconds": 3600
  }
}
```

**Plugin Error**
```json
{
  "type": "sc.plugin.error",
  "source": "sc-analyzer",
  "data": {
    "error_code": "RESOURCE_LIMIT",
    "error_message": "Memory limit exceeded",
    "stack_trace": "...",
    "request_id": "req-123"
  }
}
```

#### Command Execution Events

**Command Started**
```json
{
  "type": "sc.command.started",
  "source": "opencode-cli",
  "data": {
    "command": "/analyze",
    "arguments": {
      "path": "./src",
      "focus": "performance"
    },
    "user_id": "user-456",
    "working_directory": "/workspace/project"
  }
}
```

**Command Completed**
```json
{
  "type": "sc.command.completed",
  "source": "sc-analyzer",
  "data": {
    "command": "/analyze",
    "status": "success",
    "duration_ms": 15000,
    "artifacts_created": 3,
    "files_analyzed": 142,
    "issues_found": 8
  }
}
```

**Command Progress**
```json
{
  "type": "sc.command.progress",
  "source": "sc-analyzer",
  "data": {
    "command": "/analyze",
    "progress_percent": 45,
    "current_stage": "analyzing_components",
    "files_processed": 64,
    "total_files": 142,
    "eta_seconds": 20
  }
}
```

#### Configuration Events

**Configuration Changed**
```json
{
  "type": "sc.config.changed",
  "source": "config-manager",
  "data": {
    "config_type": "plugin",
    "plugin_name": "sc-analyzer",
    "changes": [
      {
        "key": "max_memory_mb",
        "old_value": "256",
        "new_value": "512"
      }
    ],
    "requires_restart": true
  }
}
```

**Configuration Validated**
```json
{
  "type": "sc.config.validated",
  "source": "config-manager",
  "data": {
    "config_type": "system",
    "validation_result": "success",
    "warnings": [],
    "errors": []
  }
}
```

#### File System Events

**File Changed**
```json
{
  "type": "sc.file.changed",
  "source": "file-watcher",
  "data": {
    "file_path": "./src/components/Button.tsx",
    "change_type": "modified",
    "size_bytes": 2048,
    "checksum": "sha256:abc123...",
    "git_status": "modified"
  }
}
```

**Project Structure Changed**
```json
{
  "type": "sc.project.structure_changed",
  "source": "file-watcher",
  "data": {
    "change_type": "directory_added",
    "path": "./src/features/auth",
    "files_added": 5,
    "framework_detected": "react"
  }
}
```

#### Integration Events

**MCP Server Connected**
```json
{
  "type": "sc.mcp.connected",
  "source": "mcp-client",
  "data": {
    "server_name": "superclaude-bridge",
    "protocol_version": "2025-03-26",
    "capabilities": ["tools", "prompts", "resources"],
    "connection_type": "stdio"
  }
}
```

**External Tool Invoked**
```json
{
  "type": "sc.external.tool_invoked",
  "source": "tool-runner",
  "data": {
    "tool_name": "eslint",
    "command": "eslint src/ --fix",
    "exit_code": 0,
    "output_lines": 23,
    "files_fixed": 3
  }
}
```

## Topic Structure

### NATS JetStream Topics

```
superclaudebridge/
├── plugin/
│   ├── lifecycle/
│   │   ├── started
│   │   ├── stopped
│   │   └── error
│   ├── execution/
│   │   ├── started
│   │   ├── progress
│   │   └── completed
│   └── health/
│       ├── status
│       └── metrics
├── system/
│   ├── config/
│   │   ├── changed
│   │   └── validated
│   ├── file/
│   │   ├── changed
│   │   └── project_structure
│   └── integration/
│       ├── mcp_connected
│       └── external_tool
└── telemetry/
    ├── usage
    ├── performance
    └── errors
```

### Kafka Topics (High Throughput)

```yaml
topics:
  - name: "sc-plugin-events"
    partitions: 3
    replication: 3
    retention: "24h"
    
  - name: "sc-system-events"
    partitions: 1
    replication: 3
    retention: "7d"
    
  - name: "sc-telemetry"
    partitions: 6
    replication: 3
    retention: "30d"
```

## Event Publishing

### Go Publisher Example

```go
package main

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/nats-io/nats.go"
)

type EventPublisher struct {
    nc *nats.Conn
    js nats.JetStreamContext
}

func NewEventPublisher(natsURL string) (*EventPublisher, error) {
    nc, err := nats.Connect(natsURL)
    if err != nil {
        return nil, err
    }
    
    js, err := nc.JetStream()
    if err != nil {
        return nil, err
    }
    
    return &EventPublisher{nc: nc, js: js}, nil
}

func (ep *EventPublisher) PublishPluginEvent(eventType string, source string, data interface{}) error {
    event := Event{
        ID:        generateEventID(),
        Type:      eventType,
        Source:    source,
        Timestamp: time.Now(),
        Version:   "1.0",
        Data:      data,
    }
    
    eventBytes, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    subject := fmt.Sprintf("superclaudebridge.plugin.%s", strings.ReplaceAll(eventType, ".", "_"))
    _, err = ep.js.Publish(subject, eventBytes)
    return err
}

// Usage example
func (p *MyPlugin) EmitStartedEvent() error {
    data := map[string]interface{}{
        "plugin_name":    "sc-analyzer",
        "plugin_version": "1.2.0",
        "capabilities":   []string{"filesystem", "events"},
    }
    
    return p.eventPublisher.PublishPluginEvent("sc.plugin.started", "sc-analyzer", data)
}
```

### TypeScript Publisher Example

```typescript
import { connect, JetStreamClient } from 'nats';

export class EventPublisher {
    private nc: any;
    private js: JetStreamClient;

    async connect(natsUrl: string): Promise<void> {
        this.nc = await connect({ servers: natsUrl });
        this.js = this.nc.jetstream();
    }

    async publishCommandEvent(eventType: string, source: string, data: any): Promise<void> {
        const event = {
            id: this.generateEventID(),
            type: eventType,
            source: source,
            timestamp: new Date().toISOString(),
            version: '1.0',
            data: data
        };

        const subject = `superclaudebridge.command.${eventType.replace(/\./g, '_')}`;
        await this.js.publish(subject, JSON.stringify(event));
    }

    private generateEventID(): string {
        return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    }
}

// Usage example
const publisher = new EventPublisher();
await publisher.connect('nats://localhost:4222');

await publisher.publishCommandEvent('sc.command.started', 'opencode-cli', {
    command: '/analyze',
    arguments: { path: './src', focus: 'performance' },
    user_id: 'user-456'
});
```

## Event Subscription

### Go Subscriber Example

```go
type EventSubscriber struct {
    nc *nats.Conn
    js nats.JetStreamContext
}

func (es *EventSubscriber) SubscribeToPluginEvents(handler func(*Event) error) error {
    // Create durable consumer
    consumer, err := es.js.PullSubscribe(
        "superclaudebridge.plugin.*",
        "plugin-event-processor",
        nats.DeliverAll(),
        nats.AckExplicit(),
    )
    if err != nil {
        return err
    }
    
    // Process messages
    go func() {
        for {
            msgs, err := consumer.Fetch(10, nats.MaxWait(5*time.Second))
            if err != nil {
                log.Printf("Error fetching messages: %v", err)
                continue
            }
            
            for _, msg := range msgs {
                var event Event
                if err := json.Unmarshal(msg.Data, &event); err != nil {
                    log.Printf("Error unmarshaling event: %v", err)
                    msg.Nak()
                    continue
                }
                
                if err := handler(&event); err != nil {
                    log.Printf("Error handling event: %v", err)
                    msg.Nak()
                    continue
                }
                
                msg.Ack()
            }
        }
    }()
    
    return nil
}

// Event handler example
func handlePluginEvent(event *Event) error {
    switch event.Type {
    case "sc.plugin.started":
        return handlePluginStarted(event)
    case "sc.plugin.error":
        return handlePluginError(event)
    default:
        log.Printf("Unknown event type: %s", event.Type)
        return nil
    }
}
```

### Queue Groups and Load Balancing

```go
// Multiple consumers in same queue group for load balancing
func (es *EventSubscriber) SubscribeWithQueueGroup(subject, queueGroup string, handler func(*Event) error) error {
    subscription, err := es.nc.QueueSubscribe(subject, queueGroup, func(msg *nats.Msg) {
        var event Event
        if err := json.Unmarshal(msg.Data, &event); err != nil {
            log.Printf("Error unmarshaling event: %v", err)
            return
        }
        
        if err := handler(&event); err != nil {
            log.Printf("Error handling event: %v", err)
        }
    })
    
    return err
}

// Usage: 3 workers processing events in parallel
es.SubscribeWithQueueGroup("superclaudebridge.command.*", "command-processors", handleCommandEvent)
```

## Event Routing and Filtering

### Content-Based Routing

```go
type EventRouter struct {
    routes map[string][]RouteHandler
}

type RouteHandler struct {
    Filter  func(*Event) bool
    Handler func(*Event) error
}

func (er *EventRouter) AddRoute(eventType string, filter func(*Event) bool, handler func(*Event) error) {
    if er.routes[eventType] == nil {
        er.routes[eventType] = []RouteHandler{}
    }
    
    er.routes[eventType] = append(er.routes[eventType], RouteHandler{
        Filter:  filter,
        Handler: handler,
    })
}

func (er *EventRouter) RouteEvent(event *Event) error {
    handlers, exists := er.routes[event.Type]
    if !exists {
        return nil // No handlers for this event type
    }
    
    for _, handler := range handlers {
        if handler.Filter == nil || handler.Filter(event) {
            if err := handler.Handler(event); err != nil {
                log.Printf("Handler error for event %s: %v", event.ID, err)
            }
        }
    }
    
    return nil
}

// Usage example
router := &EventRouter{routes: make(map[string][]RouteHandler)}

// Route only error events from specific plugin
router.AddRoute("sc.plugin.error", func(event *Event) bool {
    return event.Source == "sc-analyzer"
}, handleAnalyzerErrors)

// Route all command completion events
router.AddRoute("sc.command.completed", nil, updateMetrics)
```

## Error Handling and Dead Letter Queues

### NATS JetStream Configuration

```go
// Create stream with dead letter queue
streamConfig := &nats.StreamConfig{
    Name:      "superclaudebridge-events",
    Subjects:  []string{"superclaudebridge.*"},
    Storage:   nats.FileStorage,
    Retention: nats.LimitsPolicy,
    MaxAge:    24 * time.Hour,
    
    // Dead letter queue configuration
    MaxDelivery: 3,
    DeadLetter: &nats.DeadLetterPolicy{
        Stream: "superclaudebridge-dlq",
    },
}

stream, err := js.AddStream(streamConfig)
```

### Error Recovery

```go
func (es *EventSubscriber) ProcessFailedEvents() error {
    // Subscribe to dead letter queue
    consumer, err := es.js.PullSubscribe(
        "superclaudebridge-dlq.*",
        "dlq-processor",
    )
    if err != nil {
        return err
    }
    
    for {
        msgs, err := consumer.Fetch(10, nats.MaxWait(30*time.Second))
        if err != nil {
            continue
        }
        
        for _, msg := range msgs {
            var event Event
            if err := json.Unmarshal(msg.Data, &event); err != nil {
                log.Printf("Cannot parse DLQ event: %v", err)
                msg.Ack() // Acknowledge to remove from DLQ
                continue
            }
            
            // Try to recover or log for manual intervention
            if err := es.recoverFailedEvent(&event); err != nil {
                log.Printf("Failed to recover event %s: %v", event.ID, err)
                // Could republish to another queue for manual processing
            }
            
            msg.Ack()
        }
    }
}
```

## Event Persistence and Replay

### JetStream Stream Configuration

```yaml
streams:
  - name: "superclaudebridge-events"
    subjects: ["superclaudebridge.*"]
    storage: "file"
    retention: "limits"
    max_age: "24h"
    max_bytes: "1GB"
    max_msgs: 1000000
    
  - name: "superclaudebridge-audit"
    subjects: ["superclaudebridge.audit.*"]
    storage: "file"
    retention: "limits"
    max_age: "30d"
    max_bytes: "10GB"
```

### Event Replay

```go
func (es *EventSubscriber) ReplayEvents(fromTime time.Time, eventTypes []string) error {
    for _, eventType := range eventTypes {
        subject := fmt.Sprintf("superclaudebridge.%s", eventType)
        
        consumer, err := es.js.PullSubscribe(
            subject,
            "replay-consumer",
            nats.DeliverByStartTime(fromTime),
            nats.ReplayInstant(),
        )
        if err != nil {
            return err
        }
        
        // Process historical events
        go es.processReplayedEvents(consumer)
    }
    
    return nil
}
```

## Monitoring and Observability

### Event Metrics

```go
type EventMetrics struct {
    eventsPublished   prometheus.Counter
    eventsConsumed    prometheus.Counter
    processingLatency prometheus.Histogram
    errorRate         prometheus.Counter
}

func (em *EventMetrics) RecordEventPublished(eventType string) {
    em.eventsPublished.WithLabelValues(eventType).Inc()
}

func (em *EventMetrics) RecordProcessingTime(eventType string, duration time.Duration) {
    em.processingLatency.WithLabelValues(eventType).Observe(duration.Seconds())
}
```

### Health Checks

```go
func (es *EventSystem) HealthCheck() *HealthStatus {
    status := &HealthStatus{
        Status: "healthy",
        Checks: map[string]interface{}{},
    }
    
    // Check NATS connection
    if !es.nc.IsConnected() {
        status.Status = "unhealthy"
        status.Checks["nats_connection"] = "disconnected"
    } else {
        status.Checks["nats_connection"] = "connected"
    }
    
    // Check JetStream
    if _, err := es.js.AccountInfo(); err != nil {
        status.Status = "degraded"
        status.Checks["jetstream"] = fmt.Sprintf("error: %v", err)
    } else {
        status.Checks["jetstream"] = "available"
    }
    
    // Check consumer lag
    consumers, _ := es.js.Consumers("superclaudebridge-events")
    for consumer := range consumers {
        info, _ := consumer.Info()
        if info.NumPending > 1000 {
            status.Status = "degraded"
            status.Checks[fmt.Sprintf("consumer_%s_lag", info.Name)] = info.NumPending
        }
    }
    
    return status
}
```

## Configuration

### NATS JetStream Configuration

```yaml
# nats-server.conf
jetstream {
    store_dir: "/data/jetstream"
    max_memory_store: 1GB
    max_file_store: 10GB
}

accounts {
    superclaudebridge {
        jetstream: enabled
        users: [
            {user: "sc-publisher", password: "pub-secret"}
            {user: "sc-consumer", password: "con-secret"}
        ]
    }
}
```

### Event System Configuration

```json
{
  "event_system": {
    "provider": "nats",
    "nats": {
      "url": "nats://localhost:4222",
      "credentials": "/etc/superclaudebridge/nats.creds",
      "jetstream": {
        "enabled": true,
        "storage": "file",
        "max_age": "24h"
      }
    },
    "kafka": {
      "brokers": ["localhost:9092"],
      "security": {
        "protocol": "SASL_SSL",
        "mechanism": "PLAIN"
      }
    },
    "topics": {
      "plugin_events": "superclaudebridge.plugin",
      "system_events": "superclaudebridge.system",
      "telemetry": "superclaudebridge.telemetry"
    },
    "consumers": {
      "plugin_processor": {
        "queue_group": "plugin-processors",
        "max_concurrent": 5,
        "ack_wait": "30s"
      }
    }
  }
}
```

## Integration Examples

### Plugin Event Integration

```go
// Plugin publishes execution progress
func (p *AnalyzerPlugin) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
    // Emit start event
    p.eventPublisher.PublishPluginEvent("sc.command.started", p.name, map[string]interface{}{
        "command": req.Command,
        "arguments": req.Arguments,
    })
    
    // Process with progress updates
    totalFiles := len(req.Context.FilePaths)
    for i, file := range req.Context.FilePaths {
        // Analyze file...
        
        // Emit progress
        progress := float64(i+1) / float64(totalFiles) * 100
        p.eventPublisher.PublishPluginEvent("sc.command.progress", p.name, map[string]interface{}{
            "command": req.Command,
            "progress_percent": progress,
            "current_file": file,
        })
    }
    
    // Emit completion
    p.eventPublisher.PublishPluginEvent("sc.command.completed", p.name, map[string]interface{}{
        "command": req.Command,
        "status": "success",
        "files_analyzed": totalFiles,
    })
    
    return &ExecuteResponse{...}, nil
}
```

### Configuration Hot Reload

```go
// Configuration service listens for config changes and notifies plugins
func (cs *ConfigService) handleConfigChanged(event *Event) error {
    data := event.Data.(map[string]interface{})
    pluginName := data["plugin_name"].(string)
    
    // Notify specific plugin about config change
    return cs.eventPublisher.PublishPluginEvent("sc.config.reload", "config-service", map[string]interface{}{
        "plugin_name": pluginName,
        "config_path": data["config_path"],
        "requires_restart": data["requires_restart"],
    })
}

// Plugin handles config reload
func (p *MyPlugin) handleConfigReload(event *Event) error {
    data := event.Data.(map[string]interface{})
    configPath := data["config_path"].(string)
    
    // Reload configuration
    newConfig, err := p.loadConfig(configPath)
    if err != nil {
        return err
    }
    
    p.config = newConfig
    
    // Emit acknowledgment
    return p.eventPublisher.PublishPluginEvent("sc.config.reloaded", p.name, map[string]interface{}{
        "plugin_name": p.name,
        "config_version": newConfig.Version,
    })
}
```

## Testing

### Event Testing Framework

```go
type EventTestSuite struct {
    eventBus    *EventBus
    mockNATS    *MockNATSConnection
    testEvents  []Event
}

func (ets *EventTestSuite) SetupTest() {
    ets.mockNATS = NewMockNATSConnection()
    ets.eventBus = NewEventBus(ets.mockNATS)
    ets.testEvents = []Event{}
}

func (ets *EventTestSuite) TestPluginEventFlow() {
    // Setup event capture
    ets.eventBus.Subscribe("sc.plugin.*", func(event *Event) error {
        ets.testEvents = append(ets.testEvents, *event)
        return nil
    })
    
    // Trigger plugin events
    plugin := &TestPlugin{eventBus: ets.eventBus}
    plugin.Start()
    plugin.Execute("test-command", map[string]string{})
    plugin.Stop()
    
    // Assert expected events
    assert.Len(ets.T(), ets.testEvents, 3)
    assert.Equal(ets.T(), "sc.plugin.started", ets.testEvents[0].Type)
    assert.Equal(ets.T(), "sc.command.completed", ets.testEvents[1].Type)
    assert.Equal(ets.T(), "sc.plugin.stopped", ets.testEvents[2].Type)
}
```

## References

- [NATS JetStream Documentation](https://docs.nats.io/jetstream)
- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
- [Event-Driven Architecture Patterns](../EDA.md)
- [Plugin API Specification](./plugin-api.md)
- [MCP Protocol Integration](./mcp-api.md)