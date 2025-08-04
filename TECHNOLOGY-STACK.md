# Technology Stack Reference

## Overview

This document provides comprehensive details about all technologies, frameworks, and tools used across the five OpenCode-SuperClaude integration plans. Each technology includes version requirements, use cases, setup instructions, and plan-specific usage.

---

## Core Programming Languages

### Go
**Versions**: Go 1.21+ (required), Go 1.22+ (recommended)

**Primary Usage**: 
- OpenCode core implementation
- Plugin host systems
- MCP server implementations
- Event bus management

**Plans Using Go**: All plans (1, 2, 3, 4, 5)

**Key Features for Our Use Case**:
- Excellent concurrency with goroutines
- Fast compilation and deployment
- Strong CLI application ecosystem
- Cross-platform binary distribution
- Mature plugin frameworks

**Installation**:
```bash
# Using official installer
curl -L https://go.dev/dl/go1.22.0.linux-amd64.tar.gz | sudo tar -C /usr/local -xzf -

# Using package managers
# macOS
brew install go

# Ubuntu/Debian
sudo apt install golang-go

# Verify installation
go version
```

**Required Dependencies by Plan**:
```go
// Plan 1, 3, 4: Plugin systems
"github.com/hashicorp/go-plugin" v1.6.0+

// Plan 3, 4: Event systems  
"github.com/nats-io/nats.go" v1.31.0+
"github.com/segmentio/kafka-go" v0.4.47+

// All plans: gRPC communication
"google.golang.org/grpc" v1.59.0+
"google.golang.org/protobuf" v1.31.0+

// All plans: Configuration
"github.com/spf13/viper" v1.17.0+
"github.com/spf13/cobra" v1.8.0+

// Plan 3, 4, 5: Observability
"go.opentelemetry.io/otel" v1.21.0+
"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc" v0.46.0+
```

### TypeScript
**Versions**: TypeScript 5.2+ (required for advanced features)

**Primary Usage**:
- SuperClaude prompt management
- File watching and live-reload
- UI component integration
- Template processing

**Plans Using TypeScript**: Plans 2, 3, 5 (primary), Plan 4 (optional)

**Installation & Setup**:
```bash
# Global installation
npm install -g typescript@^5.2.0

# Project-specific
npm install --save-dev typescript@^5.2.0 @types/node

# Initialize TypeScript project
npx tsc --init
```

**Required Dependencies by Plan**:
```json
{
  "dependencies": {
    "chokidar": "^3.5.3",
    "gray-matter": "^4.0.3", 
    "marked": "^9.1.2",
    "handlebars": "^4.7.8",
    "zod": "^3.22.4",
    "rxjs": "^7.8.1",
    "winston": "^3.11.0"
  },
  "devDependencies": {
    "@types/marked": "^6.0.0",
    "@types/handlebars": "^4.1.0",
    "vitest": "^1.0.0",
    "typescript": "^5.2.0"
  }
}
```

**TypeScript Configuration** (`tsconfig.json`):
```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "lib": ["ES2022"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "paths": {
      "@superclaude/*": ["./src/superclaude/*"],
      "@opencode/*": ["./src/opencode/*"]
    }
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist", "tests"]
}
```

---

## Communication Protocols

### gRPC
**Version**: gRPC 1.59.0+ (Go), @grpc/grpc-js 1.9.0+ (Node.js)

**Usage**: Inter-process communication for plugin systems

**Plans**: 1, 3, 4 (primary), 2, 5 (optional for performance)

**Protocol Buffer Setup**:
```bash
# Install protoc compiler
# macOS
brew install protobuf

# Ubuntu/Debian  
sudo apt install protobuf-compiler

# Install language-specific plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

**Example Service Definition** (`plugin.proto`):
```protobuf
syntax = "proto3";

package opencode.plugin.v1;

option go_package = "github.com/opencode/plugin/v1;pluginv1";

service PluginService {
  rpc Execute(ExecuteRequest) returns (ExecuteResponse);
  rpc GetMetadata(GetMetadataRequest) returns (GetMetadataResponse);
  rpc Health(HealthRequest) returns (HealthResponse);
}

