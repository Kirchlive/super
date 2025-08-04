# MCP Integration Example

This example demonstrates how to implement a Model Context Protocol (MCP) server that integrates with OpenCode.

## 🎯 Overview

This example shows:
- Basic MCP server implementation in Go
- JSON-RPC communication over STDIO
- Tool registration and invocation
- Integration with OpenCode's command system

## 📁 Structure

```
mcp-integration/
├── README.md           # This file
├── Makefile           # Build automation
├── server/            # MCP server implementation
│   ├── main.go       # Server entry point
│   ├── handler.go    # JSON-RPC handler
│   └── tools.go      # Tool implementations
├── client/            # Test client
│   └── main.go       # Client for testing
└── config/            # Configuration
    └── server.json    # Server configuration
```

## 🚀 Quick Start

### Build and Run
```bash
# Build the MCP server
make build

# Run the server
make run

# Test with client
make test-client

# Full integration test
make test-mcp
```

### Expected Output
```
[MCP] Server starting...
[MCP] Registered tool: explain
[MCP] Registered tool: implement
[MCP] Registered tool: optimize
[MCP] Listening on STDIO...
> {"jsonrpc":"2.0","method":"initialize","params":{},"id":1}
< {"jsonrpc":"2.0","result":{"protocolVersion":"2024-10-07"},"id":1}
```

## 💻 Code Walkthrough

### 1. MCP Server (`server/main.go`)
- Implements JSON-RPC server over STDIO
- Handles MCP protocol methods
- Manages tool registration

### 2. Tool Implementation (`server/tools.go`)
```go
type Tool struct {
    Name        string
    Description string
    InputSchema json.RawMessage
    Handler     func(args map[string]interface{}) (interface{}, error)
}
```

### 3. Protocol Handler (`server/handler.go`)
- Processes JSON-RPC requests
- Routes to appropriate handlers
- Manages request/response cycle

## 🔧 MCP Protocol

### Initialize
```json
{
  "jsonrpc": "2.0",
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-10-07",
    "capabilities": {}
  },
  "id": 1
}
```

### List Tools
```json
{
  "jsonrpc": "2.0",
  "method": "tools/list",
  "id": 2
}
```

### Call Tool
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "explain",
    "arguments": {
      "code": "function example() { return 42; }",
      "language": "javascript"
    }
  },
  "id": 3
}
```

## 🧪 Testing

```bash
# Unit tests
make test

# Integration test with OpenCode
make test-integration

# Manual testing with netcat
echo '{"jsonrpc":"2.0","method":"initialize","params":{},"id":1}' | ./mcp-server
```

## 🔒 Security

- Input validation on all tool calls
- Rate limiting for DOS protection
- Capability-based access control
- Audit logging for all operations

## 📚 Next Steps

1. Add more sophisticated tools
2. Implement streaming responses
3. Add authentication
4. Integrate with plugin system (see [simple-plugin example](../simple-plugin/))
5. Add event notifications (see [event-handler example](../event-handler/))

## 🔗 References

- [MCP Specification](https://modelcontextprotocol.io/docs)
- [MCP API Documentation](../../API-SPECIFICATIONS/mcp-api.md)
- [INTEGRATION.md](../../INTEGRATION.md) - Quick integration plan