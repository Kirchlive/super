# 📚 OpenCode-SuperClaude Integration Examples

This directory contains practical implementation examples for the OpenCode-SuperClaude integration project. Each example demonstrates key concepts from our architectural plans.

## 🗂️ Example Structure

### 1. [simple-plugin/](./simple-plugin/)
**Basic Plugin Implementation**
- Demonstrates the HashiCorp go-plugin architecture
- Simple "hello world" plugin with process isolation
- Shows plugin lifecycle management
- Based on concepts from [ARCHITECTURE.md](../ARCHITECTURE.md)

### 2. [mcp-integration/](./mcp-integration/)
**MCP Server Integration**
- Implements a basic MCP server in Go
- Shows JSON-RPC communication over STDIO
- Demonstrates tool registration and invocation
- Based on [MCP API Specification](../API-SPECIFICATIONS/mcp-api.md)

### 3. [event-handler/](./event-handler/)
**Event-Driven Architecture Example**
- NATS JetStream event publishing and subscription
- Configuration hot-reload using events
- Plugin lifecycle event handling
- Based on [EDA.md](../EDA.md) and [Event API](../API-SPECIFICATIONS/event-api.md)

## 🚀 Quick Start

Each example includes:
- Complete source code with comments
- README with setup instructions
- Makefile for building and running
- Tests demonstrating usage

### Prerequisites
```bash
# Go 1.21+
go version

# Node.js 18+ (for TypeScript examples)
node --version

# NATS Server (for event examples)
# Install via: brew install nats-server
```

### Running Examples

```bash
# Simple Plugin
cd simple-plugin
make build
make run

# MCP Integration
cd mcp-integration
make build
make test-mcp

# Event Handler
cd event-handler
docker-compose up -d  # Start NATS
make run
```

## 📖 Learning Path

1. **Start with `simple-plugin`** to understand basic plugin architecture
2. **Move to `mcp-integration`** to see MCP protocol implementation
3. **Explore `event-handler`** for event-driven patterns
4. **Combine concepts** to build your own integration

## 🔗 Related Documentation

- [API Specifications](../API-SPECIFICATIONS/) - Detailed API documentation
- [TECHNICAL-GLOSSARY.md](../TECHNICAL-GLOSSARY.md) - Technical terms explained
- [TESTING-STRATEGY.md](../TESTING-STRATEGY.md) - Testing approaches
- [CONTRIBUTING.md](../CONTRIBUTING.md) - How to contribute examples

## 💡 Contributing Examples

We welcome new examples! Please:
1. Follow the existing structure
2. Include comprehensive comments
3. Add tests for your example
4. Update this README
5. Submit a PR following our [contribution guidelines](../CONTRIBUTING.md)

---

*These examples are designed to be educational and serve as starting points for your own implementations.*