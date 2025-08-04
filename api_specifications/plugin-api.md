# Plugin Interface API Specification

Plugin architecture API specification for OpenCode SuperClaude integration using HashiCorp go-plugin pattern with gRPC communication.

## Overview

The Plugin API enables isolated, process-separated extensions that integrate with OpenCode through a standardized interface. Plugins communicate via gRPC and are managed by the Plugin Manager within the MCP Server Bridge.

## Architecture

```
┌─────────────────┐    gRPC     ┌─────────────────┐    Events    ┌─────────────────┐
│ Plugin Manager  │◄───────────►│ Plugin Process  │◄────────────►│ Event Bus       │
│ (Host)          │             │ (Isolated)      │              │ (NATS/Kafka)    │
└─────────────────┘             └─────────────────┘              └─────────────────┘
```

## Core Plugin Interface

### Go Interface Definition

```go
// ICommandPlugin defines the core contract for all SuperClaude plugins
type ICommandPlugin interface {
    // Metadata returns plugin information
    GetMetadata() (*PluginMetadata, error)
    
    // Lifecycle management
    Initialize(ctx context.Context, config *PluginConfig) error
    Shutdown(ctx context.Context) error
    
    // Core execution
    Execute(ctx context.Context, request *ExecuteRequest) (*ExecuteResponse, error)
    
    // Optional capabilities
    ValidateRequest(ctx context.Context, request *ExecuteRequest) (*ValidationResult, error)
    GetSchema() (*PluginSchema, error)
    
    // Health and monitoring
    HealthCheck(ctx context.Context) (*HealthStatus, error)
}
```

### gRPC Service Definition

```protobuf
syntax = "proto3";

package superclaudebridge.v1;

option go_package = "github.com/opencode/superclaudebridge/proto/v1";

// SuperClaude Plugin Service
service PluginService {
    // Lifecycle methods
    rpc Initialize(InitializeRequest) returns (InitializeResponse);
    rpc Shutdown(ShutdownRequest) returns (ShutdownResponse);
    
    // Core functionality
    rpc Execute(ExecuteRequest) returns (ExecuteResponse);
    rpc ValidateRequest(ValidateRequest) returns (ValidateResponse);
    rpc GetMetadata(MetadataRequest) returns (MetadataResponse);
    rpc GetSchema(SchemaRequest) returns (SchemaResponse);
    
    // Health and monitoring
    rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
    
    // Event handling
    rpc HandleEvent(EventRequest) returns (EventResponse);
}

// Core data structures
message PluginMetadata {
    string name = 1;
    string version = 2;
    string description = 3;
    string author = 4;
    repeated string tags = 5;
    PluginCapabilities capabilities = 6;
    map<string, string> configuration = 7;
}

message PluginCapabilities {
    bool supports_streaming = 1;
    bool supports_cancellation = 2;
    bool supports_events = 3;
    bool requires_filesystem = 4;
    bool requires_network = 5;
    repeated string supported_contexts = 6;
}

message ExecuteRequest {
    string command = 1;
    map<string, string> arguments = 2;
    ExecutionContext context = 3;
    map<string, string> environment = 4;
}

message ExecuteResponse {
    ExecutionResult result = 1;
    repeated LogEntry logs = 2;
    map<string, string> metadata = 3;
    bool requires_followup = 4;
}

message ExecutionContext {
    string working_directory = 1;
    string project_root = 2;
    repeated string file_paths = 3;
    string user_id = 4;
    string session_id = 5;
    map<string, string> environment_vars = 6;
}

message ExecutionResult {
    enum Status {
        SUCCESS = 0;
        ERROR = 1;
        PARTIAL = 2;
        CANCELLED = 3;
    }
    
    Status status = 1;
    string output = 2;
    string error_message = 3;
    repeated Artifact artifacts = 4;
    map<string, string> metrics = 5;
}

message Artifact {
    enum Type {
        FILE = 0;
        DIRECTORY = 1;
        CODE_CHANGE = 2;
        DOCUMENTATION = 3;
        CONFIGURATION = 4;
    }
    
    Type type = 1;
    string path = 2;
    string content = 3;
    map<string, string> metadata = 4;
}

message ValidationResult {
    bool is_valid = 1;
    repeated ValidationError errors = 2;
    repeated string warnings = 3;
}

message ValidationError {
    string field = 1;
    string message = 2;
    string code = 3;
}
```

## Plugin Lifecycle

### 1. Discovery and Loading

