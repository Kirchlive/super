# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This repository contains architectural and implementation plans for integrating OpenCode with SuperClaude capabilities. It currently consists of planning documents without active code implementation.

## Integration Plans Structure

The repository contains 5 distinct integration approaches, each with different complexity levels and timelines:

1. **ARCHITECTURE.md** - Plugin Ecosystem Architecture (8-10 weeks, Very High complexity)
   - HashiCorp go-plugin based extensible platform
   - Community-driven plugin marketplace vision
   - Process isolation for security

2. **INTEGRATION.md** - SuperClaude Integration (10 days, Medium complexity)
   - TypeScript-based rapid integration
   - Live-reload prompt templates
   - Sub-100ms performance targets

3. **EVENTDRIVEN.md** - MCP Bridge + Plugin + EventDriven (4-8 months, Very High complexity)
   - Three-tier enterprise architecture
   - Full MCP (Model Context Protocol) standard compliance
   - Real-time event capabilities

4. **EDA.md** - MCP Plugin Event-Driven Architecture (8 weeks, High complexity)
   - NATS JetStream for events
   - Production-ready plugin system
   - Edge computing ready

5. **IMPLEMENTATION.md** - SuperClaude Implementation (6 weeks, Medium-High complexity)
   - Intelligent prompt broker service
   - Comprehensive monitoring
   - Enterprise-grade configuration

## Key Technical Decisions

### Architecture Patterns
- **MCP (Model Context Protocol)**: Anthropic's standard for AI tool integration
- **Plugin Systems**: HashiCorp go-plugin for process isolation
- **Event Systems**: NATS JetStream or Kafka for event-driven architecture
- **Languages**: Go for core systems, TypeScript for integrations

### Technology Stack References
- **Go**: Core OpenCode implementation and plugin system
- **TypeScript**: SuperClaude integration and prompt management
- **gRPC**: Inter-process communication for plugins
- **NATS/Kafka**: Event bus implementations
- **Docker**: Containerization for deployment

## Plan Selection Guide

Use **PLAN-USAGE-GUIDE.md** as the primary reference for choosing between plans. Key decision factors:
- Timeline constraints (10 days to 8 months)
- Team size and expertise (3-8 developers)
- Architectural requirements (simple integration vs. full ecosystem)
- Strategic goals (rapid deployment vs. community platform)

## Working with Plans

When modifying or implementing plans:
1. Maintain consistency with the chosen architectural approach
2. Consider migration paths between plans (documented in PLAN-USAGE-GUIDE.md)
3. Preserve the German language sections where present (indicates international team context)
4. Reference specific plan sections when discussing implementation details

## Common Contexts

- **OpenCode**: The CLI tool being extended ([Repository](https://github.com/Kirchlive/opencode))
- **SuperClaude**: AI prompt management system being integrated ([Repository](https://github.com/Kirchlive/superclaude))
- **MCP**: Model Context Protocol - Anthropic's standard for AI tool integration
- **Plugin Architecture**: Isolated, extensible component system
- **Event-Driven Architecture (EDA)**: Reactive, loosely-coupled system design