# 🧪 OpenCode SuperClaude Testing Strategy

## 📋 Overview

This comprehensive testing strategy provides detailed guidelines for ensuring the reliability, performance, and security of OpenCode SuperClaude integrations. It covers all aspects from unit testing to end-to-end validation, with specific approaches for each integration architecture.

## 🎯 Testing Philosophy

**Core Principles**:
- **Test-Driven Development**: Write tests before implementation
- **Pyramid Strategy**: Heavy unit tests, moderate integration tests, selective E2E tests
- **Risk-Based Testing**: Focus testing efforts on highest-risk components
- **Continuous Testing**: Automated testing in CI/CD pipeline
- **Performance-First**: Performance testing integrated throughout development

**Quality Standards**:
- **Unit Test Coverage**: ≥80% code coverage
- **Integration Test Coverage**: ≥70% critical path coverage
- **E2E Test Coverage**: 100% user journey coverage
- **Performance Targets**: <200ms API response, <3s UI load time
- **Reliability Target**: 99.9% uptime, <0.1% error rate

---

## 🏗️ Testing Architecture

### Testing Pyramid Structure

```
        E2E Tests (5%)
    ┌─────────────────────┐
    │  User Journeys      │
    │  Cross-Browser      │
    │  Performance        │
    └─────────────────────┘
            │
    Integration Tests (25%)
    ┌─────────────────────┐
    │  API Integration    │
    │  Plugin System      │
    │  Event Handling     │
    └─────────────────────┘
            │
        Unit Tests (70%)
    ┌─────────────────────┐
    │  Component Logic    │
    │  Business Rules     │
    │  Utility Functions  │
    └─────────────────────┘
```

### Test Environment Strategy

```typescript
// testing/config/environments.ts
export const TestEnvironments = {
  unit: {
    description: 'Isolated component testing',
    mocking: 'comprehensive',
    database: 'in-memory',
    network: 'mocked',
    isolation: 'complete'
  },
  
  integration: {
    description: 'Component interaction testing',
    mocking: 'external-services-only',
    database: 'test-database',
    network: 'stubbed',
    isolation: 'service-level'
  },
  
  staging: {
    description: 'Production-like environment',
    mocking: 'minimal',
    database: 'staging-database',
    network: 'real-with-stubs',
    isolation: 'environment-level'
  },
  
  production: {
    description: 'Live system validation',
    mocking: 'none',
    database: 'production',
    network: 'real',
    isolation: 'monitoring-only'
  }
};
```

---

## 🔬 Unit Testing Guidelines

### TypeScript Unit Testing Framework

#### Test Structure & Organization
```typescript
// testing/unit/template-broker.test.ts
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { TemplateBroker } from '../src/broker/template-broker';
import { MockFileWatcher } from './mocks/file-watcher.mock';
import { MockTemplateParser } from './mocks/template-parser.mock';

describe('TemplateBroker', () => {
  let broker: TemplateBroker;
  let mockWatcher: MockFileWatcher;
  let mockParser: MockTemplateParser;

  beforeEach(() => {
    mockWatcher = new MockFileWatcher();
    mockParser = new MockTemplateParser();
    
    broker = new TemplateBroker({
      watcher: mockWatcher,
      parser: mockParser,
      config: {
        watchPaths: ['/test/templates'],
        debounceMs: 100,
        maxTemplateSize: '1MB'
      }
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Template Loading', () => {
    it('should load valid templates successfully', async () => {
      // Arrange
      const templateContent = `---
name: test-template
category: code
---
This is a test template with {{variable}}`;

      mockWatcher.mockFileContent('/test/templates/test.md', templateContent);
      mockParser.mockParseResult({
        frontmatter: { name: 'test-template', category: 'code' },
        content: 'This is a test template with {{variable}}'
      });

      // Act
      await broker.loadTemplate('/test/templates/test.md');

      // Assert
      const template = broker.getTemplate('test-template');
      expect(template).toBeDefined();
      expect(template.name).toBe('test-template');
      expect(template.category).toBe('code');
      expect(template.content).toContain('{{variable}}');
    });

    it('should reject templates exceeding size limits', async () => {
      // Arrange
      const largeContent = 'x'.repeat(2 * 1024 * 1024); // 2MB content
      mockWatcher.mockFileContent('/test/templates/large.md', largeContent);

      // Act & Assert
      await expect(broker.loadTemplate('/test/templates/large.md'))
        .rejects
        .toThrow('Template size exceeds maximum allowed size');
    });

    it('should handle malformed frontmatter gracefully', async () => {
      // Arrange
      const malformedContent = `---