message ExecuteRequest {
  string command = 1;
  map<string, string> context = 2;
  repeated string args = 3;
}

message ExecuteResponse {
  string result = 1;
  int32 exit_code = 2;
  string error = 3;
}
```

**Code Generation**:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       plugin.proto
```

### JSON-RPC
**Version**: JSON-RPC 2.0 specification

**Usage**: MCP communication, fallback for gRPC

**Plans**: All plans (MCP standard requirement)

**Implementation Example**:
```go
type JSONRPCRequest struct {
    JSONRPC string                 `json:"jsonrpc"`
    Method  string                 `json:"method"`
    Params  map[string]interface{} `json:"params,omitempty"`
    ID      interface{}            `json:"id"`
}

type JSONRPCResponse struct {
    JSONRPC string      `json:"jsonrpc"`
    Result  interface{} `json:"result,omitempty"`
    Error   *RPCError   `json:"error,omitempty"`
    ID      interface{} `json:"id"`
}
```

---

## Event Systems

### NATS JetStream
**Version**: NATS Server 2.10.0+, nats.go client 1.31.0+

**Usage**: Primary event bus for low-to-medium throughput scenarios (<10k msg/s)

**Plans**: 3, 4 (primary), 1 (optional)

**Installation**:
```bash
# Server installation
# macOS
brew install nats-server

# Ubuntu/Debian
curl -L https://github.com/nats-io/nats-server/releases/download/v2.10.7/nats-server-v2.10.7-linux-amd64.zip -o nats-server.zip
unzip nats-server.zip
sudo mv nats-server-v2.10.7-linux-amd64/nats-server /usr/local/bin/

# Docker
docker run -p 4222:4222 -p 8222:8222 nats:latest --jetstream
```

**Configuration** (`nats-server.conf`):
```
port: 4222
http_port: 8222

jetstream {
  store_dir: "/data/jetstream"
  max_memory: 1G
  max_file: 10G
}

accounts {
  opencode: {
    jetstream: enabled
    users: [
      {user: opencode, password: secret}
    ]
  }
}
```

**Go Client Setup**:
```go
import "github.com/nats-io/nats.go"

nc, err := nats.Connect(
    "nats://localhost:4222",
    nats.UserInfo("opencode", "secret"),
)
if err != nil {
    log.Fatal(err)
}
defer nc.Close()

js, err := nc.JetStream()
if err != nil {
    log.Fatal(err)
}

// Create stream
stream := &nats.StreamConfig{
    Name:     "superclaude-events",
    Subjects: []string{"sc.*", "oc.*", "plugin.*"},
    Storage:  nats.FileStorage,
    MaxAge:   24 * time.Hour,
}

_, err = js.AddStream(stream)
```

### Apache Kafka
**Version**: Kafka 3.6.0+, kafka-go 0.4.47+ (Go client)

**Usage**: High-throughput event streaming (>10k msg/s) or existing Kafka infrastructure

**Plans**: 3, 4 (alternative to NATS)

**Installation**:
```bash
# Using Docker Compose
version: '3.8'
services:
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: 'zookeeper:2181'
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_INTERNAL:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092,PLAINTEXT_INTERNAL://broker:29092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
```

**Go Client Setup**:
```go
import "github.com/segmentio/kafka-go"

writer := &kafka.Writer{
    Addr:     kafka.TCP("localhost:9092"),
    Topic:    "superclaude-events",
    Balancer: &kafka.LeastBytes{},
}

reader := kafka.NewReader(kafka.ReaderConfig{
    Brokers: []string{"localhost:9092"},
    Topic:   "superclaude-events",
    GroupID: "opencode-consumer",
})
```

---

## Plugin Systems

### HashiCorp go-plugin
**Version**: v1.6.0+ (required for latest features)

**Usage**: Process-isolated plugin architecture

**Plans**: 1, 3, 4 (core component)

**Installation**:
```bash
go get github.com/hashicorp/go-plugin@latest
```

