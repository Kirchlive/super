# Contributing to OpenCode-SuperClaude Integration

## 🤝 Welcome Contributors!

Thank you for your interest in contributing to the OpenCode-SuperClaude Integration project! This document provides guidelines for contributing to our repository.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Environment Setup](#development-environment-setup)
- [Coding Standards](#coding-standards)
- [Commit Message Conventions](#commit-message-conventions)
- [Pull Request Process](#pull-request-process)
- [Documentation Standards](#documentation-standards)
- [Testing Requirements](#testing-requirements)
- [Review Process](#review-process)

## 🤖 Code of Conduct

### Our Pledge

We pledge to make participation in our project a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, sex characteristics, gender identity and expression, level of experience, education, socio-economic status, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Standards

**Examples of behavior that contributes to a positive environment:**

- Using welcoming and inclusive language
- Being respectful of differing viewpoints and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

**Examples of unacceptable behavior:**

- The use of sexualized language or imagery and unwelcome sexual attention or advances
- Trolling, insulting/derogatory comments, and personal or political attacks
- Public or private harassment
- Publishing others' private information without explicit permission
- Other conduct which could reasonably be considered inappropriate in a professional setting

### Enforcement

Instances of abusive, harassing, or otherwise unacceptable behavior may be reported by contacting the project maintainers. All complaints will be reviewed and investigated promptly and fairly.

## 🚀 Getting Started

### Prerequisites

- **Go**: Version 1.21 or higher
- **Node.js**: Version 18 or higher (for TypeScript development)
- **Git**: Version 2.20 or higher
- **Docker**: For containerized development (optional)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:

```bash
git clone https://github.com/yourusername/opencode-superclaude.git
cd opencode-superclaude
```

3. Add the upstream repository:

```bash
git remote add upstream https://github.com/original/opencode-superclaude.git
```

## 🛠️ Development Environment Setup

### Go Development Setup

1. **Install Go dependencies:**

```bash
go mod download
go mod tidy
```

2. **Install development tools:**

```bash
# Code formatting and linting
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Testing tools
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install github.com/onsi/gomega/...@latest
```

3. **Set up pre-commit hooks:**

```bash
git config core.hooksPath .githooks
chmod +x .githooks/pre-commit
```

### TypeScript Development Setup

1. **Install Node.js dependencies:**

```bash
npm install
# or
yarn install
```

2. **Install TypeScript tools:**

```bash
npm install -g typescript@latest
npm install -g @typescript-eslint/parser @typescript-eslint/eslint-plugin
```

### Environment Configuration

1. **Copy environment template:**

```bash
cp .env.example .env
```

2. **Configure local settings:**

```bash
# .env
OPENCODE_LOG_LEVEL=debug
OPENCODE_PLUGIN_DIR=./plugins
MCP_SERVER_PORT=8080
```

## 📝 Coding Standards

### Go Coding Standards

#### Code Organization

```go
// Package structure
package main

import (
    // Standard library imports first
    "context"
    "fmt"
    "log"
    
    // Third-party imports
    "github.com/hashicorp/go-plugin"
    
    // Local imports
    "github.com/opencode/internal/plugin"
)
```

#### Naming Conventions

- **Packages**: Short, lowercase, single word when possible
- **Functions/Methods**: CamelCase, start with verb
- **Variables**: camelCase for local, PascalCase for exported
- **Constants**: ALL_CAPS with underscores
- **Interfaces**: End with "er" when possible (e.g., `CommandRunner`)

#### Error Handling

```go
// Always handle errors explicitly
result, err := someFunction()
if err != nil {
    return fmt.Errorf("failed to execute function: %w", err)
}

// Use context for cancellation
func processCommand(ctx context.Context, cmd string) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Process command
    }
}
```

#### Documentation

```go
// CommandPlugin defines the interface for OpenCode plugins.
// All plugins must implement this interface to be discoverable
// by the plugin manager.
type CommandPlugin interface {
    // Name returns the unique identifier for this plugin.
    // This name is used for command routing and plugin discovery.
    Name() string
    
    // Execute runs the plugin with the provided context and arguments.
    // It returns a result object or an error if execution fails.
    Execute(ctx context.Context, args []string) (*Result, error)
}
```

### TypeScript Coding Standards

#### Project Structure

```
src/
├── types/           # Type definitions
├── services/        # Business logic
├── utils/           # Utility functions
├── integration/     # MCP integration layer
└── __tests__/       # Test files
```

#### Type Definitions

```typescript
// Use interfaces for object shapes
interface PluginConfig {
  name: string;
  version: string;
  commands: CommandDefinition[];
  dependencies?: string[];
}

// Use type aliases for unions and primitives
type LogLevel = 'debug' | 'info' | 'warn' | 'error';
type PluginStatus = 'active' | 'inactive' | 'error';

// Use enums for constants
enum MessageType {
  Command = 'command',
  Response = 'response',
  Error = 'error'
}
```

#### Async/Await Patterns

```typescript
// Prefer async/await over Promises
async function executePlugin(name: string, args: string[]): Promise<PluginResult> {
  try {
    const plugin = await pluginManager.getPlugin(name);
    const result = await plugin.execute(args);
    return result;
  } catch (error) {
    logger.error(`Plugin execution failed: ${error.message}`);
    throw new PluginExecutionError(error.message);
  }
}
```

## 📝 Commit Message Conventions

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

### Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **chore**: Changes to the build process or auxiliary tools

### Examples

```
feat(plugin): add hot-reload capability for plugins

Implement hot-reload functionality that allows plugins to be
reloaded without restarting the MCP server.

- Add file system watcher for plugin directory
- Implement graceful plugin shutdown and restart
- Update plugin manager to handle reload events

Closes #123
```

```
fix(mcp): resolve connection timeout issues

- Increase default timeout from 5s to 30s
- Add exponential backoff for connection retries
- Improve error messages for connection failures

Fixes #456
```

## 🔄 Pull Request Process

### Before Creating a PR

1. **Create a feature branch:**

```bash
git checkout -b feat/plugin-hot-reload
```

2. **Make your changes and commit:**

```bash
git add .
git commit -m "feat(plugin): add hot-reload capability"
```

3. **Sync with upstream:**

```bash
git fetch upstream
git rebase upstream/main
```

4. **Run tests and linting:**

```bash
# Go
go test ./...
golangci-lint run

# TypeScript
npm run test
npm run lint
```

### PR Guidelines

#### PR Title Format

Use the same format as commit messages:

```
feat(plugin): add hot-reload capability for plugins
```

#### PR Description Template

```markdown
## Description
Brief description of the changes and motivation.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes

## Screenshots (if applicable)

## Related Issues
Closes #123
```

### Review Process

1. **Automated Checks**: All PRs must pass automated tests and linting
2. **Code Review**: At least one maintainer must review and approve
3. **Documentation Review**: Documentation changes require tech writer review
4. **Security Review**: Security-related changes require security team review

## 📚 Documentation Standards

### General Guidelines

- Write clear, concise documentation
- Use examples liberally
- Keep documentation up-to-date with code changes
- Follow the established documentation structure

### API Documentation

```go
// GetPlugin retrieves a plugin by name from the registry.
//
// Parameters:
//   - name: The unique identifier of the plugin to retrieve
//
// Returns:
//   - *Plugin: The plugin instance if found
//   - error: ErrPluginNotFound if the plugin doesn't exist, or other errors
//
// Example:
//   plugin, err := manager.GetPlugin("superclaude")
//   if err != nil {
//       log.Fatal(err)
//   }
//   result, err := plugin.Execute(ctx, []string{"analyze", "code.go"})
func (m *PluginManager) GetPlugin(name string) (*Plugin, error) {
    // Implementation
}
```

### README Updates

When adding new features:

1. Update relevant README sections
2. Add usage examples
3. Update the feature matrix
4. Include migration notes if needed

## 🧪 Testing Requirements

### Go Testing

#### Unit Tests

```go
func TestPluginManager_GetPlugin(t *testing.T) {
    tests := []struct {
        name     string
        pluginID string
        want     *Plugin
        wantErr  bool
    }{
        {
            name:     "existing plugin",
            pluginID: "superclaude",
            want:     &Plugin{Name: "superclaude"},
            wantErr:  false,
        },
        {
            name:     "non-existent plugin",
            pluginID: "nonexistent",
            want:     nil,
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := &PluginManager{/* setup */}
            got, err := m.GetPlugin(tt.pluginID)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetPlugin() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetPlugin() got = %v, want %v", got, tt.want)
            }
        })
    }
}
```

#### Integration Tests

```go
func TestPluginIntegration(t *testing.T) {
    // Set up test environment
    ctx := context.Background()
    manager := setupTestPluginManager(t)
    
    // Load test plugin
    err := manager.LoadPlugin("testdata/hello-plugin")
    require.NoError(t, err)
    
    // Execute plugin
    result, err := manager.ExecutePlugin(ctx, "hello", []string{"world"})
    require.NoError(t, err)
    assert.Equal(t, "Hello, world!", result.Output)
}
```

### TypeScript Testing

#### Unit Tests (Jest)

```typescript
describe('PluginManager', () => {
  let pluginManager: PluginManager;
  
  beforeEach(() => {
    pluginManager = new PluginManager();
  });
  
  test('should load plugin successfully', async () => {
    const config: PluginConfig = {
      name: 'test-plugin',
      version: '1.0.0',
      commands: [{ name: 'hello', description: 'Say hello' }]
    };
    
    await expect(pluginManager.loadPlugin(config)).resolves.not.toThrow();
    
    const loadedPlugin = pluginManager.getPlugin('test-plugin');
    expect(loadedPlugin).toBeDefined();
    expect(loadedPlugin?.config.name).toBe('test-plugin');
  });
});
```

### Test Coverage Requirements

- **Unit Tests**: Minimum 80% code coverage
- **Integration Tests**: All major workflows covered
- **E2E Tests**: Critical user journeys tested

### Running Tests

```bash
# Go tests
go test ./... -v -race -coverprofile=coverage.out
go tool cover -html=coverage.out

# TypeScript tests
npm run test
npm run test:coverage

# Integration tests
npm run test:integration

# All tests
make test
```

## 🔍 Review Process

### Review Criteria

#### Code Quality
- [ ] Code follows established patterns and conventions
- [ ] Appropriate error handling and logging
- [ ] Performance considerations addressed
- [ ] Security best practices followed

#### Testing
- [ ] Adequate test coverage (>80%)
- [ ] Tests are meaningful and test the right things
- [ ] Edge cases are covered
- [ ] Integration tests for new features

#### Documentation
- [ ] Public APIs are documented
- [ ] README updates if needed
- [ ] Code comments explain "why" not "what"
- [ ] Migration guides for breaking changes

#### Architecture
- [ ] Changes align with overall architecture
- [ ] Dependencies are justified and minimal
- [ ] Interfaces are well-defined
- [ ] Backwards compatibility maintained

### Review Assignments

- **Core maintainers**: Overall code review and architecture decisions
- **Domain experts**: Specialized reviews (security, performance, etc.)
- **Community reviewers**: General code quality and documentation

### Feedback and Iteration

- Be constructive and specific in feedback
- Explain the reasoning behind suggestions
- Be open to discussion and alternative approaches
- Focus on the code, not the person

## 🚀 Release Process

### Release Preparation

1. **Version Bump**: Update version numbers in relevant files
2. **Changelog**: Update CHANGELOG.md with new features and fixes
3. **Documentation**: Ensure all documentation is up-to-date
4. **Testing**: Run full test suite including integration tests

### Release Types

- **Major**: Breaking changes (v1.0.0 → v2.0.0)
- **Minor**: New features, backwards compatible (v1.0.0 → v1.1.0)
- **Patch**: Bug fixes, backwards compatible (v1.0.0 → v1.0.1)

## 🙋 Getting Help

### Resources

- **Documentation**: Check existing documentation first
- **Issues**: Search existing issues before creating new ones
- **Discussions**: Use GitHub Discussions for questions and ideas
- **FAQ**: Check [FAQ.md](../guides_and_documentation/FAQ.md) for common questions

### Contact

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and community discussion
- **Security Issues**: See [SECURITY.md](../guides_and_documentation/SECURITY.md) for reporting security vulnerabilities

## 📄 License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project (MIT License).

---

Thank you for contributing to the OpenCode-SuperClaude Integration project! Your contributions help make this tool better for everyone. 🎉