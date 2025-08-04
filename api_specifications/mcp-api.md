# MCP Protocol API Specification

Model Context Protocol (MCP) API specification for OpenCode integration with SuperClaude framework.

## Overview

The MCP Protocol provides standardized communication between OpenCode (client) and SuperClaude tools (server) using JSON-RPC over STDIO/HTTP transport layers. This specification defines the complete protocol interface, message formats, and integration patterns.

## Architecture

```
┌─────────────┐    JSON-RPC    ┌─────────────┐    gRPC/Events   ┌─────────────┐
│ OpenCode    │◄──────────────►│ MCP Server  │◄────────────────►│ Plugin      │
│ (Client)    │                │ Bridge      │                  │ Runtime     │
└─────────────┘                └─────────────┘                  └─────────────┘
```

## Protocol Specification

### 1. Transport Layer

#### STDIO Transport
```json
{
  "jsonrpc": "2.0",
  "method": "initialize",
  "params": {
    "protocolVersion": "2025-03-26",
    "capabilities": {
      "roots": {"listChanged": true},
      "sampling": {}
    },
    "clientInfo": {
      "name": "opencode",
      "version": "1.0.0"
    }
  },
  "id": 1
}
```

#### HTTP Transport (Optional)
```http
POST /mcp HTTP/1.1
Content-Type: application/json
Authorization: Bearer <token>

{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {...},
  "id": 1
}
```

### 2. Core Protocol Methods

#### Initialize Handshake

**Request: `initialize`**
```json
{
  "jsonrpc": "2.0",
  "method": "initialize",
  "params": {
    "protocolVersion": "2025-03-26",
    "capabilities": {
      "roots": {"listChanged": true},
      "sampling": {},
      "experimental": {}
    },
    "clientInfo": {
      "name": "opencode",
      "version": "1.0.0"
    }
  },
  "id": 1
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "protocolVersion": "2025-03-26",
    "capabilities": {
      "logging": {},
      "tools": {"listChanged": true},
      "prompts": {"listChanged": true},
      "resources": {"subscribe": true, "listChanged": true}
    },
    "serverInfo": {
      "name": "superclaude-mcp-bridge",
      "version": "1.0.0"
    }
  },
  "id": 1
}
```

#### Tool Discovery and Execution

**Request: `tools/list`**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/list",
  "id": 2
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "tools": [
      {
        "name": "sc-analyze",
        "description": "Analyze code with SuperClaude intelligence",
        "inputSchema": {
          "type": "object",
          "properties": {
            "path": {"type": "string", "description": "File or directory path"},
            "focus": {"type": "string", "enum": ["performance", "security", "quality"]},
            "depth": {"type": "string", "enum": ["shallow", "deep", "comprehensive"]}
          },
          "required": ["path"]
        }
      },
      {
        "name": "sc-implement",
        "description": "Implement features using SuperClaude patterns",
        "inputSchema": {
          "type": "object",
          "properties": {
            "description": {"type": "string"},
            "framework": {"type": "string"},
            "type": {"type": "string", "enum": ["component", "api", "service", "feature"]}
          },
          "required": ["description"]
        }
      }
    ]
  },
  "id": 2
}
```

**Request: `tools/call`**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "sc-analyze",
    "arguments": {
      "path": "./src/components",
      "focus": "performance",
      "depth": "deep"
    }
  },
  "id": 3
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Performance analysis completed:\n\n## Critical Issues\n- Bundle size: 2.3MB (target: <500KB)\n- Unused dependencies: 15 packages\n- Render optimization opportunities: 8 components\n\n## Recommendations\n1. Implement code splitting\n2. Add React.memo to heavy components\n3. Optimize image loading"
      }
    ],
    "isError": false
  },
  "id": 3
}
```

#### Prompt Management

**Request: `prompts/list`**
```json
{
  "jsonrpc": "2.0",
  "method": "prompts/list",
  "id": 4
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "prompts": [
      {
        "name": "sc-context-builder",
        "description": "Build comprehensive context for SuperClaude analysis",
        "arguments": [
          {
            "name": "scope",
            "description": "Analysis scope",
            "required": true
          }
        ]
      }
    ]
  },
  "id": 4
}
```

### 3. Resource Management

#### Resource Discovery
```json
{
  "jsonrpc": "2.0",
  "method": "resources/list",
  "id": 5
}
```

#### Resource Subscription
```json
{
  "jsonrpc": "2.0",
  "method": "resources/subscribe",
  "params": {
    "uri": "file://./src"
  },
  "id": 6
}
```

### 4. Notification System

#### Server Notifications
```json
{
  "jsonrpc": "2.0",
  "method": "notifications/tools/list_changed",
  "params": {}
}
```

#### Client Notifications
```json
{
  "jsonrpc": "2.0",
  "method": "notifications/cancelled",
  "params": {
    "requestId": 3,
    "reason": "User cancelled operation"
  }
}
```