```go
// Plugin Manager discovers plugins in directory
func (pm *PluginManager) DiscoverPlugins(directory string) error {
    files, err := filepath.Glob(filepath.Join(directory, "*.exe"))
    if err != nil {
        return err
    }
    
    for _, file := range files {
        plugin, err := pm.LoadPlugin(file)
        if err != nil {
            log.Warnf("Failed to load plugin %s: %v", file, err)
            continue
        }
        
        pm.plugins[plugin.Name] = plugin
    }
    
    return nil
}
```

### 2. Plugin Registration

```go
// Plugin registers with HashiCorp go-plugin
func main() {
    plugin.Serve(&plugin.ServeConfig{
        HandshakeConfig: handshakeConfig,
        Plugins: map[string]plugin.Plugin{
            "superclaude": &SuperClaudePlugin{
                Impl: &MyPlugin{}, // Your plugin implementation
            },
        },
        GRPCServer: plugin.DefaultGRPCServer,
    })
}

var handshakeConfig = plugin.HandshakeConfig{
    ProtocolVersion:  1,
    MagicCookieKey:   "SUPERCLAUDE_PLUGIN",
    MagicCookieValue: "superclaudebridge",
}
```

### 3. Initialization

```go
func (p *MyPlugin) Initialize(ctx context.Context, config *PluginConfig) error {
    // Setup plugin state
    p.config = config
    p.eventClient = NewEventClient(config.EventBusURL)
    
    // Register for events
    if err := p.eventClient.Subscribe("config.changed", p.handleConfigChange); err != nil {
        return fmt.Errorf("failed to subscribe to events: %w", err)
    }
    
    // Initialize resources
    return p.setupResources()
}
```

## Plugin Development Kit (PDK)

### Go Plugin Template

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/hashicorp/go-plugin"
    "github.com/opencode/superclaudebridge/sdk"
)

// MyPlugin implements the ICommandPlugin interface
type MyPlugin struct {
    config *sdk.PluginConfig
    client *sdk.EventClient
}

func (p *MyPlugin) GetMetadata() (*sdk.PluginMetadata, error) {
    return &sdk.PluginMetadata{
        Name:        "my-superclaude-plugin",
        Version:     "1.0.0", 
        Description: "Example SuperClaude plugin",
        Author:      "Your Name",
        Tags:        []string{"example", "template"},
        Capabilities: &sdk.PluginCapabilities{
            SupportsStreaming:     false,
            SupportsCancellation: true,
            SupportsEvents:       true,
            RequiresFilesystem:   true,
            RequiresNetwork:      false,
            SupportedContexts:    []string{"file", "directory", "project"},
        },
    }, nil
}

func (p *MyPlugin) Execute(ctx context.Context, req *sdk.ExecuteRequest) (*sdk.ExecuteResponse, error) {
    switch req.Command {
    case "analyze":
        return p.handleAnalyze(ctx, req)
    case "generate":
        return p.handleGenerate(ctx, req)
    default:
        return nil, fmt.Errorf("unknown command: %s", req.Command)
    }
}

func (p *MyPlugin) handleAnalyze(ctx context.Context, req *sdk.ExecuteRequest) (*sdk.ExecuteResponse, error) {
    // Your analysis logic here
    
    return &sdk.ExecuteResponse{
        Result: &sdk.ExecutionResult{
            Status: sdk.ExecutionResult_SUCCESS,
            Output: "Analysis completed successfully",
            Artifacts: []*sdk.Artifact{
                {
                    Type:    sdk.Artifact_DOCUMENTATION,
                    Path:    "analysis-report.md",
                    Content: "# Analysis Report\n\nDetails...",
                },
            },
        },
    }, nil
}

// Plugin setup
func main() {
    plugin.Serve(&plugin.ServeConfig{
        HandshakeConfig: sdk.HandshakeConfig,
        Plugins: map[string]plugin.Plugin{
            "superclaude": &sdk.SuperClaudePlugin{
                Impl: &MyPlugin{},
            },
        },
        GRPCServer: plugin.DefaultGRPCServer,
    })
}
```

### TypeScript Plugin Template

```typescript
import { PluginService } from './generated/plugin_grpc_pb';
import { 
    ExecuteRequest, 
    ExecuteResponse, 
    PluginMetadata,
    ExecutionResult 
} from './generated/plugin_pb';

export class MySuperClaudePlugin implements PluginService {
    private config: PluginConfig;
    private eventClient: EventClient;

    async initialize(config: PluginConfig): Promise<void> {
        this.config = config;
        this.eventClient = new EventClient(config.eventBusUrl);
        
        // Subscribe to events
        await this.eventClient.subscribe('config.changed', this.handleConfigChange.bind(this));
    }

    async getMetadata(): Promise<PluginMetadata> {
        const metadata = new PluginMetadata();
        metadata.setName('my-ts-plugin');
        metadata.setVersion('1.0.0');
        metadata.setDescription('TypeScript SuperClaude plugin');
        metadata.setAuthor('Your Name');
        
        return metadata;
    }