**Plugin Interface Definition**:
```go
// Shared interface between host and plugin
type CommandPlugin interface {
    Name() string
    Version() string
    Execute(context map[string]string) (string, error) 
    Metadata() PluginMetadata
}

type PluginMetadata struct {
    Name        string            `json:"name"`
    Version     string            `json:"version"`
    Description string            `json:"description"`
    Author      string            `json:"author"`
    Capabilities []string         `json:"capabilities"`
    Config      map[string]string `json:"config"`
}
```

**Plugin Implementation**:
```go
// Plugin side (separate binary)
type MyPlugin struct{}

func (p *MyPlugin) Execute(context map[string]string) (string, error) {
    // Plugin logic here
    return "result", nil
}

func main() {
    plugin.Serve(&plugin.ServeConfig{
        HandshakeConfig: plugin.HandshakeConfig{
            ProtocolVersion:  1,
            MagicCookieKey:   "OPENCODE_PLUGIN",
            MagicCookieValue: "opencode",
        },
        Plugins: map[string]plugin.Plugin{
            "command": &CommandPluginImpl{
                Impl: &MyPlugin{},
            },
        },
        GRPCServer: plugin.DefaultGRPCServer,
    })
}
```

**Host Implementation**:
```go
// Host side (OpenCode)
client := plugin.NewClient(&plugin.ClientConfig{
    HandshakeConfig: handshakeConfig,
    Plugins: pluginMap,
    Cmd: exec.Command("./my-plugin"),
    AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
})
defer client.Kill()

rpcClient, err := client.Client()
if err != nil {
    return err
}

raw, err := rpcClient.Dispense("command")
if err != nil {
    return err
}

commandPlugin := raw.(CommandPlugin)
result, err := commandPlugin.Execute(context)
```

---

## File System & Watching

### Chokidar (Node.js/TypeScript)
**Version**: 3.5.3+

**Usage**: File watching for live-reload functionality

**Plans**: 2, 3, 5 (TypeScript components)

**Installation**:
```bash
npm install chokidar @types/chokidar
```

**Implementation**:
```typescript
import chokidar from 'chokidar';

const watcher = chokidar.watch('prompts/**/*.md', {
  ignored: /(^|[\/\\])\../,
  persistent: true,
  ignoreInitial: false,
  awaitWriteFinish: {
    stabilityThreshold: 100,
    pollInterval: 50
  }
});

watcher
  .on('add', path => console.log(`File ${path} has been added`))
  .on('change', path => console.log(`File ${path} has been changed`))
  .on('unlink', path => console.log(`File ${path} has been removed`));
```

### fsnotify (Go)
**Version**: v1.7.0+

**Usage**: File watching in Go components

**Plans**: 1, 4 (Go-based file watching)

```go
import "github.com/fsnotify/fsnotify"

watcher, err := fsnotify.NewWatcher()
if err != nil {
    log.Fatal(err)
}
defer watcher.Close()

go func() {
    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }
            if event.Has(fsnotify.Write) {
                log.Println("modified file:", event.Name)
                // Reload logic here
            }
        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            log.Println("error:", err)
        }
    }
}()

err = watcher.Add("/path/to/watch")
```

---

## Configuration & Parsing

### Viper (Go)
**Version**: v1.17.0+

**Usage**: Configuration management across all Go components

**Plans**: 1, 3, 4 (Go configuration)

```go
import "github.com/spf13/viper"

func initConfig() {
    viper.SetConfigName("opencode")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("$HOME/.opencode")
    viper.AddConfigPath(".")
    
    // Set defaults
    viper.SetDefault("plugins.directory", "~/.opencode/plugins")
    viper.SetDefault("events.broker", "nats")
    viper.SetDefault("log.level", "info")
    
    // Environment variables
    viper.SetEnvPrefix("OPENCODE")
    viper.AutomaticEnv()
    
    if err := viper.ReadInConfig(); err != nil {
        log.Printf("Config file not found: %v", err)
    }
}
```

### Gray-matter (TypeScript)
**Version**: 4.0.3+

**Usage**: YAML frontmatter parsing in SuperClaude templates

**Plans**: 2, 3, 5