## Error Handling

### Standard Error Codes
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32602,
    "message": "Invalid params",
    "data": {
      "parameter": "focus",
      "expected": ["performance", "security", "quality"],
      "received": "invalid-focus"
    }
  },
  "id": 3
}
```

### Error Code Reference
| Code | Message | Description |
|------|---------|-------------|
| -32700 | Parse error | Invalid JSON |
| -32600 | Invalid Request | Invalid JSON-RPC |
| -32601 | Method not found | Unknown method |
| -32602 | Invalid params | Invalid parameters |
| -32603 | Internal error | Server error |
| -32000 | Plugin error | Plugin execution failed |
| -32001 | Timeout error | Operation timeout |

## Integration with OpenCode

### Client Configuration
```json
{
  "mcpServers": {
    "superclaudebridge": {
      "command": "superclaudebridge",
      "args": ["--config", "~/.opencode/superclaudebridge.json"],
      "env": {
        "OPENCODE_VERSION": "1.0.0"
      },
      "initializationOptions": {
        "pluginDirectory": "~/.opencode/plugins",
        "eventBusUrl": "nats://localhost:4222"
      }
    }
  }
}
```

### Authentication
```json
{
  "jsonrpc": "2.0",
  "method": "initialize",
  "params": {
    "protocolVersion": "2025-03-26",
    "capabilities": {},
    "clientInfo": {
      "name": "opencode",
      "version": "1.0.0"
    },
    "auth": {
      "type": "bearer",
      "token": "<jwt-token>"
    }
  },
  "id": 1
}
```

## Performance Considerations

### Batching Requests
```json
[
  {
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {"name": "sc-analyze", "arguments": {...}},
    "id": 1
  },
  {
    "jsonrpc": "2.0",
    "method": "tools/call", 
    "params": {"name": "sc-lint", "arguments": {...}},
    "id": 2
  }
]
```

### Streaming Responses
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Analysis progress: 25%"
      }
    ],
    "isError": false,
    "streaming": true
  },
  "id": 3
}
```

## Security Model

### Capability-Based Access
```json
{
  "capabilities": {
    "filesystem": {
      "read": ["./src", "./docs"],
      "write": ["./src"]
    },
    "network": {
      "allowed_hosts": ["api.github.com"]
    },
    "environment": {
      "allowed_vars": ["NODE_ENV", "OPENCODE_*"]
    }
  }
}
```

### Sandboxing
- All plugins run in isolated processes
- Resource access controlled by capability declarations
- Network requests limited to allowlisted domains
- File system access restricted to declared paths

## OpenCode Integration Examples

### Command Integration
```go
// OpenCode command handler
func (c *MCPCommand) Execute(ctx context.Context, args []string) error {
    client := mcp.NewClient(c.serverPath)
    
    result, err := client.CallTool(ctx, &mcp.CallToolRequest{
        Name: "sc-analyze",
        Arguments: map[string]interface{}{
            "path": args[0],
            "focus": "performance",
        },
    })
    
    if err != nil {
        return fmt.Errorf("MCP call failed: %w", err)
    }
    
    c.ui.Output(result.Content[0].Text)
    return nil
}
```

### Event Integration
```go
// Listen for MCP notifications
go func() {
    for notification := range client.Notifications() {
        switch notification.Method {
        case "notifications/tools/list_changed":
            c.refreshToolCache()
        case "notifications/progress":
            c.updateProgressBar(notification.Params)
        }
    }
}()
```

## Version Compatibility

### Protocol Versioning
- Current version: `2025-03-26`
- Backward compatibility maintained for 2 major versions
- Feature detection via capabilities negotiation
- Graceful degradation for unsupported features

### Migration Guide
```json
{
  "migration": {
    "from": "2024-11-05",
    "to": "2025-03-26",
    "changes": [
      {
        "type": "added",
        "feature": "resource_subscription",
        "description": "Added real-time resource monitoring"
      },
      {
        "type": "deprecated", 
        "feature": "legacy_prompts",
        "replacement": "prompts/list",
        "removal_version": "2025-09-01"
      }
    ]
  }
}
```

## Testing and Development

### Protocol Testing
```bash
# Test MCP server compliance
mcp-test-suite --server ./superclaudebridge --spec 2025-03-26

# Interactive testing
mcp-inspector --connect stdio --command ./superclaudebridge
```

### Development Tools
- MCP Inspector: Interactive protocol debugging
- Schema validation: JSON schema validation for all messages
- Performance profiler: Latency and throughput analysis
- Mock server: Testing client implementations

## References

- [Model Context Protocol Specification](https://modelcontextprotocol.io/docs/concepts/architecture)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)
- [OpenCode MCP Integration Guide](../INTEGRATION.md)
- [SuperClaude Plugin Development](./plugin-api.md)