    async execute(request: ExecuteRequest): Promise<ExecuteResponse> {
        const command = request.getCommand();
        const response = new ExecuteResponse();
        
        switch (command) {
            case 'process':
                return this.handleProcess(request);
            case 'validate':
                return this.handleValidate(request);
            default:
                throw new Error(`Unknown command: ${command}`);
        }
    }

    private async handleProcess(request: ExecuteRequest): Promise<ExecuteResponse> {
        // Your processing logic here
        const result = new ExecutionResult();
        result.setStatus(ExecutionResult.Status.SUCCESS);
        result.setOutput('Processing completed');
        
        const response = new ExecuteResponse();
        response.setResult(result);
        
        return response;
    }
}

// Plugin registration
if (require.main === module) {
    const server = new grpc.Server();
    server.addService(PluginService, new MySuperClaudePlugin());
    
    const port = process.env.PLUGIN_PORT || '50051';
    server.bindAsync(`0.0.0.0:${port}`, grpc.ServerCredentials.createInsecure(), () => {
        console.log(`Plugin server running on port ${port}`);
        server.start();
    });
}
```

## Communication Protocols

### Host-Plugin Communication

```go
// Plugin Manager communicates with plugins via gRPC
type PluginClient struct {
    client pb.PluginServiceClient
    conn   *grpc.ClientConn
}

func (pc *PluginClient) CallPlugin(ctx context.Context, command string, args map[string]string) (*pb.ExecuteResponse, error) {
    req := &pb.ExecuteRequest{
        Command:   command,
        Arguments: args,
        Context: &pb.ExecutionContext{
            WorkingDirectory: pc.getWorkingDir(),
            ProjectRoot:      pc.getProjectRoot(),
            SessionId:        pc.getSessionID(),
        },
    }
    
    return pc.client.Execute(ctx, req)
}
```

### Event-Driven Communication

```go
// Plugins can emit and receive events
func (p *MyPlugin) EmitEvent(eventType string, data map[string]interface{}) error {
    event := &EventMessage{
        Type:      eventType,
        Source:    p.GetMetadata().Name,
        Data:      data,
        Timestamp: time.Now(),
    }
    
    return p.eventClient.Publish(event)
}

func (p *MyPlugin) handleConfigChange(event *EventMessage) error {
    // React to configuration changes
    return p.reloadConfiguration()
}
```

## Plugin Security Model

### Sandboxing and Isolation

```go
// Plugin configuration with security constraints
type PluginConfig struct {
    // Filesystem access
    AllowedReadPaths  []string `json:"allowed_read_paths"`
    AllowedWritePaths []string `json:"allowed_write_paths"`
    
    // Network access
    AllowedHosts []string `json:"allowed_hosts"`
    NetworkEnabled bool   `json:"network_enabled"`
    
    // Environment access
    AllowedEnvVars []string `json:"allowed_env_vars"`
    
    // Resource limits
    MaxMemoryMB     int           `json:"max_memory_mb"`
    MaxExecutionTime time.Duration `json:"max_execution_time"`
    MaxFileSize     int64         `json:"max_file_size"`
}
```

### Permission System

```go
// Permission checks before plugin execution
func (pm *PluginManager) ValidatePermissions(plugin *Plugin, request *ExecuteRequest) error {
    // Check filesystem access
    for _, path := range request.Context.FilePaths {
        if !pm.isPathAllowed(plugin, path) {
            return fmt.Errorf("plugin %s not allowed to access path: %s", plugin.Name, path)
        }
    }
    
    // Check network access
    if request.RequiresNetwork && !plugin.Config.NetworkEnabled {
        return fmt.Errorf("plugin %s not allowed network access", plugin.Name)
    }
    
    // Check resource limits
    if err := pm.checkResourceLimits(plugin, request); err != nil {
        return fmt.Errorf("resource limit violation: %w", err)
    }
    
    return nil
}
```

## Plugin Management

### Hot Reloading

```go
// Hot reload plugin without restarting host
func (pm *PluginManager) ReloadPlugin(name string) error {
    // Stop existing plugin
    if existing, exists := pm.plugins[name]; exists {
        if err := existing.Shutdown(context.Background()); err != nil {
            log.Warnf("Error shutting down plugin %s: %v", name, err)
        }
    }
    
    // Load new version
    newPlugin, err := pm.LoadPlugin(filepath.Join(pm.pluginDir, name+".exe"))
    if err != nil {
        return fmt.Errorf("failed to reload plugin %s: %w", name, err)
    }
    
    // Initialize new plugin
    if err := newPlugin.Initialize(context.Background(), pm.getPluginConfig(name)); err != nil {
        return fmt.Errorf("failed to initialize reloaded plugin %s: %w", name, err)
    }
    
    pm.plugins[name] = newPlugin
    return nil
}
```

### Plugin Registry

```go
// Plugin registry for metadata and discovery
type PluginRegistry struct {
    plugins map[string]*PluginEntry
    mu      sync.RWMutex
}