```typescript
import matter from 'gray-matter';

const content = `
---
name: "explain"
version: "1.0.0"
description: "Explain code"
requires:
  selectedCode: true
---

# Code Explanation
Explain this code: {{selectedCode}}
`;

const { data, content: markdownContent } = matter(content);
console.log(data.name); // "explain"
console.log(markdownContent); // "# Code Explanation..."
```

### Zod (TypeScript)
**Version**: 3.22.4+

**Usage**: Schema validation and type safety

**Plans**: 2, 3, 5

```typescript
import { z } from 'zod';

const PromptMetadataSchema = z.object({
  name: z.string(),
  version: z.string().default('1.0.0'),
  description: z.string(),
  category: z.enum(['code', 'analysis', 'generation']),
  requires: z.object({
    selectedCode: z.boolean().default(false),
    filePath: z.boolean().default(false),
  }).default({}),
});

type PromptMetadata = z.infer<typeof PromptMetadataSchema>;

// Validation
const metadata = PromptMetadataSchema.parse(yamlData);
```

---

## Template Processing

### Handlebars (TypeScript)
**Version**: 4.7.8+

**Usage**: Template rendering for SuperClaude prompts

**Plans**: 2, 3, 5

```typescript
import Handlebars from 'handlebars';

// Register custom helpers
Handlebars.registerHelper('codeBlock', (code: string, language = '') => {
  return new Handlebars.SafeString(`\`\`\`${language}\n${code}\n\`\`\``);
});

Handlebars.registerHelper('truncate', (str: string, length: number) => {
  return str.length > length ? str.substring(0, length) + '...' : str;
});

// Compile and render
const template = Handlebars.compile('Explain this code: {{codeBlock selectedCode "javascript"}}');
const result = template({ selectedCode: 'console.log("hello");' });
```

### text/template (Go)
**Usage**: Template processing in Go components

```go
import "text/template"

const promptTemplate = `
Analysis for {{.FileName}}:
{{range .Issues}}
- {{.Severity}}: {{.Description}}
{{end}}
`

tmpl, err := template.New("prompt").Parse(promptTemplate)
if err != nil {
    log.Fatal(err)
}

data := struct {
    FileName string
    Issues   []Issue
}{
    FileName: "main.go",
    Issues:   issues,
}

err = tmpl.Execute(os.Stdout, data)
```

---

## Observability & Monitoring

### OpenTelemetry
**Versions**: 
- Go: v1.21.0+
- Node.js: @opentelemetry/api v1.7.0+

**Usage**: Distributed tracing and metrics collection

**Plans**: 3, 4, 5 (production deployments)

**Go Setup**:
```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/trace"
)

func initTracing() {
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint())
    if err != nil {
        log.Fatal(err)
    }
    
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String("opencode"),
        )),
    )
    
    otel.SetTracerProvider(tp)
}

// Usage in code
func executePlugin(ctx context.Context, name string) error {
    tracer := otel.Tracer("plugin-manager")
    ctx, span := tracer.Start(ctx, "execute-plugin")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("plugin.name", name),
        attribute.String("plugin.version", version),
    )
    
    // Plugin execution logic
    return nil
}
```

### Prometheus (Metrics)
**Usage**: Metrics collection and alerting

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    pluginExecutions = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "opencode_plugin_executions_total",
            Help: "Total number of plugin executions",
        },
        []string{"plugin", "status"},
    )
    
    executionDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "opencode_plugin_duration_seconds",
            Help: "Plugin execution duration",
        },
        []string{"plugin"},
    )
)

// Usage
timer := prometheus.NewTimer(executionDuration.WithLabelValues(pluginName))
defer timer.ObserveDuration()

pluginExecutions.WithLabelValues(pluginName, "success").Inc()
```

---

## Testing Frameworks

### Vitest (TypeScript)
**Version**: 1.0.0+

**Usage**: Fast unit testing for TypeScript components

**Plans**: 2, 3, 5

**Configuration** (`vitest.config.ts`):
```typescript
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    globals: true,
    environment: 'node',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: ['node_modules/', 'dist/', 'tests/'],
    },
  },
  resolve: {
    alias: {
      '@': '/src',
    },
  },
});
```