name: test
invalid-yaml: [unclosed
---
Template content`;

      mockWatcher.mockFileContent('/test/templates/malformed.md', malformedContent);
      mockParser.mockParseError(new Error('Invalid YAML syntax'));

      // Act & Assert
      await expect(broker.loadTemplate('/test/templates/malformed.md'))
        .rejects
        .toThrow('Failed to parse template frontmatter');
    });
  });

  describe('Template Rendering', () => {
    it('should render templates with context variables', async () => {
      // Arrange
      const template = {
        name: 'greeting',
        content: 'Hello {{name}}, your role is {{role}}.',
        frontmatter: { category: 'general' }
      };

      broker.addTemplate(template);

      const context = { name: 'Alice', role: 'developer' };

      // Act
      const rendered = await broker.renderTemplate('greeting', context);

      // Assert
      expect(rendered).toBe('Hello Alice, your role is developer.');
    });

    it('should handle missing context variables', async () => {
      // Arrange
      const template = {
        name: 'incomplete',
        content: 'Hello {{name}}, your role is {{role}}.',
        frontmatter: { strictMode: true }
      };

      broker.addTemplate(template);
      const incompleteContext = { name: 'Bob' }; // missing 'role'

      // Act & Assert
      await expect(broker.renderTemplate('incomplete', incompleteContext))
        .rejects
        .toThrow('Missing required template variable: role');
    });
  });

  describe('Performance Requirements', () => {
    it('should load templates within performance targets', async () => {
      // Arrange
      const startTime = performance.now();
      const templateCount = 100;
      
      // Create mock templates
      for (let i = 0; i < templateCount; i++) {
        mockWatcher.mockFileContent(
          `/test/templates/template-${i}.md`,
          `---\nname: template-${i}\n---\nContent ${i}`
        );
      }

      // Act
      const loadPromises = Array.from({ length: templateCount }, (_, i) =>
        broker.loadTemplate(`/test/templates/template-${i}.md`)
      );
      
      await Promise.all(loadPromises);
      const endTime = performance.now();

      // Assert
      const loadTime = endTime - startTime;
      expect(loadTime).toBeLessThan(1000); // Should load 100 templates in <1s
      
      // Verify all templates loaded
      expect(broker.getTemplateCount()).toBe(templateCount);
    });

    it('should render templates within performance targets', async () => {
      // Arrange
      const template = {
        name: 'performance-test',
        content: 'Template with {{var1}} and {{var2}} and {{var3}}',
        frontmatter: {}
      };
      
      broker.addTemplate(template);
      const context = { var1: 'value1', var2: 'value2', var3: 'value3' };

      // Act - render multiple times to test performance
      const iterations = 1000;
      const startTime = performance.now();
      
      for (let i = 0; i < iterations; i++) {
        await broker.renderTemplate('performance-test', context);
      }
      
      const endTime = performance.now();

      // Assert
      const avgRenderTime = (endTime - startTime) / iterations;
      expect(avgRenderTime).toBeLessThan(10); // <10ms average render time
    });
  });
});
```

#### Mock Utilities & Test Helpers
```typescript
// testing/mocks/file-watcher.mock.ts
export class MockFileWatcher {
  private fileContents = new Map<string, string>();
  private watchCallbacks = new Map<string, Function[]>();

  mockFileContent(path: string, content: string): void {
    this.fileContents.set(path, content);
  }

  mockFileChange(path: string, content: string): void {
    this.fileContents.set(path, content);
    const callbacks = this.watchCallbacks.get(path) || [];
    callbacks.forEach(callback => callback('change', path));
  }

  watch(path: string, callback: Function): void {
    if (!this.watchCallbacks.has(path)) {
      this.watchCallbacks.set(path, []);
    }
    this.watchCallbacks.get(path)!.push(callback);
  }

  readFile(path: string): Promise<string> {
    const content = this.fileContents.get(path);
    if (!content) {
      return Promise.reject(new Error(`File not found: ${path}`));
    }
    return Promise.resolve(content);
  }
}

// testing/helpers/test-helpers.ts
export class TestHelpers {
  static createMockTemplate(overrides: Partial<Template> = {}): Template {
    return {
      id: 'test-template',
      name: 'Test Template',
      category: 'code',
      content: 'Default test content with {{variable}}',
      frontmatter: {
        name: 'Test Template',
        category: 'code',
        version: '1.0.0'
      },
      createdAt: new Date(),
      updatedAt: new Date(),
      ...overrides
    };
  }

  static createMockContext(overrides: Record<string, any> = {}): TemplateContext {
    return {
      selectedCode: 'console.log("test");',
      filePath: '/test/file.ts',
      projectContext: { name: 'test-project', type: 'typescript' },
      userInput: 'Test user input',
      ...overrides
    };
  }

  static async waitFor(condition: () => boolean, timeout: number = 5000): Promise<void> {
    const start = Date.now();
    while (!condition() && (Date.now() - start) < timeout) {
      await new Promise(resolve => setTimeout(resolve, 10));
    }
    if (!condition()) {
      throw new Error(`Condition not met within ${timeout}ms`);
    }
  }
}
```

### Go Unit Testing Framework

#### Go Test Structure
```go
// internal/plugin/manager_test.go
package plugin

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
)

type PluginManagerTestSuite struct {
    suite.Suite
    manager    *PluginManager
    mockPlugin *MockPlugin
    ctx        context.Context
    cancel     context.CancelFunc
}

func (suite *PluginManagerTestSuite) SetupTest() {
    suite.ctx, suite.cancel = context.WithCancel(context.Background())
    suite.mockPlugin = NewMockPlugin(suite.T())
    
    config := &PluginManagerConfig{
        PluginDir:     "/tmp/test-plugins",
        MaxPlugins:    5,
        HealthTimeout: 5 * time.Second,
    }
    
    suite.manager = NewPluginManager(config)
}

func (suite *PluginManagerTestSuite) TearDownTest() {
    if suite.cancel != nil {
        suite.cancel()
    }
    if suite.manager != nil {
        suite.manager.Shutdown()
    }
}

func (suite *PluginManagerTestSuite) TestLoadPlugin() {
    // Arrange
    pluginPath := "/tmp/test-plugins/test-plugin"
    expectedMetadata := &PluginMetadata{
        Name:    "test-plugin",
        Version: "1.0.0",
        Commands: []string{"test-command"},
    }

    suite.mockPlugin.On("GetMetadata").Return(expectedMetadata, nil)
    suite.mockPlugin.On("Start", mock.Anything).Return(nil)
    suite.mockPlugin.On("Health").Return(&HealthStatus{Status: "healthy"}, nil)

    // Act
    plugin, err := suite.manager.LoadPlugin(suite.ctx, pluginPath)

    // Assert
    require.NoError(suite.T(), err)
    assert.NotNil(suite.T(), plugin)
    assert.Equal(suite.T(), "test-plugin", plugin.Name())
    
    suite.mockPlugin.AssertExpectations(suite.T())
}

func (suite *PluginManagerTestSuite) TestLoadPluginTimeout() {
    // Arrange
    pluginPath := "/tmp/test-plugins/slow-plugin"
    
    suite.mockPlugin.On("Start", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
        // Simulate slow plugin startup
        time.Sleep(10 * time.Second)
    })

    ctx, cancel := context.WithTimeout(suite.ctx, 1*time.Second)
    defer cancel()

    // Act
    plugin, err := suite.manager.LoadPlugin(ctx, pluginPath)

    // Assert
    assert.Error(suite.T(), err)
    assert.Nil(suite.T(), plugin)
    assert.Contains(suite.T(), err.Error(), "timeout")
}

func (suite *PluginManagerTestSuite) TestPluginHealthMonitoring() {
    // Arrange
    suite.mockPlugin.On("GetMetadata").Return(&PluginMetadata{
        Name: "health-test-plugin",
    }, nil)
    suite.mockPlugin.On("Start", mock.Anything).Return(nil)
    
    // First health check - healthy
    suite.mockPlugin.On("Health").Return(&HealthStatus{
        Status: "healthy",
    }, nil).Once()
    
    // Second health check - unhealthy
    suite.mockPlugin.On("Health").Return(&HealthStatus{
        Status: "unhealthy",
        Error:  "plugin crashed",
    }, nil).Once()
    
    suite.mockPlugin.On("Stop").Return(nil)

    // Act
    plugin, err := suite.manager.LoadPlugin(suite.ctx, "/tmp/test-plugins/health-test")
    require.NoError(suite.T(), err)

    // Start health monitoring
    suite.manager.StartHealthMonitoring(1*time.Second)

    // Wait for health check cycles
    time.Sleep(3 * time.Second)

    // Assert
    assert.False(suite.T(), suite.manager.IsPluginHealthy(plugin.Name()))
    suite.mockPlugin.AssertExpectations(suite.T())
}

func (suite *PluginManagerTestSuite) TestConcurrentPluginLoading() {
    // Arrange
    pluginCount := 5
    plugins := make([]*MockPlugin, pluginCount)
    
    for i := 0; i < pluginCount; i++ {
        plugins[i] = NewMockPlugin(suite.T())
        plugins[i].On("GetMetadata").Return(&PluginMetadata{
            Name: fmt.Sprintf("plugin-%d", i),
        }, nil)
        plugins[i].On("Start", mock.Anything).Return(nil)
        plugins[i].On("Health").Return(&HealthStatus{Status: "healthy"}, nil)
    }

    // Act
    var wg sync.WaitGroup
    results := make(chan error, pluginCount)
    
    for i := 0; i < pluginCount; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()
            _, err := suite.manager.LoadPlugin(suite.ctx, fmt.Sprintf("/tmp/test-plugins/plugin-%d", index))
            results <- err
        }(i)
    }

    wg.Wait()
    close(results)

    // Assert
    errorCount := 0
    for err := range results {
        if err != nil {
            errorCount++
        }
    }
    
    assert.Equal(suite.T(), 0, errorCount, "All plugins should load successfully")
    assert.Equal(suite.T(), pluginCount, suite.manager.GetPluginCount())
}

func TestPluginManagerTestSuite(t *testing.T) {
    suite.Run(t, new(PluginManagerTestSuite))
}

// Benchmark tests for performance validation
func BenchmarkPluginManager_LoadPlugin(b *testing.B) {
    config := &PluginManagerConfig{
        PluginDir:     "/tmp/bench-plugins",
        MaxPlugins:    100,
        HealthTimeout: 5 * time.Second,
    }
    
    manager := NewPluginManager(config)
    defer manager.Shutdown()

    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        ctx := context.Background()
        pluginPath := fmt.Sprintf("/tmp/bench-plugins/plugin-%d", i%10)
        
        _, err := manager.LoadPlugin(ctx, pluginPath)
        if err != nil {
            b.Fatalf("Failed to load plugin: %v", err)
        }
    }
}

func BenchmarkPluginManager_ExecuteCommand(b *testing.B) {
    // Setup
    manager := setupBenchmarkManager(b)
    defer manager.Shutdown()

    request := &CommandRequest{
        Command: "test-command",
        Args:    map[string]interface{}{"input": "test"},
    }

    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        ctx := context.Background()
        _, err := manager.ExecuteCommand(ctx, "test-plugin", request)
        if err != nil {
            b.Fatalf("Failed to execute command: %v", err)
        }
    }
}
```

#### Go Mock Generation
```go
// testing/mocks/generate.go
//go:generate mockery --name=CommandPlugin --output=../mocks --outpkg=mocks
//go:generate mockery --name=EventBus --output=../mocks --outpkg=mocks
//go:generate mockery --name=HealthChecker --output=../mocks --outpkg=mocks

// testing/mocks/plugin_mock.go
// Generated by mockery. Manual modifications will be overwritten.
type MockPlugin struct {
    mock.Mock
}

func (m *MockPlugin) Name() string {
    ret := m.Called()
    return ret.String(0)
}

func (m *MockPlugin) Execute(ctx context.Context, req *CommandRequest) (*CommandResponse, error) {
    ret := m.Called(ctx, req)
    return ret.Get(0).(*CommandResponse), ret.Error(1)
}

func (m *MockPlugin) Health() (*HealthStatus, error) {
    ret := m.Called()
    return ret.Get(0).(*HealthStatus), ret.Error(1)
}
```

---

## 🔗 Integration Testing

### API Integration Testing

#### TypeScript API Integration Tests
```typescript
// testing/integration/api-integration.test.ts
import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import { TestServer } from '../helpers/test-server';
import { TestDatabase } from '../helpers/test-database';
import { SuperClaudeClient } from '../clients/superclaude-client';

describe('SuperClaude API Integration', () => {
  let server: TestServer;
  let database: TestDatabase;
  let client: SuperClaudeClient;

  beforeAll(async () => {
    // Setup test environment
    database = new TestDatabase();
    await database.setup();
    
    server = new TestServer({
      database: database.getConnection(),
      port: 0, // Random available port
      logLevel: 'error' // Reduce noise in tests
    });
    
    await server.start();
    
    client = new SuperClaudeClient({
      baseUrl: server.getUrl(),
      timeout: 10000
    });
  });

  afterAll(async () => {
    await server.stop();
    await database.cleanup();
  });

  describe('Template Management API', () => {
    it('should create, retrieve, update, and delete templates', async () => {
      // Create template
      const createRequest = {
        name: 'integration-test-template',
        category: 'code',
        content: 'Test template with {{variable}}',
        frontmatter: {
          requires: ['selectedCode'],
          optional: { userInput: true }
        }
      };

      const createResponse = await client.createTemplate(createRequest);
      expect(createResponse.status).toBe(201);
      expect(createResponse.data.id).toBeDefined();
      
      const templateId = createResponse.data.id;

      // Retrieve template
      const getResponse = await client.getTemplate(templateId);
      expect(getResponse.status).toBe(200);
      expect(getResponse.data.name).toBe(createRequest.name);
      expect(getResponse.data.content).toBe(createRequest.content);

      // Update template
      const updateRequest = {
        ...createRequest,
        content: 'Updated test template with {{variable}}'
      };

      const updateResponse = await client.updateTemplate(templateId, updateRequest);
      expect(updateResponse.status).toBe(200);
      expect(updateResponse.data.content).toBe(updateRequest.content);

      // Delete template
      const deleteResponse = await client.deleteTemplate(templateId);
      expect(deleteResponse.status).toBe(204);

      // Verify deletion
      await expect(client.getTemplate(templateId))
        .rejects
        .toThrow('Template not found');
    });

    it('should handle concurrent template operations', async () => {
      const concurrency = 10;
      const templates = Array.from({ length: concurrency }, (_, i) => ({
        name: `concurrent-template-${i}`,
        category: 'code',
        content: `Template ${i} with {{variable}}`,
        frontmatter: {}
      }));

      // Create templates concurrently
      const createPromises = templates.map(template => 
        client.createTemplate(template)
      );

      const createResults = await Promise.all(createPromises);
      
      // Verify all created successfully
      createResults.forEach((result, i) => {
        expect(result.status).toBe(201);
        expect(result.data.name).toBe(templates[i].name);
      });

      // Cleanup
      const deletePromises = createResults.map(result =>
        client.deleteTemplate(result.data.id)
      );
      
      await Promise.all(deletePromises);
    });
  });

  describe('Template Rendering API', () => {
    it('should render templates with context', async () => {
      // Create template
      const template = {
        name: 'render-test-template',
        category: 'code',
        content: 'Hello {{name}}, please {{action}} the {{target}}.',
        frontmatter: {
          requires: ['name', 'action', 'target']
        }
      };

      const createResponse = await client.createTemplate(template);
      const templateId = createResponse.data.id;

      // Render template
      const context = {
        name: 'Developer',
        action: 'refactor',
        target: 'function'
      };

      const renderResponse = await client.renderTemplate(templateId, context);
      
      expect(renderResponse.status).toBe(200);
      expect(renderResponse.data.rendered).toBe(
        'Hello Developer, please refactor the function.'
      );
      expect(renderResponse.data.metadata.renderTime).toBeLessThan(100);

      // Cleanup
      await client.deleteTemplate(templateId);
    });

    it('should validate required context variables', async () => {
      // Create template with required variables
      const template = {
        name: 'validation-test-template',
        category: 'code',
        content: 'Required: {{required1}} and {{required2}}',
        frontmatter: {
          requires: ['required1', 'required2'],
          strictMode: true
        }
      };

      const createResponse = await client.createTemplate(template);
      const templateId = createResponse.data.id;

      // Attempt render with missing variables
      const incompleteContext = { required1: 'value1' };

      await expect(client.renderTemplate(templateId, incompleteContext))
        .rejects
        .toThrow('Missing required variable: required2');

      // Cleanup
      await client.deleteTemplate(templateId);
    });
  });

  describe('File Watcher Integration', () => {
    it('should detect and reload template changes', async () => {
      const testFilePath = '/tmp/integration-test-template.md';
      const initialContent = `---
name: file-watcher-test
category: code
---
Initial content with {{variable}}`;

      // Write initial template file
      await fs.promises.writeFile(testFilePath, initialContent);

      // Wait for file watcher to detect and load
      await TestHelpers.waitFor(
        () => client.getTemplateByName('file-watcher-test') !== null,
        5000
      );

      // Verify initial load
      const initialTemplate = await client.getTemplateByName('file-watcher-test');
      expect(initialTemplate.content).toContain('Initial content');

      // Update file
      const updatedContent = `---
name: file-watcher-test
category: code
---
Updated content with {{variable}}`;

      await fs.promises.writeFile(testFilePath, updatedContent);

      // Wait for reload
      await TestHelpers.waitFor(
        async () => {
          const template = await client.getTemplateByName('file-watcher-test');
          return template?.content?.includes('Updated content') || false;
        },
        5000
      );

      // Verify update
      const updatedTemplate = await client.getTemplateByName('file-watcher-test');
      expect(updatedTemplate.content).toContain('Updated content');

      // Cleanup
      await fs.promises.unlink(testFilePath);
    });
  });

  describe('Performance Integration', () => {
    it('should handle high load template operations', async () => {
      const requestCount = 1000;
      const maxConcurrency = 50;
      
      // Create test template
      const template = {
        name: 'performance-test-template',
        category: 'code',
        content: 'Performance test with {{iteration}}',
        frontmatter: {}
      };

      const createResponse = await client.createTemplate(template);
      const templateId = createResponse.data.id;

      // Generate render requests
      const renderRequests = Array.from({ length: requestCount }, (_, i) => ({
        templateId,
        context: { iteration: i }
      }));

      // Execute with concurrency limit
      const startTime = performance.now();
      
      let results: any[] = [];
      for (let i = 0; i < renderRequests.length; i += maxConcurrency) {
        const batch = renderRequests.slice(i, i + maxConcurrency);
        const batchPromises = batch.map(req => 
          client.renderTemplate(req.templateId, req.context)
        );
        
        const batchResults = await Promise.all(batchPromises);
        results = results.concat(batchResults);
      }

      const endTime = performance.now();
      const totalTime = endTime - startTime;

      // Performance assertions
      expect(results).toHaveLength(requestCount);
      expect(totalTime).toBeLessThan(30000); // Should complete in <30s
      
      const avgResponseTime = totalTime / requestCount;
      expect(avgResponseTime).toBeLessThan(30); // <30ms average

      // Verify all renders successful
      results.forEach((result, i) => {
        expect(result.status).toBe(200);
        expect(result.data.rendered).toContain(`${i}`);
      });

      // Cleanup
      await client.deleteTemplate(templateId);
    });
  });
});
```

### Plugin System Integration Testing

#### Go Plugin Integration Tests
```go
// testing/integration/plugin_integration_test.go
package integration

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
    
    "github.com/opencode/superclaude/internal/plugin"
    "github.com/opencode/superclaude/testing/helpers"
)

type PluginIntegrationSuite struct {
    suite.Suite
    tempDir    string
    manager    *plugin.PluginManager
    ctx        context.Context
    cancel     context.CancelFunc
}

func (suite *PluginIntegrationSuite) SetupSuite() {
    var err error
    suite.tempDir, err = os.MkdirTemp("", "plugin-integration-test")
    require.NoError(suite.T(), err)

    // Build test plugins
    err = suite.buildTestPlugins()
    require.NoError(suite.T(), err)
}

func (suite *PluginIntegrationSuite) TearDownSuite() {
    os.RemoveAll(suite.tempDir)
}

func (suite *PluginIntegrationSuite) SetupTest() {
    suite.ctx, suite.cancel = context.WithCancel(context.Background())
    
    config := &plugin.PluginManagerConfig{
        PluginDir:     filepath.Join(suite.tempDir, "plugins"),
        MaxPlugins:    10,
        HealthTimeout: 5 * time.Second,
        EventBus:      helpers.NewTestEventBus(),
    }
    
    suite.manager = plugin.NewPluginManager(config)
    require.NoError(suite.T(), suite.manager.Start(suite.ctx))
}

func (suite *PluginIntegrationSuite) TearDownTest() {
    if suite.manager != nil {
        suite.manager.Shutdown()
    }
    if suite.cancel != nil {
        suite.cancel()
    }
}

func (suite *PluginIntegrationSuite) TestPluginDiscoveryAndLoading() {
    // Act
    err := suite.manager.DiscoverAndLoadPlugins()
    require.NoError(suite.T(), err)

    // Wait for plugins to load
    time.Sleep(2 * time.Second)

    // Assert
    plugins := suite.manager.GetLoadedPlugins()
    assert.GreaterOrEqual(suite.T(), len(plugins), 1, "Should load at least one test plugin")

    // Verify specific test plugin
    testPlugin := suite.manager.GetPlugin("test-plugin")
    assert.NotNil(suite.T(), testPlugin)
    assert.Equal(suite.T(), "test-plugin", testPlugin.Name())
}

func (suite *PluginIntegrationSuite) TestPluginCommandExecution() {
    // Arrange
    err := suite.manager.DiscoverAndLoadPlugins()
    require.NoError(suite.T(), err)

    time.Sleep(1 * time.Second) // Wait for plugin to load

    request := &plugin.CommandRequest{
        Command: "echo",
        Args: map[string]interface{}{
            "message": "Hello from integration test",
        },
    }

    // Act
    response, err := suite.manager.ExecuteCommand(suite.ctx, "test-plugin", request)

    // Assert
    require.NoError(suite.T(), err)
    assert.NotNil(suite.T(), response)
    assert.Equal(suite.T(), "success", response.Status)
    assert.Equal(suite.T(), "Hello from integration test", response.Data["echo"])
}

func (suite *PluginIntegrationSuite) TestPluginHealthMonitoring() {
    // Arrange
    err := suite.manager.DiscoverAndLoadPlugins()
    require.NoError(suite.T(), err)

    time.Sleep(1 * time.Second)

    // Start health monitoring
    suite.manager.StartHealthMonitoring(1 * time.Second)

    // Act - wait for a few health check cycles
    time.Sleep(3 * time.Second)

    // Assert
    plugins := suite.manager.GetLoadedPlugins()
    for _, pluginName := range plugins {
        isHealthy := suite.manager.IsPluginHealthy(pluginName)
        assert.True(suite.T(), isHealthy, fmt.Sprintf("Plugin %s should be healthy", pluginName))
    }
}

func (suite *PluginIntegrationSuite) TestPluginCrashRecovery() {
    // Arrange
    err := suite.manager.DiscoverAndLoadPlugins()
    require.NoError(suite.T(), err)

    time.Sleep(1 * time.Second)

    // Act - send crash command to test plugin
    crashRequest := &plugin.CommandRequest{
        Command: "crash",
        Args:    map[string]interface{}{},
    }

    _, err = suite.manager.ExecuteCommand(suite.ctx, "test-plugin", crashRequest)
    // Plugin should crash, so we expect an error
    assert.Error(suite.T(), err)

    // Wait for recovery
    time.Sleep(3 * time.Second)

    // Assert - plugin should be restarted and healthy
    echoRequest := &plugin.CommandRequest{
        Command: "echo",
        Args: map[string]interface{}{
            "message": "Recovery test",
        },
    }

    response, err := suite.manager.ExecuteCommand(suite.ctx, "test-plugin", echoRequest)
    require.NoError(suite.T(), err)
    assert.Equal(suite.T(), "Recovery test", response.Data["echo"])
}

func (suite *PluginIntegrationSuite) TestPluginEventHandling() {
    // Arrange
    err := suite.manager.DiscoverAndLoadPlugins()
    require.NoError(suite.T(), err)

    time.Sleep(1 * time.Second)

    eventBus := suite.manager.GetEventBus()
    
    // Subscribe to plugin events
    eventChan := make(chan *plugin.Event, 10)
    eventBus.Subscribe("plugin.*", func(event *plugin.Event) {
        eventChan <- event
    })

    // Act - trigger plugin reload
    err = suite.manager.ReloadPlugin("test-plugin")
    require.NoError(suite.T(), err)

    // Assert - should receive reload events
    var events []*plugin.Event
    timeout := time.After(5 * time.Second)
    
    for len(events) < 2 { // Expect plugin.stopping and plugin.started events
        select {
        case event := <-eventChan:
            events = append(events, event)
        case <-timeout:
            break
        }
    }

    assert.GreaterOrEqual(suite.T(), len(events), 2)
    
    // Check for expected event types
    eventTypes := make(map[string]bool)
    for _, event := range events {
        eventTypes[event.Type] = true
    }
    
    assert.True(suite.T(), eventTypes["plugin.stopping"] || eventTypes["plugin.stopped"])
    assert.True(suite.T(), eventTypes["plugin.started"])
}

func (suite *PluginIntegrationSuite) TestConcurrentPluginOperations() {
    // Arrange
    err := suite.manager.DiscoverAndLoadPlugins()
    require.NoError(suite.T(), err)

    time.Sleep(1 * time.Second)

    concurrency := 20
    requestsPerGoroutine := 10

    // Act - concurrent command execution
    results := make(chan error, concurrency*requestsPerGoroutine)
    
    for i := 0; i < concurrency; i++ {
        go func(goroutineID int) {
            for j := 0; j < requestsPerGoroutine; j++ {
                request := &plugin.CommandRequest{
                    Command: "echo",
                    Args: map[string]interface{}{
                        "message": fmt.Sprintf("Goroutine %d, Request %d", goroutineID, j),
                    },
                }

                _, err := suite.manager.ExecuteCommand(suite.ctx, "test-plugin", request)
                results <- err
            }
        }(i)
    }

    // Collect results
    successCount := 0
    totalRequests := concurrency * requestsPerGoroutine
    
    for i := 0; i < totalRequests; i++ {
        err := <-results
        if err == nil {
            successCount++
        }
    }

    // Assert
    successRate := float64(successCount) / float64(totalRequests)
    assert.GreaterOrEqual(suite.T(), successRate, 0.95, "Success rate should be at least 95%")
}

func (suite *PluginIntegrationSuite) buildTestPlugins() error {
    // Build test plugin binary
    pluginDir := filepath.Join(suite.tempDir, "plugins")
    err := os.MkdirAll(pluginDir, 0755)
    if err != nil {
        return err
    }

    // Copy test plugin source and build
    testPluginSrc := filepath.Join("testdata", "test-plugin")
    testPluginBin := filepath.Join(pluginDir, "test-plugin")
    
    return helpers.BuildPlugin(testPluginSrc, testPluginBin)
}

func TestPluginIntegrationSuite(t *testing.T) {
    suite.Run(t, new(PluginIntegrationSuite))
}
```

### Event System Integration Testing

#### Event-Driven Architecture Tests
```go
// testing/integration/event_integration_test.go
package integration

import (
    "context"
    "testing"
    "time"

    "github.com/nats-io/nats.go"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"

    "github.com/opencode/superclaude/internal/events"
)

type EventIntegrationSuite struct {
    suite.Suite
    natsServer *nats.Server
    eventBus   *events.EventBus
    ctx        context.Context
    cancel     context.CancelFunc
}

func (suite *EventIntegrationSuite) SetupSuite() {
    // Start embedded NATS server for testing
    opts := &nats.Options{
        Port:      -1, // Random port
        JetStream: true,
    }

    var err error
    suite.natsServer, err = nats.NewServer(opts)
    require.NoError(suite.T(), err)

    go suite.natsServer.Start()
    
    // Wait for server to be ready
    if !suite.natsServer.ReadyForConnections(5 * time.Second) {
        suite.T().Fatal("NATS server not ready")
    }
}

func (suite *EventIntegrationSuite) TearDownSuite() {
    if suite.natsServer != nil {
        suite.natsServer.Shutdown()
    }
}

func (suite *EventIntegrationSuite) SetupTest() {
    suite.ctx, suite.cancel = context.WithCancel(context.Background())

    config := &events.Config{
        NatsURL: suite.natsServer.ClientURL(),
        StreamConfig: events.StreamConfig{
            Name:      "test-stream",
            Subjects:  []string{"test.*"},
            Retention: nats.WorkQueuePolicy,
        },
    }

    var err error
    suite.eventBus, err = events.NewEventBus(config)
    require.NoError(suite.T(), err)

    err = suite.eventBus.Start(suite.ctx)
    require.NoError(suite.T(), err)
}

func (suite *EventIntegrationSuite) TearDownTest() {
    if suite.eventBus != nil {
        suite.eventBus.Stop()
    }
    if suite.cancel != nil {
        suite.cancel()
    }
}

func (suite *EventIntegrationSuite) TestEventPublishAndSubscribe() {
    // Arrange
    eventChan := make(chan *events.Event, 10)
    
    err := suite.eventBus.Subscribe("test.example", func(event *events.Event) {
        eventChan <- event
    })
    require.NoError(suite.T(), err)

    testEvent := &events.Event{
        Type:      "test.example",
        Source:    "integration-test",
        Data:      map[string]interface{}{"message": "Hello World"},
        Timestamp: time.Now(),
    }

    // Act
    err = suite.eventBus.Publish(testEvent)
    require.NoError(suite.T(), err)

    // Assert
    select {
    case receivedEvent := <-eventChan:
        assert.Equal(suite.T(), testEvent.Type, receivedEvent.Type)
        assert.Equal(suite.T(), testEvent.Source, receivedEvent.Source)
        assert.Equal(suite.T(), testEvent.Data["message"], receivedEvent.Data["message"])
    case <-time.After(5 * time.Second):
        suite.T().Fatal("Timeout waiting for event")
    }
}

func (suite *EventIntegrationSuite) TestEventPersistence() {
    // Arrange
    testEvent := &events.Event{
        Type:   "test.persistence",
        Source: "integration-test",
        Data:   map[string]interface{}{"persistent": true},
    }

    // Act - publish event
    err := suite.eventBus.Publish(testEvent)
    require.NoError(suite.T(), err)

    // Stop and restart event bus
    suite.eventBus.Stop()
    
    // Create new event bus instance
    config := &events.Config{
        NatsURL: suite.natsServer.ClientURL(),
        StreamConfig: events.StreamConfig{
            Name:      "test-stream",
            Subjects:  []string{"test.*"},
            Retention: nats.WorkQueuePolicy,
        },
    }

    newEventBus, err := events.NewEventBus(config)
    require.NoError(suite.T(), err)

    err = newEventBus.Start(suite.ctx)
    require.NoError(suite.T(), err)

    // Subscribe to events
    eventChan := make(chan *events.Event, 10)
    err = newEventBus.Subscribe("test.persistence", func(event *events.Event) {
        eventChan <- event
    })
    require.NoError(suite.T(), err)

    // Assert - should receive previously published event
    select {
    case receivedEvent := <-eventChan:
        assert.Equal(suite.T(), testEvent.Type, receivedEvent.Type)
        assert.Equal(suite.T(), testEvent.Data["persistent"], receivedEvent.Data["persistent"])
    case <-time.After(5 * time.Second):
        suite.T().Fatal("Timeout waiting for persisted event")
    }

    newEventBus.Stop()
}

func (suite *EventIntegrationSuite) TestHighVolumeEventProcessing() {
    // Arrange
    eventCount := 10000
    receivedEvents := make(chan *events.Event, eventCount)
    
    err := suite.eventBus.Subscribe("test.volume.*", func(event *events.Event) {
        receivedEvents <- event
    })
    require.NoError(suite.T(), err)

    // Act - publish events concurrently
    startTime := time.Now()
    
    for i := 0; i < eventCount; i++ {
        go func(index int) {
            event := &events.Event{
                Type:   "test.volume.event",
                Source: "load-test",
                Data:   map[string]interface{}{"index": index},
            }
            suite.eventBus.Publish(event)
        }(i)
    }

    // Assert - collect all events
    receivedCount := 0
    timeout := time.After(30 * time.Second)

    for receivedCount < eventCount {
        select {
        case <-receivedEvents:
            receivedCount++
        case <-timeout:
            break
        }
    }

    endTime := time.Now()
    duration := endTime.Sub(startTime)

    assert.Equal(suite.T(), eventCount, receivedCount, "Should receive all published events")
    
    throughput := float64(eventCount) / duration.Seconds()
    assert.Greater(suite.T(), throughput, 1000.0, "Should process at least 1000 events/second")
}

func TestEventIntegrationSuite(t *testing.T) {
    suite.Run(t, new(EventIntegrationSuite))
}
```

---

## 🌐 End-to-End Testing with Playwright

### E2E Test Framework Setup

#### Playwright Configuration
```typescript
// playwright.config.ts
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './testing/e2e',
  outputDir: './testing/e2e-results',
  
  // Parallel execution
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,

  // Test settings
  timeout: 30 * 1000, // 30 seconds
  expect: { timeout: 5 * 1000 }, // 5 seconds
  
  use: {
    baseURL: process.env.BASE_URL || 'http://localhost:3000',
    trace: 'on-first-retry', // Collect trace on retry
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },

  // Test projects for different browsers
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
    {
      name: 'mobile-chrome',
      use: { ...devices['Pixel 5'] },
    },
    {
      name: 'mobile-safari',
      use: { ...devices['iPhone 12'] },
    },
  ],

  // Web server setup for tests
  webServer: {
    command: process.env.CI ? 'npm run start:test' : 'npm run dev',
    port: 3000,
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000, // 2 minutes
  },

  // Reporting
  reporter: [
    ['html', { outputFolder: './testing/e2e-reports' }],
    ['junit', { outputFile: './testing/e2e-results/junit.xml' }],
    ['github'], // GitHub Actions integration
  ],
});
```

#### E2E Test Structure
```typescript
// testing/e2e/superclaude-workflow.spec.ts
import { test, expect, Page } from '@playwright/test';
import { SuperClaudeE2EHelper } from '../helpers/superclaude-e2e-helper';

test.describe('SuperClaude Integration Workflow', () => {
  let helper: SuperClaudeE2EHelper;

  test.beforeEach(async ({ page }) => {
    helper = new SuperClaudeE2EHelper(page);
    await helper.setup();
  });

  test.afterEach(async () => {
    await helper.cleanup();
  });

  test('should complete full template creation and usage workflow', async ({ page }) => {
    // Step 1: Navigate to template management
    await test.step('Navigate to template management', async () => {
      await page.goto('/templates');
      await expect(page.locator('h1')).toContainText('Template Management');
    });

    // Step 2: Create new template
    await test.step('Create new template', async () => {
      await page.click('[data-testid="create-template-button"]');
      
      const modal = page.locator('[data-testid="template-modal"]');
      await expect(modal).toBeVisible();

      await page.fill('[data-testid="template-name"]', 'E2E Test Template');
      await page.selectOption('[data-testid="template-category"]', 'code');
      await page.fill('[data-testid="template-content"]', 
        'Please {{action}} the following {{target}}:\n\n{{selectedCode}}');
      
      await page.click('[data-testid="save-template"]');
      
      // Wait for success notification
      await expect(page.locator('[data-testid="success-notification"]'))
        .toContainText('Template created successfully');
    });

    // Step 3: Verify template appears in list
    await test.step('Verify template in list', async () => {
      await expect(page.locator('[data-testid="template-list"]'))
        .toContainText('E2E Test Template');
    });

    // Step 4: Use template in code editor
    await test.step('Use template in code editor', async () => {
      await page.goto('/editor');
      
      // Add some code to editor
      const editor = page.locator('[data-testid="code-editor"]');
      await editor.fill(`function testFunction() {
  console.log('This needs refactoring');
  return true;
}`);
      
      // Select the code
      await editor.selectText();
      
      // Open template selector
      await page.click('[data-testid="template-selector"]');
      
      // Select our created template
      await page.click('[data-testid="template-option"][data-template="E2E Test Template"]');
      
      // Fill template variables
      await page.fill('[data-testid="template-variable-action"]', 'refactor');
      await page.fill('[data-testid="template-variable-target"]', 'function');
      
      // Apply template
      await page.click('[data-testid="apply-template"]');
      
      // Verify rendered output
      const renderedOutput = page.locator('[data-testid="rendered-output"]');
      await expect(renderedOutput).toContainText('Please refactor the following function:');
      await expect(renderedOutput).toContainText('function testFunction()');
    });

    // Step 5: Test live reload functionality
    await test.step('Test live reload functionality', async () => {
      // Modify template file directly (simulating external edit)
      await helper.modifyTemplateFile('E2E Test Template', 
        'Updated: Please {{action}} the {{target}}:\n\n{{selectedCode}}');
      
      // Wait for live reload
      await page.waitForTimeout(1000);
      
      // Re-apply template
      await page.click('[data-testid="apply-template"]');
      
      // Verify updated content
      const renderedOutput = page.locator('[data-testid="rendered-output"]');
      await expect(renderedOutput).toContainText('Updated: Please refactor');
    });
  });

  test('should handle template errors gracefully', async ({ page }) => {
    await test.step('Navigate to template creation', async () => {
      await page.goto('/templates');
      await page.click('[data-testid="create-template-button"]');
    });

    await test.step('Create template with invalid syntax', async () => {
      await page.fill('[data-testid="template-name"]', 'Invalid Template');
      await page.selectOption('[data-testid="template-category"]', 'code');
      await page.fill('[data-testid="template-content"]', 
        'Invalid template with {{unclosed variable');
      
      await page.click('[data-testid="save-template"]');
      
      // Should show error notification
      await expect(page.locator('[data-testid="error-notification"]'))
        .toContainText('Invalid template syntax');
    });

    await test.step('Fix template and save successfully', async () => {
      await page.fill('[data-testid="template-content"]', 
        'Fixed template with {{variable}}');
      
      await page.click('[data-testid="save-template"]');
      
      await expect(page.locator('[data-testid="success-notification"]'))
        .toContainText('Template created successfully');
    });
  });

  test('should support concurrent template usage', async ({ browser }) => {
    // Create multiple browser contexts to simulate concurrent users
    const contexts = await Promise.all([
      browser.newContext(),
      browser.newContext(),
      browser.newContext()
    ]);

    const pages = await Promise.all(contexts.map(context => context.newPage()));

    try {
      // Each user creates and uses templates concurrently
      await Promise.all(pages.map(async (page, index) => {
        await page.goto('/templates');
        
        // Create unique template
        await page.click('[data-testid="create-template-button"]');
        await page.fill('[data-testid="template-name"]', `Concurrent Template ${index}`);
        await page.selectOption('[data-testid="template-category"]', 'code');
        await page.fill('[data-testid="template-content"]', 
          `User ${index} template with {{variable}}`);
        await page.click('[data-testid="save-template"]');
        
        // Use template
        await page.goto('/editor');
        await page.fill('[data-testid="code-editor"]', `// User ${index} code`);
        await page.click('[data-testid="template-selector"]');
        await page.click(`[data-testid="template-option"][data-template="Concurrent Template ${index}"]`);
        await page.fill('[data-testid="template-variable-variable"]', `value-${index}`);
        await page.click('[data-testid="apply-template"]');
        
        // Verify result
        await expect(page.locator('[data-testid="rendered-output"]'))
          .toContainText(`User ${index} template with value-${index}`);
      }));
    } finally {
      // Cleanup
      await Promise.all(contexts.map(context => context.close()));
    }
  });
});
```

#### Performance E2E Testing
```typescript
// testing/e2e/performance.spec.ts
import { test, expect } from '@playwright/test';