type PluginEntry struct {
    Metadata    *PluginMetadata
    BinaryPath  string
    Config      *PluginConfig
    Status      PluginStatus
    LastUpdated time.Time
}

func (pr *PluginRegistry) Register(plugin *PluginEntry) error {
    pr.mu.Lock()
    defer pr.mu.Unlock()
    
    pr.plugins[plugin.Metadata.Name] = plugin
    return nil
}

func (pr *PluginRegistry) List() []*PluginEntry {
    pr.mu.RLock()
    defer pr.mu.RUnlock()
    
    var entries []*PluginEntry
    for _, entry := range pr.plugins {
        entries = append(entries, entry)
    }
    
    return entries
}
```

## Error Handling and Monitoring

### Error Propagation

```go
// Structured error handling
type PluginError struct {
    Plugin    string    `json:"plugin"`
    Command   string    `json:"command"`
    Message   string    `json:"message"`
    Code      ErrorCode `json:"code"`
    Timestamp time.Time `json:"timestamp"`
    Context   map[string]interface{} `json:"context"`
}

const (
    ErrorCodeInvalidRequest ErrorCode = "INVALID_REQUEST"
    ErrorCodePermissionDenied ErrorCode = "PERMISSION_DENIED"  
    ErrorCodeResourceLimit ErrorCode = "RESOURCE_LIMIT"
    ErrorCodePluginCrash ErrorCode = "PLUGIN_CRASH"
    ErrorCodeTimeout ErrorCode = "TIMEOUT"
)
```

### Health Monitoring

```go
// Health check implementation
func (p *MyPlugin) HealthCheck(ctx context.Context) (*HealthStatus, error) {
    status := &HealthStatus{
        Status:    HealthStatus_HEALTHY,
        Timestamp: time.Now().Unix(),
        Checks: map[string]*HealthCheck{
            "memory": {
                Status:  HealthCheck_PASS,
                Message: fmt.Sprintf("Memory usage: %d MB", p.getMemoryUsage()),
            },
            "dependencies": {
                Status:  HealthCheck_PASS,
                Message: "All dependencies available",
            },
        },
    }
    
    // Check if we can access required resources
    if err := p.checkDependencies(); err != nil {
        status.Status = HealthStatus_UNHEALTHY
        status.Checks["dependencies"] = &HealthCheck{
            Status:  HealthCheck_FAIL,
            Message: fmt.Sprintf("Dependency check failed: %v", err),
        }
    }
    
    return status, nil
}
```

## Build and Deployment

### Build Configuration

```yaml
# .goreleaser.yml
builds:
  - id: plugin
    main: ./cmd/plugin
    binary: my-superclaude-plugin
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}

archives:
  - id: plugin-archive
    builds:
      - plugin
    format: tar.gz
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: 'checksums.txt'

signs:
  - artifacts: checksum
    args: ["--batch", "-u", "{{ .Env.GPG_FINGERPRINT }}", "--output", "${signature}", "--detach-sign", "${artifact}"]
```

### Plugin Manifest

```json
{
  "name": "my-superclaude-plugin",
  "version": "1.0.0",
  "description": "Example SuperClaude plugin",
  "author": "Your Name",
  "license": "MIT",
  "repository": "https://github.com/yourorg/my-superclaude-plugin",
  "binary": {
    "linux-amd64": "my-superclaude-plugin-linux-amd64",
    "darwin-amd64": "my-superclaude-plugin-darwin-amd64", 
    "windows-amd64": "my-superclaude-plugin-windows-amd64.exe"
  },
  "checksum": {
    "sha256": "abc123...",
    "file": "checksums.txt"
  },
  "capabilities": {
    "filesystem": true,
    "network": false,
    "events": true
  },
  "configuration": {
    "max_memory_mb": 512,
    "timeout_seconds": 300
  }
}
```

## References

- [HashiCorp go-plugin Documentation](https://github.com/hashicorp/go-plugin)
- [gRPC Protocol Buffers Guide](https://grpc.io/docs/languages/go/basics/)
- [OpenCode Plugin Architecture](../ARCHITECTURE.md)
- [Event System Specification](./event-api.md)
- [MCP Protocol Integration](./mcp-api.md)