**Example Test**:
```typescript
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { PromptBroker } from '../src/broker/PromptBroker';

describe('PromptBroker', () => {
  let broker: PromptBroker;

  beforeEach(() => {
    broker = new PromptBroker({
      promptsDirectory: './test-prompts',
      watchEnabled: false,
    });
  });

  it('should execute prompts with context', async () => {
    const result = await broker.executePrompt('explain', {
      selectedCode: 'console.log("test");',
      filePath: 'test.js',
    });

    expect(result).toContain('console.log');
    expect(result).toContain('test.js');
  });

  it('should handle missing templates', async () => {
    await expect(
      broker.executePrompt('nonexistent', {})
    ).rejects.toThrow('Template \'nonexistent\' not found');
  });
});
```

### Go Testing (Standard Library)
**Usage**: Unit and integration testing for Go components

```go
func TestPluginManager(t *testing.T) {
    manager := NewPluginManager(Config{
        PluginDir: "./test-plugins",
    })
    
    err := manager.LoadPlugin("test-plugin")
    assert.NoError(t, err)
    
    plugin := manager.GetPlugin("test-plugin")
    assert.NotNil(t, plugin)
    
    result, err := plugin.Execute(map[string]string{
        "input": "test data",
    })
    assert.NoError(t, err)
    assert.Equal(t, "expected result", result)
}

func TestPluginIsolation(t *testing.T) {
    // Test that plugin crashes don't affect host
    manager := NewPluginManager(Config{})
    
    // Load plugin that will crash
    err := manager.LoadPlugin("crash-plugin")
    assert.NoError(t, err)
    
    // Execute should return error, not crash host
    _, err = manager.ExecutePlugin("crash-plugin", map[string]string{})
    assert.Error(t, err)
    
    // Host should still be functional
    assert.True(t, manager.IsHealthy())
}
```

---

## Containerization & Deployment

### Docker
**Version**: Docker 24.0+, Docker Compose 2.21+

**Usage**: Containerization for deployment and development environments

**Plans**: All plans (deployment), 3, 4 (development)

**Dockerfile Example** (Go MCP Server):
```dockerfile
# Multi-stage build
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o mcp-server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/mcp-server .
COPY --from=builder /app/configs/ ./configs/

EXPOSE 8080
CMD ["./mcp-server"]
```

**Docker Compose** (Development Environment):
```yaml
version: '3.8'

services:
  nats:
    image: nats:2.10-alpine
    ports:
      - "4222:4222"
      - "8222:8222"
    command: ["-js", "-m", "8222"]
    volumes:
      - nats_data:/data
    
  mcp-server:
    build: 
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    environment:
      - NATS_URL=nats://nats:4222
      - LOG_LEVEL=debug
    depends_on:
      - nats
    volumes:
      - ./configs:/app/configs
      - ./plugins:/app/plugins
    
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    
  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"
      - "14268:14268"
    environment:
      - COLLECTOR_OTLP_ENABLED=true

volumes:
  nats_data:
```

---

## Development Environment Setup

### Prerequisites by Plan

**Plan 1 (Plugin Ecosystem)**:
```bash
# Core requirements
go version # 1.21+
docker --version # 24.0+
protoc --version # 3.21+

# Install Go dependencies
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Optional: Development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/air-verse/air@latest # Live reload
```

**Plan 2 (SuperClaude Integration)**:
```bash
# Core requirements
node --version # 18+
npm --version # 9+
tsc --version # 5.2+

# Install TypeScript globally
npm install -g typescript@^5.2.0

# Project setup
npm init -y
npm install chokidar gray-matter marked handlebars zod winston
npm install -D @types/node vitest
```

**Plan 3 (MCP Bridge + Plugin + EDA)**:
```bash
# Requires both Go and TypeScript environments
# Plus additional infrastructure
docker-compose --version # 2.21+

# Start development infrastructure
docker-compose -f docker-compose.dev.yml up -d

# Verify services
curl http://localhost:8222/healthz # NATS
curl http://localhost:9090 # Prometheus
curl http://localhost:16686 # Jaeger
```