test.describe('SuperClaude Performance Tests', () => {
  test('should meet performance targets for template operations', async ({ page }) => {
    // Navigate to application
    await page.goto('/');

    // Measure initial page load
    const loadStartTime = performance.now();
    await page.waitForLoadState('networkidle');
    const loadEndTime = performance.now();
    
    const pageLoadTime = loadEndTime - loadStartTime;
    expect(pageLoadTime).toBeLessThan(3000); // <3s page load

    // Measure template creation performance
    await page.goto('/templates');
    
    const createStartTime = performance.now();
    await page.click('[data-testid="create-template-button"]');
    await page.fill('[data-testid="template-name"]', 'Performance Test Template');
    await page.selectOption('[data-testid="template-category"]', 'code');
    await page.fill('[data-testid="template-content"]', 
      'Performance test template with {{variable}}');
    await page.click('[data-testid="save-template"]');
    
    await expect(page.locator('[data-testid="success-notification"]'))
      .toBeVisible();
    const createEndTime = performance.now();
    
    const createTime = createEndTime - createStartTime;
    expect(createTime).toBeLessThan(2000); // <2s template creation

    // Measure template rendering performance
    await page.goto('/editor');
    await page.fill('[data-testid="code-editor"]', 'test code');
    
    const renderStartTime = performance.now();
    await page.click('[data-testid="template-selector"]');
    await page.click('[data-testid="template-option"][data-template="Performance Test Template"]');
    await page.fill('[data-testid="template-variable-variable"]', 'test value');
    await page.click('[data-testid="apply-template"]');
    
    await expect(page.locator('[data-testid="rendered-output"]'))
      .toContainText('Performance test template with test value');
    const renderEndTime = performance.now();
    
    const renderTime = renderEndTime - renderStartTime;
    expect(renderTime).toBeLessThan(500); // <500ms template rendering
  });

  test('should handle high load template operations', async ({ page }) => {
    await page.goto('/editor');
    
    // Load large template content
    const largeContent = 'x'.repeat(10000); // 10KB template
    await page.fill('[data-testid="code-editor"]', largeContent);
    
    // Apply template multiple times rapidly
    const iterations = 10;
    const startTime = performance.now();
    
    for (let i = 0; i < iterations; i++) {
      await page.click('[data-testid="template-selector"]');
      await page.click('[data-testid="template-option"]:first-child');
      await page.click('[data-testid="apply-template"]');
      
      // Wait for result
      await expect(page.locator('[data-testid="rendered-output"]'))
        .not.toBeEmpty();
    }
    
    const endTime = performance.now();
    const avgTime = (endTime - startTime) / iterations;
    
    expect(avgTime).toBeLessThan(1000); // <1s average per operation
  });
});
```

### Cross-Browser Testing Strategy

#### Browser Compatibility Matrix
```typescript
// testing/e2e/cross-browser.spec.ts
import { test, expect, devices } from '@playwright/test';

