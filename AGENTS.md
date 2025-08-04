# AGENTS.md

This guide is for AI agents working in this repository.

## Overview

This repository contains architectural and implementation plans for integrating OpenCode with SuperClaude capabilities. The project is in a pre-implementation phase, with a focus on Go for core systems and TypeScript for integrations.

## Commands

As there is no code yet, no build, lint, or test commands are defined. When implementing, please add standard commands and document them here:

*   **Go**: `go build`, `go test`, `gofmt`
*   **TypeScript**: `npm install`, `npm run build`, `npm run test`, `npm run lint`
*   **Integration Tests**: `npm run test:integration` (for end-to-end CLI tests)

## Code Style Guidelines

### General
*   **Languages**: Use Go for core systems (like the MCP server and plugin manager) and TypeScript for integrations and prompt management.
*   **Architecture**: Implementations must align with the chosen plan (e.g., MCP, Plugin Systems, EDA). Refer to `PLAN-USAGE-GUIDE.md` for details.
*   **Dependencies**: Minimize dependencies. New libraries must be vetted for security and performance.
*   **Comments**: Explain the *why*, not the *what*. Document complex logic, public APIs, and any non-obvious architectural decisions.

### Go
*   **Formatting**: Follow standard `gofmt` formatting.
*   **Dependency Management**: Use Go modules (`go.mod`/`go.sum`).
*   **Error Handling**: Implement comprehensive and idiomatic error handling.
*   **Plugin System**: When implementing a plugin system, use `hashicorp/go-plugin` for gRPC-based communication and process isolation.

### TypeScript
*   **Style Guide**: Adhere to a common style guide like Airbnb, enforced by ESLint and Prettier.
*   **Typing**: Employ strict typing (`"strict": true` in `tsconfig.json`). Use `zod` for runtime validation of data structures like prompt metadata.
*   **Modules**: Use ES modules and configure path aliases (e.g., `@superclaude/*`) for clean imports.
*   **File Watching**: Use `chokidar` for file watching and live-reloading of prompt templates.
*   **Reactive Programming**: Use `rxjs` for handling asynchronous events, especially for file changes.

### Naming Conventions
*   **Variables & Functions**: Use clear, descriptive names in camelCase (TypeScript) or CamelCase (Go public) / camelCase (Go private).
*   **Packages & Files**: Use lowercase, hyphenated names for packages (e.g., `prompt-broker`) and files.

## Project Structure
*   Maintain a monorepo structure with separate packages for different concerns (e.g., `@opencode/prompt-broker`, `@opencode/opencode-core`).
*   Store SuperClaude prompt templates in a dedicated directory (e.g., `superclaude-prompts/`) with subdirectories for commands and personas.