**Plan 4 (MCP Plugin EDA)**:
```bash
# Go environment + NATS
go version # 1.21+
docker run -p 4222:4222 nats:latest --jetstream

# Test NATS connection
go run -c "
package main
import (
    \"github.com/nats-io/nats.go\"
    \"log\"
)
func main() {
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil { log.Fatal(err) }
    defer nc.Close()
    log.Println(\"NATS connected successfully\")
}
"
```

**Plan 5 (SuperClaude Implementation)**:
```bash
# TypeScript + Node.js with testing setup
node --version # 18+
npm install -g typescript vitest

# Verify testing framework
npx vitest --version
```

### IDE Setup

**VS Code Extensions**:
```json
{
  "recommendations": [
    "golang.go",
    "bradlc.vscode-tailwindcss", 
    "ms-vscode.vscode-typescript-next",
    "zxh404.vscode-proto3",
    "ms-vscode.vscode-json",
    "redhat.vscode-yaml",
    "ms-vscode.vscode-docker"
  ]
}
```

**VS Code Settings** (`.vscode/settings.json`):
```json
{
  "go.buildTags": "integration",
  "go.testTimeout": "60s",
  "go.lintTool": "golangci-lint",
  "typescript.preferences.importModuleSpecifier": "relative",
  "files.associations": {
    "*.proto": "proto3"
  }
}
```

---

## Build & Deployment Tools

### Go Build Tools

**Makefile Example**:
```makefile
.PHONY: build test clean proto docker

# Variables
BINARY_NAME=opencode
BUILD_DIR=./bin
PROTO_DIR=./proto
GO_FILES=$(shell find . -name "*.go" -type f)

# Build
build: proto
	CGO_ENABLED=0 GOOS=linux go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go

# Testing
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Protocol Buffers
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

# Docker
docker:
	docker build -t opencode:latest .

# Clean
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
```

### CI/CD Pipeline (GitHub Actions)

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test-go:
    runs-on: ubuntu-latest
    services:
      nats:
        image: nats:latest
        ports:
          - 4222:4222
        options: --health-cmd="wget --quiet --tries=1 --spider http://localhost:8222/healthz || exit 1"
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.22'
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3

  test-typescript:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '18'
        cache: 'npm'
    
    - name: Install dependencies
      run: npm ci
    
    - name: Run tests
      run: npm run test:coverage
    
    - name: Type check
      run: npm run type-check

  build-and-push:
    needs: [test-go, test-typescript]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v3
    
    - name: Login to Docker Registry
      uses: docker/login-action@v3
      with:
        registry: ghcr.io
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
    
    - name: Build and push
      uses: docker/build-push-action@v5
      with:
        context: .
        push: true
        tags: ghcr.io/${{ github.repository }}:latest
```

---

## Version Dependencies Summary

### Minimum Required Versions
| Technology | Minimum Version | Recommended | Plans |
|------------|----------------|-------------|-------|
| Go | 1.21.0 | 1.22.0+ | 1,3,4 |
| TypeScript | 5.2.0 | 5.3.0+ | 2,3,5 |
| Node.js | 18.0.0 | 20.0.0+ | 2,3,5 |
| Docker | 24.0.0 | 24.0.0+ | All |
| gRPC | 1.59.0 | 1.60.0+ | 1,3,4 |
| NATS Server | 2.10.0 | 2.10.7+ | 3,4 |
| Kafka | 3.6.0 | 3.6.1+ | 3,4 |
| Protobuf | 3.21.0 | 4.25.0+ | 1,3,4 |

### Breaking Changes to Watch
- **Go 1.21**: Generic type inference improvements
- **TypeScript 5.2**: New decorator implementation
- **NATS 2.10**: JetStream API changes
- **gRPC 1.60**: Transport security updates

---

This comprehensive technology stack reference provides all the information needed to set up, develop, and deploy any of the five OpenCode-SuperClaude integration plans. Each technology choice is justified by specific use cases and requirements across the different architectural approaches.