const BROWSERS = [
  { name: 'Chrome', device: devices['Desktop Chrome'] },
  { name: 'Firefox', device: devices['Desktop Firefox'] },
  { name: 'Safari', device: devices['Desktop Safari'] },
  { name: 'Edge', device: devices['Desktop Edge'] },
];

const MOBILE_DEVICES = [
  { name: 'iPhone', device: devices['iPhone 12'] },
  { name: 'Android', device: devices['Pixel 5'] },
  { name: 'iPad', device: devices['iPad Pro'] },
];

BROWSERS.forEach(({ name, device }) => {
  test.describe(`${name} Browser Tests`, () => {
    test.use(device);

    test('should work correctly in ' + name, async ({ page }) => {
      await page.goto('/');
      
      // Test basic functionality
      await expect(page.locator('h1')).toBeVisible();
      
      // Test template creation
      await page.goto('/templates');
      await page.click('[data-testid="create-template-button"]');
      
      const modal = page.locator('[data-testid="template-modal"]');
      await expect(modal).toBeVisible();
      
      // Browser-specific assertions
      if (name === 'Safari') {
        // Safari-specific checks
        await expect(page).toHaveTitle(/OpenCode SuperClaude/);
      }
    });
  });
});

MOBILE_DEVICES.forEach(({ name, device }) => {
  test.describe(`${name} Mobile Tests`, () => {
    test.use(device);

    test('should be responsive on ' + name, async ({ page }) => {
      await page.goto('/');
      
      // Check mobile navigation
      const mobileMenu = page.locator('[data-testid="mobile-menu-toggle"]');
      if (await mobileMenu.isVisible()) {
        await mobileMenu.click();
        await expect(page.locator('[data-testid="mobile-menu"]')).toBeVisible();
      }
      
      // Test mobile template creation
      await page.goto('/templates');
      
      // Should adapt to mobile layout
      await expect(page.locator('[data-testid="template-list"]'))
        .toHaveCSS('flex-direction', 'column');
    });
  });
});
```

---

## ⚡ Performance Testing

### Load Testing Strategy

#### Artillery.js Configuration
```yaml
# testing/performance/load-test.yml
config:
  target: 'http://localhost:3000'
  phases:
    - duration: 60
      arrivalRate: 5
    - duration: 120
      arrivalRate: 10
    - duration: 60
      arrivalRate: 20
  defaults:
    headers:
      Content-Type: 'application/json'

