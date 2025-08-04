# Simple Plugin Example

This example demonstrates a basic plugin implementation using HashiCorp's go-plugin architecture.

## 🎯 Overview

This example shows:
- Basic plugin interface definition
- Plugin implementation with process isolation
- Host application that loads and manages plugins
- Secure RPC communication between host and plugin

## 📁 Structure

```
simple-plugin/
├── README.md           # This file
├── Makefile           # Build automation
├── shared/            # Shared interfaces
│   └── interface.go   # Plugin interface definition
├── plugin/            # Plugin implementation
│   ├── main.go       # Plugin entry point
│   └── hello.go      # Plugin logic
├── host/              # Host application
│   ├── main.go       # Host entry point
│   └── manager.go    # Plugin manager
└── plugin-hello/      # Built plugin binary
```

## 🚀 Quick Start

### Build and Run
```bash
# Build everything
make build

# Run the example
make run

# Clean build artifacts
make clean
```

### Expected Output
```
[HOST] Starting plugin system...
[HOST] Loading plugin: plugin-hello
[PLUGIN] Hello plugin initialized
[HOST] Executing command: greet
[PLUGIN] Hello from SuperClaude integration!
[HOST] Plugin response: Hello from SuperClaude integration!
[HOST] Shutting down...
```

## 💻 Code Walkthrough

### 1. Plugin Interface (`shared/interface.go`)
```go
type CommandPlugin interface {
    Name() string
    Execute(args map[string]interface{}) (string, error)
}
```

### 2. Plugin Implementation (`plugin/hello.go`)
```go
type HelloPlugin struct{}

func (p *HelloPlugin) Name() string {
    return "hello"
}

func (p *HelloPlugin) Execute(args map[string]interface{}) (string, error) {
    name := args["name"].(string)
    return fmt.Sprintf("Hello %s from SuperClaude!", name), nil
}
```

### 3. Host Application (`host/manager.go`)
- Discovers plugins in the plugins directory
- Starts plugin processes
- Manages plugin lifecycle
- Handles RPC communication

## 🔧 Configuration

### Plugin Discovery
Plugins are discovered from:
- `~/.opencode/plugins/` (production)
- `./plugins/` (development)

### Plugin Manifest
Each plugin can have a `plugin.json`:
```json
{
  "name": "hello",
  "version": "1.0.0",
  "description": "Simple greeting plugin",
  "author": "OpenCode Team",
  "capabilities": ["greet", "welcome"]
}
```

## 🧪 Testing

```bash
# Run unit tests
make test

# Run integration tests
make test-integration
```

## 🔒 Security

- Plugins run in separate processes
- Communication via secure RPC
- No shared memory between host and plugins
- Plugins can crash without affecting the host

## 📚 Next Steps

1. Add more complex plugin functionality
2. Implement plugin configuration
3. Add event handling (see [event-handler example](../event-handler/))
4. Integrate with MCP protocol (see [mcp-integration example](../mcp-integration/))

## 🔗 References

- [HashiCorp go-plugin](https://github.com/hashicorp/go-plugin)
- [ARCHITECTURE.md](../../ARCHITECTURE.md) - Full plugin architecture
- [Plugin API Specification](../../API-SPECIFICATIONS/plugin-api.md)