scenarios:
  - name: 'Template Operations'
    weight: 70
    flow:
      - get:
          url: '/api/templates'
      - post:
          url: '/api/templates'
          json:
            name: 'Load Test Template {{ $randomString() }}'
            category: 'code'
            content: 'Load test template with {{variable}}'
      - post:
          url: '/api/templates/{{ $randomInt(1, 100) }}/render'
          json:
            context:
              variable: 'test value'

  - name: 'File Watcher Operations'
    weight: 20
    flow:
      - get:
          url: '/api/watcher/status'
      - post:
          url: '/api/watcher/reload'

  - name: 'Health Checks'
    weight: 10
    flow:
      - get:
          url: '/api/health'
      - get:
          url: '/api/metrics'
```

#### Go Load Testing
```go
// testing/performance/load_test.go
package performance

import (
    "context"
    "fmt"
    "sync"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPluginManagerLoad(t *testing.T) {
    // Setup
    manager := setupTestPluginManager(t)
    defer manager.Shutdown()

    // Load test parameters
    concurrency := 50
    requestsPerWorker := 100
    totalRequests := concurrency * requestsPerWorker

    // Metrics collection
    var (
        successCount int64
        errorCount   int64
        totalLatency time.Duration
        mu           sync.Mutex
    )

    // Start time
    startTime := time.Now()

    // Worker pool
    var wg sync.WaitGroup
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()

            for j := 0; j < requestsPerWorker; j++ {
                ctx := context.Background()
                
                request := &CommandRequest{
                    Command: "echo",
                    Args: map[string]interface{}{
                        "message": fmt.Sprintf("Worker %d, Request %d", workerID, j),
                    },
                }

                requestStart := time.Now()
                _, err := manager.ExecuteCommand(ctx, "test-plugin", request)
                requestLatency := time.Since(requestStart)

                mu.Lock()
                if err != nil {
                    errorCount++
                } else {
                    successCount++
                }
                totalLatency += requestLatency
                mu.Unlock()
            }
        }(i)
    }

    wg.Wait()
    totalTime := time.Since(startTime)

    // Assertions
    successRate := float64(successCount) / float64(totalRequests)
    assert.GreaterOrEqual(t, successRate, 0.95, "Success rate should be at least 95%")

    avgLatency := totalLatency / time.Duration(successCount)
    assert.Less(t, avgLatency, 100*time.Millisecond, "Average latency should be under 100ms")

    throughput := float64(successCount) / totalTime.Seconds()
    assert.GreaterOrEqual(t, throughput, 500.0, "Throughput should be at least 500 req/s")

    t.Logf("Load Test Results:")
    t.Logf("  Total Requests: %d", totalRequests)
    t.Logf("  Successful: %d", successCount)
    t.Logf("  Failed: %d", errorCount)
    t.Logf("  Success Rate: %.2f%%", successRate*100)
    t.Logf("  Average Latency: %v", avgLatency)
    t.Logf("  Throughput: %.2f req/s", throughput)
    t.Logf("  Total Time: %v", totalTime)
}

func BenchmarkEventSystemThroughput(b *testing.B) {
    eventBus := setupTestEventBus(b)
    defer eventBus.Stop()

    // Setup subscriber
    messageCount := 0
    eventBus.Subscribe("benchmark.*", func(event *Event) {
        messageCount++
    })

    b.ResetTimer()
    b.ReportAllocs()

    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            event := &Event{
                Type:   "benchmark.event",
                Source: "load-test",
                Data:   map[string]interface{}{"index": i},
            }
            
            eventBus.Publish(event)
            i++
        }
    })

    // Wait for all events to be processed
    time.Sleep(1 * time.Second)
    
    b.Logf("Messages processed: %d", messageCount)
}
```

### Memory and Resource Testing

#### Memory Leak Detection
```go
// testing/performance/memory_test.go
package performance

import (
    "runtime"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestMemoryLeaks(t *testing.T) {
    // Force garbage collection
    runtime.GC()
    runtime.GC()
    
    var initialStats runtime.MemStats
    runtime.ReadMemStats(&initialStats)

    // Setup system under test
    manager := setupTestPluginManager(t)
    defer manager.Shutdown()

    // Load plugins and execute operations
    for i := 0; i < 1000; i++ {
        // Simulate plugin operations that might cause leaks
        ctx := context.Background()
        request := &CommandRequest{
            Command: "echo",
            Args:    map[string]interface{}{"data": make([]byte, 1024)},
        }
        
        manager.ExecuteCommand(ctx, "test-plugin", request)
        
        // Periodic garbage collection
        if i%100 == 0 {
            runtime.GC()
        }
    }

    // Final memory check
    runtime.GC()
    runtime.GC()
    
    var finalStats runtime.MemStats
    runtime.ReadMemStats(&finalStats)

    // Memory growth should be reasonable
    memoryGrowth := finalStats.Alloc - initialStats.Alloc
    assert.Less(t, memoryGrowth, uint64(50*1024*1024), "Memory growth should be less than 50MB")

    t.Logf("Memory Stats:")
    t.Logf("  Initial Alloc: %d bytes", initialStats.Alloc)
    t.Logf("  Final Alloc: %d bytes", finalStats.Alloc)
    t.Logf("  Growth: %d bytes", memoryGrowth)
    t.Logf("  GC Cycles: %d", finalStats.NumGC-initialStats.NumGC)
}

func TestResourceCleanup(t *testing.T) {
    // Track resource allocation
    initialGoroutines := runtime.NumGoroutine()
    
    // Create and destroy multiple managers
    for i := 0; i < 10; i++ {
        manager := setupTestPluginManager(t)
        
        // Use the manager
        time.Sleep(100 * time.Millisecond)
        
        // Shutdown
        manager.Shutdown()
    }

    // Allow time for cleanup
    time.Sleep(1 * time.Second)
    runtime.GC()

    finalGoroutines := runtime.NumGoroutine()
    
    // Should not leak goroutines
    assert.LessOrEqual(t, finalGoroutines, initialGoroutines+2, 
        "Should not leak more than 2 goroutines")

    t.Logf("Goroutine Stats:")
    t.Logf("  Initial: %d", initialGoroutines)
    t.Logf("  Final: %d", finalGoroutines)
    t.Logf("  Difference: %d", finalGoroutines-initialGoroutines)
}
```

---

## 🛡️ Security Testing

### Security Test Framework

#### Security Validation Tests
```typescript
// testing/security/security.test.ts
import { describe, it, expect, beforeEach } from 'vitest';
import { SecurityValidator } from '../src/security/security-validator';
import { TemplateBroker } from '../src/broker/template-broker';

describe('Security Validation', () => {
  let validator: SecurityValidator;
  let broker: TemplateBroker;

  beforeEach(() => {
    validator = new SecurityValidator();
    broker = new TemplateBroker();
  });

  describe('Template Injection Prevention', () => {
    it('should prevent code injection in templates', async () => {
      const maliciousTemplate = {
        name: 'malicious-template',
        content: 'Hello {{name}}{{#each this}}{{@key}}: {{this}}{{/each}}<script>alert("XSS")</script>',
        frontmatter: {}
      };

      const result = await validator.validateTemplate(maliciousTemplate);
      
      expect(result.isSecure).toBe(false);
      expect(result.issues).toContainEqual(
        expect.objectContaining({
          type: 'code-injection',
          severity: 'high'
        })
      );
    });

    it('should sanitize template output', async () => {
      const template = {
        name: 'test-template',
        content: 'User input: {{userInput}}',
        frontmatter: {}
      };

      const maliciousContext = {
        userInput: '<script>alert("XSS")</script>'
      };

      const rendered = await broker.renderTemplate('test-template', maliciousContext);
      
      expect(rendered).not.toContain('<script>');
      expect(rendered).toContain('&lt;script&gt;');
    });
  });

  describe('Path Traversal Prevention', () => {
    it('should prevent path traversal in template paths', async () => {
      const maliciousPaths = [
        '../../../etc/passwd',
        '..\\..\\windows\\system32\\config\\sam',
        '/etc/shadow',
        'C:\\Windows\\system.ini'
      ];

      for (const path of maliciousPaths) {
        const result = await validator.validateTemplatePath(path);
        expect(result.isSecure).toBe(false);
        expect(result.issues).toContainEqual(
          expect.objectContaining({
            type: 'path-traversal',
            severity: 'critical'
          })
        );
      }
    });
  });

  describe('Resource Limits', () => {
    it('should enforce template size limits', async () => {
      const largeTemplate = {
        name: 'large-template',
        content: 'x'.repeat(10 * 1024 * 1024), // 10MB
        frontmatter: {}
      };

      await expect(validator.validateTemplate(largeTemplate))
        .rejects
        .toThrow('Template exceeds maximum size limit');
    });

    it('should prevent infinite loops in templates', async () => {
      const infiniteTemplate = {
        name: 'infinite-template',
        content: '{{#each items}}{{> recursivePartial}}{{/each}}',
        frontmatter: {}
      };

      const result = await validator.validateTemplate(infiniteTemplate);
      
      expect(result.isSecure).toBe(false);
      expect(result.issues).toContainEqual(
        expect.objectContaining({
          type: 'infinite-recursion',
          severity: 'high'
        })
      );
    });
  });

  describe('Input Validation', () => {
    it('should validate template frontmatter schema', async () => {
      const invalidTemplate = {
        name: 'invalid-template',
        content: 'Test content',
        frontmatter: {
          name: 123, // Should be string
          category: 'invalid-category', // Not in allowed list
          requires: 'not-an-array' // Should be array
        }
      };

      const result = await validator.validateTemplate(invalidTemplate);
      
      expect(result.isSecure).toBe(false);
      expect(result.issues.length).toBeGreaterThan(0);
    });
  });
});
```

#### Go Security Tests
```go
// testing/security/security_test.go
package security

import (
    "context"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPluginSecurity(t *testing.T) {
    manager := setupSecurePluginManager(t)
    defer manager.Shutdown()

    t.Run("should isolate plugin processes", func(t *testing.T) {
        // Load plugin
        ctx := context.Background()
        plugin, err := manager.LoadPlugin(ctx, "/test/plugins/isolated-plugin")
        require.NoError(t, err)

        // Plugin should run in separate process
        assert.NotEqual(t, os.Getpid(), plugin.GetPID())
        
        // Plugin should have limited permissions
        permissions := plugin.GetPermissions()
        assert.False(t, permissions.FileSystemAccess)
        assert.False(t, permissions.NetworkAccess)
        assert.True(t, permissions.ComputeOnly)
    })

    t.Run("should prevent unauthorized file access", func(t *testing.T) {
        request := &CommandRequest{
            Command: "read-file",
            Args: map[string]interface{}{
                "path": "/etc/passwd",
            },
        }

        ctx := context.Background()
        _, err := manager.ExecuteCommand(ctx, "test-plugin", request)
        
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "access denied")
    })

    t.Run("should enforce resource limits", func(t *testing.T) {
        request := &CommandRequest{
            Command: "allocate-memory",
            Args: map[string]interface{}{
                "size": "1GB", // Exceeds limit
            },
        }

        ctx := context.Background()
        _, err := manager.ExecuteCommand(ctx, "test-plugin", request)
        
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "resource limit exceeded")
    })

    t.Run("should validate command injection", func(t *testing.T) {
        maliciousCommands := []string{
            "test; rm -rf /",
            "test && cat /etc/passwd",
            "test | sh",
            "test $(cat /etc/passwd)",
        }

        for _, cmd := range maliciousCommands {
            request := &CommandRequest{
                Command: "execute",
                Args: map[string]interface{}{
                    "command": cmd,
                },
            }

            ctx := context.Background()
            _, err := manager.ExecuteCommand(ctx, "test-plugin", request)
            
            assert.Error(t, err, "Should reject malicious command: %s", cmd)
            assert.Contains(t, err.Error(), "invalid command")
        }
    })
}

func TestEventSecurity(t *testing.T) {
    eventBus := setupSecureEventBus(t)
    defer eventBus.Stop()

    t.Run("should validate event payloads", func(t *testing.T) {
        maliciousEvents := []*Event{
            {
                Type: "test",
                Data: map[string]interface{}{
                    "script": "<script>alert('xss')</script>",
                },
            },
            {
                Type: "test",
                Data: map[string]interface{}{
                    "path": "../../../etc/passwd",
                },
            },
        }

        for _, event := range maliciousEvents {
            err := eventBus.Publish(event)
            assert.Error(t, err, "Should reject malicious event")
        }
    })

    t.Run("should enforce event size limits", func(t *testing.T) {
        largeEvent := &Event{
            Type: "test",
            Data: map[string]interface{}{
                "payload": strings.Repeat("x", 10*1024*1024), // 10MB
            },
        }

        err := eventBus.Publish(largeEvent)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "event size exceeds limit")
    })
}
```

### Penetration Testing

#### Automated Security Scanning
```yaml
# testing/security/security-scan.yml
version: '3'
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=test
      - SECURITY_TESTING=true

  zap:
    image: owasp/zap2docker-stable
    command: >
      zap-baseline.py 
      -t http://app:3000 
      -J /zap/wrk/baseline-report.json
      -r /zap/wrk/baseline-report.html
    volumes:
      - ./security-reports:/zap/wrk
    depends_on:
      - app

  nikto:
    image: sullo/nikto
    command: >
      -h http://app:3000 
      -Format htm 
      -output /reports/nikto-report.html
    volumes:
      - ./security-reports:/reports
    depends_on:
      - app
```

---

## 📊 CI/CD Pipeline Integration

### GitHub Actions Workflow

#### Main Testing Pipeline
```yaml
# .github/workflows/test.yml
name: Test Suite

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  unit-tests-typescript:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        node-version: [18, 20]
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Node.js ${{ matrix.node-version }}
      uses: actions/setup-node@v4
      with:
        node-version: ${{ matrix.node-version }}
        cache: 'npm'
    
    - name: Install dependencies
      run: npm ci
    
    - name: Run TypeScript unit tests
      run: npm run test:unit
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage/lcov.info
        flags: typescript-unit

  unit-tests-go:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.21, 1.22]
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Go ${{ matrix.go-version }}
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run Go unit tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.out
        flags: go-unit

  integration-tests:
    runs-on: ubuntu-latest
    needs: [unit-tests-typescript, unit-tests-go]
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
          POSTGRES_DB: superclaude_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      
      nats:
        image: nats:2.10-alpine
        ports:
          - 4222:4222
        options: >-
          --health-cmd "nats-server --version"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: 20
        cache: 'npm'
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: 1.22
    
    - name: Install dependencies
      run: |
        npm ci
        go mod download
    
    - name: Build test plugins
      run: make build-test-plugins
    
    - name: Run integration tests
      run: |
        npm run test:integration
        go test -v -tags=integration ./testing/integration/...
      env:
        DATABASE_URL: postgres://postgres:test@localhost:5432/superclaude_test
        NATS_URL: nats://localhost:4222

  e2e-tests:
    runs-on: ubuntu-latest
    needs: [integration-tests]
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: 20
        cache: 'npm'
    
    - name: Install dependencies
      run: npm ci
    
    - name: Install Playwright
      run: npx playwright install --with-deps
    
    - name: Build application
      run: npm run build:test
    
    - name: Run E2E tests
      run: npx playwright test
      env:
        BASE_URL: http://localhost:3000
    
    - name: Upload E2E artifacts
      uses: actions/upload-artifact@v3
      if: failure()
      with:
        name: e2e-artifacts
        path: |
          testing/e2e-results/
          testing/e2e-reports/

  performance-tests:
    runs-on: ubuntu-latest
    needs: [integration-tests]
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: 20
        cache: 'npm'
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: 1.22
    
    - name: Install dependencies
      run: |
        npm ci
        go mod download
        npm install -g artillery
    
    - name: Build application
      run: npm run build:test
    
    - name: Run performance tests
      run: |
        npm run start:test &
        sleep 10
        artillery run testing/performance/load-test.yml
        go test -bench=. -benchmem ./testing/performance/...

  security-tests:
    runs-on: ubuntu-latest
    needs: [integration-tests]
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: 20
        cache: 'npm'
    
    - name: Install dependencies
      run: npm ci
    
    - name: Run security audit
      run: |
        npm audit --audit-level high
        go list -json -m all | nancy sleuth
    
    - name: Run SAST analysis
      uses: securecodewarrior/github-action-add-sarif@v1
      with:
        sarif-file: 'security-analysis.sarif'
    
    - name: Build and run security tests
      run: |
        npm run build:test
        npm run test:security
        go test -v ./testing/security/...

  cross-platform-tests:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: 20
        cache: 'npm'
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: 1.22
    
    - name: Install dependencies
      run: |
        npm ci
        go mod download
    
    - name: Run cross-platform tests
      run: |
        npm run test:unit
        go test -v ./...
```

### Test Coverage Requirements

#### Coverage Configuration
```typescript
// vitest.config.ts
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html', 'lcov'],
      reportsDirectory: './coverage',
      
      // Coverage thresholds
      thresholds: {
        global: {
          branches: 80,
          functions: 80,
          lines: 80,
          statements: 80
        },
        
        // Specific file patterns
        './src/broker/': {
          branches: 90,
          functions: 90,
          lines: 90,
          statements: 90
        },
        
        './src/security/': {
          branches: 95,
          functions: 95,
          lines: 95,
          statements: 95
        }
      },
      
      // Exclude patterns
      exclude: [
        'node_modules/',
        'testing/',
        '**/*.test.{ts,js}',
        '**/*.spec.{ts,js}',
        '**/mocks/',
        'dist/',
        'build/'
      ]
    }
  }
});
```

#### Go Coverage Configuration
```makefile
# Makefile for Go coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running Go tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'
	@echo "Coverage report generated: coverage.html"

.PHONY: test-coverage-check
test-coverage-check:
	@echo "Checking coverage thresholds..."
	@coverage=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ "$${coverage%.*}" -lt 80 ]; then \
		echo "Coverage $${coverage}% is below threshold of 80%"; \
		exit 1; \
	else \
		echo "Coverage $${coverage}% meets threshold"; \
	fi
```

---

## 📋 Testing Checklist

### Pre-Deployment Testing Checklist

#### Unit Testing Verification
- [ ] **TypeScript Unit Tests**
  - [ ] All components have ≥80% code coverage
  - [ ] All business logic functions tested
  - [ ] Edge cases and error conditions covered
  - [ ] Performance requirements validated
  - [ ] Mock dependencies properly isolated

- [ ] **Go Unit Tests**
  - [ ] All packages have ≥80% code coverage
  - [ ] Plugin manager thoroughly tested
  - [ ] Event system functionality verified
  - [ ] Concurrent operations tested
  - [ ] Resource management validated

#### Integration Testing Verification
- [ ] **API Integration**
  - [ ] All API endpoints tested
  - [ ] Database operations verified
  - [ ] File watcher integration working
  - [ ] Template rendering pipeline tested
  - [ ] Error handling validated

- [ ] **Plugin System Integration**
  - [ ] Plugin discovery and loading
  - [ ] Command execution and responses
  - [ ] Health monitoring and recovery
  - [ ] Event handling and propagation
  - [ ] Security isolation verified

#### End-to-End Testing Verification
- [ ] **User Workflows**
  - [ ] Template creation and management
  - [ ] Template usage in editor
  - [ ] Live reload functionality
  - [ ] Error handling and recovery
  - [ ] Performance within targets

- [ ] **Cross-Browser Compatibility**
  - [ ] Chrome/Chromium functionality
  - [ ] Firefox compatibility
  - [ ] Safari/WebKit support
  - [ ] Mobile device responsiveness
  - [ ] Accessibility compliance

#### Performance Testing Verification
- [ ] **Load Testing**
  - [ ] System handles target load
  - [ ] Response times within limits
  - [ ] Resource usage acceptable
  - [ ] No memory leaks detected
  - [ ] Graceful degradation under stress

- [ ] **Security Testing**
  - [ ] Input validation working
  - [ ] XSS prevention implemented
  - [ ] Path traversal blocked
  - [ ] Code injection prevented
  - [ ] Resource limits enforced

### Continuous Monitoring

#### Production Testing Strategy
```typescript
// monitoring/production-tests.ts
export class ProductionTestSuite {
  private healthChecker: HealthChecker;
  private performanceMonitor: PerformanceMonitor;
  private securityScanner: SecurityScanner;

  async runContinuousTests(): Promise<void> {
    // Health checks every 5 minutes
    setInterval(async () => {
      const health = await this.healthChecker.checkSystemHealth();
      if (health.status !== 'healthy') {
        await this.alertManager.sendAlert('System health degraded', health);
      }
    }, 5 * 60 * 1000);

    // Performance monitoring every hour
    setInterval(async () => {
      const metrics = await this.performanceMonitor.collectMetrics();
      if (metrics.responseTime > 200) {
        await this.alertManager.sendAlert('Performance degradation', metrics);
      }
    }, 60 * 60 * 1000);

    // Security scans daily
    setInterval(async () => {
      const scanResults = await this.securityScanner.runScan();
      if (scanResults.hasVulnerabilities) {
        await this.alertManager.sendAlert('Security vulnerabilities detected', scanResults);
      }
    }, 24 * 60 * 60 * 1000);
  }
}
```

---

## 🎯 Testing Success Metrics

### Quality Gates

#### Mandatory Quality Requirements
- **Unit Test Coverage**: ≥80% for all components
- **Integration Test Coverage**: ≥70% for critical paths
- **E2E Test Coverage**: 100% for user journeys
- **Performance Targets**: <200ms API response, <3s UI load
- **Error Rate**: <0.1% in production
- **Security Compliance**: Zero high-severity vulnerabilities

#### Performance Benchmarks
- **Template Loading**: <100ms for templates <1MB
- **Template Rendering**: <50ms for templates with <10 variables
- **Plugin Operations**: <200ms for standard commands
- **Event Processing**: >1000 events/second throughput
- **Memory Usage**: <256MB for plugin system
- **CPU Usage**: <50% under normal load

#### Release Criteria
- [ ] All unit tests passing
- [ ] All integration tests passing
- [ ] All E2E tests passing
- [ ] Performance benchmarks met
- [ ] Security scan clean
- [ ] Load testing successful
- [ ] Cross-browser compatibility verified
- [ ] Accessibility compliance validated
- [ ] Documentation updated
- [ ] Migration path tested

---

*This comprehensive testing strategy ensures the reliability, performance, and security of OpenCode SuperClaude integrations across all architectural approaches. Regular testing and continuous monitoring are essential for maintaining system quality and user satisfaction.*