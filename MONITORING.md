# Monitoring and Observability Guide

## 📊 Overview

This guide provides comprehensive monitoring, logging, and observability setup for the OpenCode-SuperClaude Integration project. Proper observability is crucial for maintaining performance, diagnosing issues, and ensuring reliable operation across all integration plans.

## 📋 Table of Contents

- [Observability Philosophy](#observability-philosophy)
- [Logging Standards](#logging-standards)
- [Metrics Collection](#metrics-collection)
- [Distributed Tracing](#distributed-tracing)
- [Performance Monitoring](#performance-monitoring)
- [Alert Configuration](#alert-configuration)
- [Dashboard Setup](#dashboard-setup)
- [Debugging Techniques](#debugging-techniques)
- [Plan-Specific Monitoring](#plan-specific-monitoring)

## 🎯 Observability Philosophy

### Three Pillars of Observability

1. **Metrics**: Quantitative data about system behavior
2. **Logs**: Discrete events with contextual information
3. **Traces**: Request flows through distributed systems

### Core Principles

- **Observability by Design**: Build monitoring into the system from the start
- **Context Preservation**: Maintain request context across all components
- **Actionable Information**: Every alert should be actionable
- **Performance Impact**: Monitoring should not significantly impact system performance
- **Security First**: Never log sensitive information

## 📝 Logging Standards

### Log Levels

```go
// Go logging levels
const (
    LevelDebug = "debug"  // Detailed information for debugging
    LevelInfo  = "info"   // General information about program execution
    LevelWarn  = "warn"   // Warning messages for potential issues
    LevelError = "error"  // Error messages for failures
    LevelFatal = "fatal"  // Critical errors that cause program termination
)
```

```typescript
// TypeScript logging levels
enum LogLevel {
  DEBUG = 'debug',
  INFO = 'info',
  WARN = 'warn',
  ERROR = 'error',
  FATAL = 'fatal'
}
```

### Structured Logging Format

#### Go Implementation

```go
package logging

import (
    "context"
    "encoding/json"
    "log/slog"
    "os"
    "time"
)

type Logger struct {
    *slog.Logger
}

func NewLogger(level string) *Logger {
    var logLevel slog.Level
    switch level {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    default:
        logLevel = slog.LevelInfo
    }

    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: logLevel,
        AddSource: true,
    })

    return &Logger{
        Logger: slog.New(handler),
    }
}

func (l *Logger) LogWithContext(ctx context.Context, level slog.Level, msg string, args ...any) {
    // Extract trace information from context
    if traceID := TraceIDFromContext(ctx); traceID != "" {
        args = append(args, "trace_id", traceID)
    }
    if spanID := SpanIDFromContext(ctx); spanID != "" {
        args = append(args, "span_id", spanID)
    }
    
    l.Log(ctx, level, msg, args...)
}

// Usage examples
func (s *PluginManager) LoadPlugin(ctx context.Context, name string) error {
    logger.LogWithContext(ctx, slog.LevelInfo, "Loading plugin",
        "plugin_name", name,
        "plugin_dir", s.pluginDir,
    )
    
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        logger.LogWithContext(ctx, slog.LevelInfo, "Plugin load completed",
            "plugin_name", name,
            "duration_ms", duration.Milliseconds(),
        )
    }()
    
    // Plugin loading logic...
    if err != nil {
        logger.LogWithContext(ctx, slog.LevelError, "Failed to load plugin",
            "plugin_name", name,
            "error", err.Error(),
        )
        return err
    }
    
    return nil
}
```

#### TypeScript Implementation

```typescript
import winston from 'winston';
import { v4 as uuidv4 } from 'uuid';

interface LogContext {
  traceId?: string;
  spanId?: string;
  userId?: string;
  pluginName?: string;
  operation?: string;
}

class Logger {
  private winston: winston.Logger;

  constructor(level: string = 'info') {
    this.winston = winston.createLogger({
      level,
      format: winston.format.combine(
        winston.format.timestamp(),
        winston.format.errors({ stack: true }),
        winston.format.json(),
        winston.format.printf(({ timestamp, level, message, ...meta }) => {
          return JSON.stringify({
            timestamp,
            level,
            message,
            ...meta
          });
        })
      ),
      transports: [
        new winston.transports.Console(),
        new winston.transports.File({ filename: 'opencode.log' })
      ]
    });
  }

  private formatLog(level: string, message: string, context?: LogContext, extra?: any) {
    return {
      level,
      message,
      ...context,
      ...extra,
      timestamp: new Date().toISOString()
    };
  }

  debug(message: string, context?: LogContext, extra?: any) {
    this.winston.debug(this.formatLog('debug', message, context, extra));
  }

  info(message: string, context?: LogContext, extra?: any) {
    this.winston.info(this.formatLog('info', message, context, extra));
  }

  warn(message: string, context?: LogContext, extra?: any) {
    this.winston.warn(this.formatLog('warn', message, context, extra));
  }

  error(message: string, context?: LogContext, extra?: any) {
    this.winston.error(this.formatLog('error', message, context, extra));
  }
}

// Usage example
const logger = new Logger();

async function executePlugin(name: string, args: string[]): Promise<PluginResult> {
  const traceId = uuidv4();
  const context: LogContext = { traceId, pluginName: name, operation: 'execute_plugin' };
  
  logger.info('Starting plugin execution', context, { args });
  
  const startTime = performance.now();
  try {
    const result = await pluginRegistry.execute(name, args);
    const duration = performance.now() - startTime;
    
    logger.info('Plugin execution completed', context, { 
      duration_ms: duration,
      result_status: result.status 
    });
    
    return result;
  } catch (error) {
    const duration = performance.now() - startTime;
    logger.error('Plugin execution failed', context, { 
      duration_ms: duration,
      error: error.message,
      stack: error.stack 
    });
    throw error;
  }
}
```

### Log Output Format

```json
{
  "timestamp": "2024-01-15T10:30:45.123Z",
  "level": "info",
  "message": "Plugin execution completed",
  "trace_id": "abc123def456",
  "span_id": "span789",
  "plugin_name": "superclaude",
  "operation": "execute_plugin",
  "duration_ms": 150,
  "result_status": "success"
}
```

### Security Considerations

**Never log sensitive information**:

```go
// ❌ BAD - Don't log sensitive data
logger.Info("User authenticated", "password", password, "api_key", apiKey)

// ✅ GOOD - Log safely
logger.Info("User authenticated", "user_id", userID, "method", "oauth")
```

**Sanitize log data**:

```typescript
function sanitizeLogData(data: any): any {
  const sensitive = ['password', 'token', 'api_key', 'secret', 'private_key'];
  const sanitized = { ...data };
  
  for (const key of Object.keys(sanitized)) {
    if (sensitive.some(s => key.toLowerCase().includes(s))) {
      sanitized[key] = '[REDACTED]';
    }
  }
  
  return sanitized;
}
```

## 📊 Metrics Collection

### Prometheus Integration

#### Go Metrics Setup

```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
    "time"
)

var (
    // Plugin metrics
    pluginExecutions = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "opencode_plugin_executions_total",
            Help: "Total number of plugin executions",
        },
        []string{"plugin_name", "status"},
    )
    
    pluginDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "opencode_plugin_duration_seconds",
            Help: "Plugin execution duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"plugin_name"},
    )
    
    activePlugins = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "opencode_active_plugins",
            Help: "Number of currently active plugins",
        },
    )
    
    // MCP metrics
    mcpConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "opencode_mcp_connections",
            Help: "Number of active MCP connections",
        },
    )
    
    mcpRequests = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "opencode_mcp_requests_total",
            Help: "Total number of MCP requests",
        },
        []string{"method", "status"},
    )
    
    mcpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "opencode_mcp_request_duration_seconds",
            Help: "MCP request duration in seconds",
        },
        []string{"method"},
    )
)

// Metrics middleware
func MetricsMiddleware(next http.Handler) http.Handler {
    return promhttp.InstrumentHandlerInFlight(
        promauto.NewGauge(prometheus.GaugeOpts{
            Name: "http_requests_in_flight",
            Help: "Current number of HTTP requests being served",
        }),
        promhttp.InstrumentHandlerDuration(
            promauto.NewHistogramVec(
                prometheus.HistogramOpts{
                    Name: "http_request_duration_seconds",
                    Help: "HTTP request duration in seconds",
                },
                []string{"method", "path"},
            ),
            next,
        ),
    )
}

// Plugin execution tracking
func TrackPluginExecution(pluginName string, fn func() error) error {
    timer := prometheus.NewTimer(pluginDuration.WithLabelValues(pluginName))
    defer timer.ObserveDuration()
    
    activePlugins.Inc()
    defer activePlugins.Dec()
    
    err := fn()
    
    status := "success"
    if err != nil {
        status = "error"
    }
    pluginExecutions.WithLabelValues(pluginName, status).Inc()
    
    return err
}

// Start metrics server
func StartMetricsServer(port string) error {
    http.Handle("/metrics", promhttp.Handler())
    return http.ListenAndServe(":"+port, nil)
}
```

#### TypeScript Metrics Setup

```typescript
import client from 'prom-client';
import express from 'express';

// Create a Registry to register the metrics
const register = new client.Registry();

// Add default metrics
client.collectDefaultMetrics({ 
  register,
  prefix: 'opencode_' 
});

// Custom metrics
const pluginExecutions = new client.Counter({
  name: 'opencode_plugin_executions_total',
  help: 'Total number of plugin executions',
  labelNames: ['plugin_name', 'status'],
  registers: [register]
});

const pluginDuration = new client.Histogram({
  name: 'opencode_plugin_duration_seconds',
  help: 'Plugin execution duration in seconds',
  labelNames: ['plugin_name'],
  registers: [register]
});

const activeConnections = new client.Gauge({
  name: 'opencode_active_connections',
  help: 'Number of active connections',
  registers: [register]
});

// Metrics tracking utilities
export class MetricsTracker {
  static trackPluginExecution<T>(
    pluginName: string, 
    fn: () => Promise<T>
  ): Promise<T> {
    const endTimer = pluginDuration.startTimer({ plugin_name: pluginName });
    
    return fn()
      .then(result => {
        pluginExecutions.inc({ plugin_name: pluginName, status: 'success' });
        return result;
      })
      .catch(error => {
        pluginExecutions.inc({ plugin_name: pluginName, status: 'error' });
        throw error;
      })
      .finally(() => {
        endTimer();
      });
  }
  
  static incrementConnections() {
    activeConnections.inc();
  }
  
  static decrementConnections() {
    activeConnections.dec();
  }
}

// Metrics endpoint
export function setupMetricsEndpoint(app: express.Application) {
  app.get('/metrics', async (req, res) => {
    try {
      res.set('Content-Type', register.contentType);
      res.end(await register.metrics());
    } catch (error) {
      res.status(500).end(error);
    }
  });
}
```

### Key Metrics to Track

#### System Metrics
- **CPU Usage**: `opencode_cpu_usage_percent`
- **Memory Usage**: `opencode_memory_usage_bytes`
- **Disk I/O**: `opencode_disk_io_bytes_total`
- **Network I/O**: `opencode_network_io_bytes_total`

#### Application Metrics
- **Plugin Executions**: `opencode_plugin_executions_total`
- **Plugin Duration**: `opencode_plugin_duration_seconds`
- **Active Plugins**: `opencode_active_plugins`
- **MCP Requests**: `opencode_mcp_requests_total`
- **Error Rates**: `opencode_errors_total`

#### Business Metrics
- **User Sessions**: `opencode_user_sessions_total`
- **Commands Executed**: `opencode_commands_total`
- **Success Rate**: `opencode_success_rate`

## 🔍 Distributed Tracing

### OpenTelemetry Setup

#### Go Tracing Implementation

```go
package tracing

import (
    "context"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const serviceName = "opencode-mcp-server"

func InitTracer() (*trace.TracerProvider, error) {
    // Create Jaeger exporter
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
    if err != nil {
        return nil, err
    }

    // Create resource
    resource, err := resource.Merge(
        resource.Default(),
        resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(serviceName),
            semconv.ServiceVersionKey.String("1.0.0"),
        ),
    )
    if err != nil {
        return nil, err
    }

    // Create tracer provider
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource),
        trace.WithSampler(trace.AlwaysSample()),
    )

    otel.SetTracerProvider(tp)
    return tp, nil
}

// Plugin execution with tracing
func (m *PluginManager) ExecutePluginWithTracing(ctx context.Context, name string, args []string) (*Result, error) {
    tracer := otel.Tracer(serviceName)
    ctx, span := tracer.Start(ctx, "plugin.execute",
        trace.WithAttributes(
            attribute.String("plugin.name", name),
            attribute.StringSlice("plugin.args", args),
        ),
    )
    defer span.End()

    // Add plugin metadata to span
    plugin, err := m.GetPlugin(name)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }

    span.SetAttributes(
        attribute.String("plugin.version", plugin.Version),
        attribute.String("plugin.type", plugin.Type),
    )

    // Execute plugin
    result, err := m.executePlugin(ctx, plugin, args)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }

    span.SetAttributes(
        attribute.String("result.status", result.Status),
        attribute.Int("result.size", len(result.Output)),
    )

    span.SetStatus(codes.Ok, "Plugin executed successfully")
    return result, nil
}
```

#### TypeScript Tracing Implementation

```typescript
import { NodeSDK } from '@opentelemetry/sdk-node';
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';
import { JaegerExporter } from '@opentelemetry/exporter-jaeger';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { trace, context, SpanStatusCode } from '@opentelemetry/api';

// Initialize OpenTelemetry SDK
const jaegerExporter = new JaegerExporter({
  endpoint: 'http://localhost:14268/api/traces',
});

const sdk = new NodeSDK({
  resource: new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: 'opencode-integration',
    [SemanticResourceAttributes.SERVICE_VERSION]: '1.0.0',
  }),
  traceExporter: jaegerExporter,
  instrumentations: [getNodeAutoInstrumentations()]
});

sdk.start();

// Tracing utilities
export class TracingUtils {
  private static tracer = trace.getTracer('opencode-integration');

  static async executeWithTracing<T>(
    operationName: string,
    attributes: Record<string, string | number | boolean>,
    fn: () => Promise<T>
  ): Promise<T> {
    const span = this.tracer.startSpan(operationName, { attributes });
    
    try {
      const result = await fn();
      span.setStatus({ code: SpanStatusCode.OK });
      return result;
    } catch (error) {
      span.recordException(error);
      span.setStatus({
        code: SpanStatusCode.ERROR,
        message: error.message
      });
      throw error;
    } finally {
      span.end();
    }
  }

  static createChildSpan(operationName: string, attributes?: Record<string, any>) {
    return this.tracer.startSpan(operationName, { attributes });
  }
}

// Usage in plugin execution
async function executePlugin(name: string, args: string[]): Promise<PluginResult> {
  return TracingUtils.executeWithTracing(
    'plugin.execute',
    {
      'plugin.name': name,
      'plugin.args.count': args.length,
    },
    async () => {
      const childSpan = TracingUtils.createChildSpan('plugin.load', {
        'plugin.name': name
      });
      
      try {
        const plugin = await pluginRegistry.getPlugin(name);
        childSpan.setAttributes({
          'plugin.version': plugin.version,
          'plugin.type': plugin.type
        });
        
        const result = await plugin.execute(args);
        
        childSpan.setAttributes({
          'result.status': result.status,
          'result.output.length': result.output.length
        });
        
        return result;
      } finally {
        childSpan.end();
      }
    }
  );
}
```

### Trace Context Propagation

```go
// Context keys for trace information
type contextKey string

const (
    traceIDKey contextKey = "trace_id"
    spanIDKey  contextKey = "span_id"
)

// Extract trace information from context
func TraceIDFromContext(ctx context.Context) string {
    if traceID, ok := ctx.Value(traceIDKey).(string); ok {
        return traceID
    }
    return ""
}

func SpanIDFromContext(ctx context.Context) string {
    if spanID, ok := ctx.Value(spanIDKey).(string); ok {
        return spanID
    }
    return ""
}

// Add trace information to context
func ContextWithTraceInfo(ctx context.Context, traceID, spanID string) context.Context {
    ctx = context.WithValue(ctx, traceIDKey, traceID)
    ctx = context.WithValue(ctx, spanIDKey, spanID)
    return ctx
}
```

## ⚡ Performance Monitoring

### Response Time Monitoring

```go
// Performance monitoring middleware
func PerformanceMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap response writer to capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(wrapped, r)
        
        duration := time.Since(start)
        
        // Log slow requests
        if duration > 100*time.Millisecond {
            logger.Warn("Slow request detected",
                "method", r.Method,
                "path", r.URL.Path,
                "duration_ms", duration.Milliseconds(),
                "status_code", wrapped.statusCode,
            )
        }
        
        // Record metrics
        httpRequestDuration.WithLabelValues(
            r.Method,
            r.URL.Path,
            fmt.Sprintf("%d", wrapped.statusCode),
        ).Observe(duration.Seconds())
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

### Memory Monitoring

```go
package monitoring

import (
    "runtime"
    "time"
)

type MemoryStats struct {
    Alloc        uint64    `json:"alloc"`
    TotalAlloc   uint64    `json:"total_alloc"`
    Sys          uint64    `json:"sys"`
    NumGC        uint32    `json:"num_gc"`
    Timestamp    time.Time `json:"timestamp"`
}

func GetMemoryStats() *MemoryStats {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return &MemoryStats{
        Alloc:      m.Alloc,
        TotalAlloc: m.TotalAlloc,
        Sys:        m.Sys,
        NumGC:      m.NumGC,
        Timestamp:  time.Now(),
    }
}

// Memory monitoring goroutine
func StartMemoryMonitoring(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            stats := GetMemoryStats()
            
            // Log memory usage
            logger.Debug("Memory stats",
                "alloc_mb", stats.Alloc/1024/1024,
                "total_alloc_mb", stats.TotalAlloc/1024/1024,
                "sys_mb", stats.Sys/1024/1024,
                "num_gc", stats.NumGC,
            )
            
            // Update Prometheus metrics
            memoryUsage.Set(float64(stats.Alloc))
            totalAllocated.Set(float64(stats.TotalAlloc))
            
            // Alert on high memory usage
            if stats.Alloc > 500*1024*1024 { // 500MB
                logger.Warn("High memory usage detected",
                    "alloc_mb", stats.Alloc/1024/1024,
                )
            }
        }
    }()
}
```

### CPU Monitoring

```typescript
import os from 'os';
import { EventEmitter } from 'events';

export class CPUMonitor extends EventEmitter {
  private interval: NodeJS.Timeout | null = null;
  private previousCPUUsage: NodeJS.CpuUsage = { user: 0, system: 0 };

  start(intervalMs: number = 5000) {
    this.interval = setInterval(() => {
      this.measureCPUUsage();
    }, intervalMs);
  }

  stop() {
    if (this.interval) {
      clearInterval(this.interval);
      this.interval = null;
    }
  }

  private measureCPUUsage() {
    const currentCPUUsage = process.cpuUsage();
    const cpuUsagePercent = this.calculateCPUPercent(currentCPUUsage);
    
    // Update metrics
    cpuUsageGauge.set(cpuUsagePercent);
    
    // Log high CPU usage
    if (cpuUsagePercent > 80) {
      logger.warn('High CPU usage detected', {
        cpu_percent: cpuUsagePercent,
        user_time: currentCPUUsage.user,
        system_time: currentCPUUsage.system
      });
    }
    
    // Emit event for listeners
    this.emit('cpu-usage', cpuUsagePercent);
    
    this.previousCPUUsage = currentCPUUsage;
  }

  private calculateCPUPercent(currentUsage: NodeJS.CpuUsage): number {
    const userUsage = currentUsage.user - this.previousCPUUsage.user;
    const systemUsage = currentUsage.system - this.previousCPUUsage.system;
    const totalUsage = userUsage + systemUsage;
    
    // Convert to percentage (CPU usage is in microseconds)
    return (totalUsage / 1000 / 5000) * 100; // 5000ms interval
  }
}

// Usage
const cpuMonitor = new CPUMonitor();
cpuMonitor.start(5000);

cpuMonitor.on('cpu-usage', (percent) => {
  if (percent > 90) {
    // Trigger alert
    alertManager.sendAlert({
      severity: 'critical',
      message: `High CPU usage: ${percent.toFixed(2)}%`,
      timestamp: new Date()
    });
  }
});
```

## 🚨 Alert Configuration

### Alert Rules

```yaml
# prometheus-alerts.yml
groups:
  - name: opencode-alerts
    rules:
      # High error rate
      - alert: HighErrorRate
        expr: rate(opencode_plugin_executions_total{status="error"}[5m]) > 0.1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "High plugin error rate detected"
          description: "Error rate is {{ $value }} errors per second"

      # High response time
      - alert: HighResponseTime
        expr: histogram_quantile(0.95, rate(opencode_plugin_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High plugin response time"
          description: "95th percentile response time is {{ $value }} seconds"

      # High memory usage
      - alert: HighMemoryUsage
        expr: opencode_memory_usage_bytes > 1073741824  # 1GB
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High memory usage"
          description: "Memory usage is {{ $value | humanize }}B"

      # Plugin down
      - alert: PluginDown
        expr: opencode_active_plugins == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "No active plugins"
          description: "All plugins are down"

      # MCP connection issues
      - alert: MCPConnectionLoss
        expr: opencode_mcp_connections == 0
        for: 30s
        labels:
          severity: critical
        annotations:
          summary: "MCP connection lost"
          description: "No active MCP connections"
```

### Alert Manager Configuration

```yaml
# alertmanager.yml
global:
  smtp_smarthost: 'localhost:587'
  smtp_from: 'alerts@opencode.dev'

route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'default'

receivers:
  - name: 'default'
    email_configs:
      - to: 'dev-team@opencode.dev'
        subject: 'OpenCode Alert: {{ .GroupLabels.alertname }}'
        body: |
          {{ range .Alerts }}
          Alert: {{ .Annotations.summary }}
          Description: {{ .Annotations.description }}
          Labels: {{ range .Labels }}{{ .Name }}={{ .Value }} {{ end }}
          {{ end }}

    slack_configs:
      - api_url: 'YOUR_SLACK_WEBHOOK_URL'
        channel: '#alerts'
        title: 'OpenCode Alert'
        text: |
          {{ range .Alerts }}
          *{{ .Annotations.summary }}*
          {{ .Annotations.description }}
          {{ end }}

    webhook_configs:
      - url: 'http://localhost:9093/webhook'
        send_resolved: true
```

### Custom Alert Handler

```go
package alerts

import (
    "bytes"
    "encoding/json"
    "net/http"
    "time"
)

type Alert struct {
    Status      string            `json:"status"`
    Labels      map[string]string `json:"labels"`
    Annotations map[string]string `json:"annotations"`
    StartsAt    time.Time         `json:"startsAt"`
    EndsAt      time.Time         `json:"endsAt"`
}

type AlertPayload struct {
    Version           string    `json:"version"`
    GroupKey          string    `json:"groupKey"`
    TruncatedAlerts   int       `json:"truncatedAlerts"`
    Status            string    `json:"status"`
    Receiver          string    `json:"receiver"`
    GroupLabels       map[string]string `json:"groupLabels"`
    CommonLabels      map[string]string `json:"commonLabels"`
    CommonAnnotations map[string]string `json:"commonAnnotations"`
    ExternalURL       string    `json:"externalURL"`
    Alerts            []Alert   `json:"alerts"`
}

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
    var payload AlertPayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for _, alert := range payload.Alerts {
        processAlert(alert)
    }

    w.WriteHeader(http.StatusOK)
}

func processAlert(alert Alert) {
    logger.Info("Processing alert",
        "status", alert.Status,
        "alertname", alert.Labels["alertname"],
        "severity", alert.Labels["severity"],
        "summary", alert.Annotations["summary"],
    )

    // Custom alert processing logic
    switch alert.Labels["alertname"] {
    case "HighErrorRate":
        handleHighErrorRate(alert)
    case "HighMemoryUsage":
        handleHighMemoryUsage(alert)
    case "PluginDown":
        handlePluginDown(alert)
    }
}

func handleHighErrorRate(alert Alert) {
    // Implement automatic error rate response
    // e.g., restart failing plugins, enable debug mode
    logger.Warn("High error rate detected, enabling debug mode")
    // Enable debug logging
    // Restart problematic plugins
}
```

## 📊 Dashboard Setup

### Grafana Dashboard Configuration

```json
{
  "dashboard": {
    "id": null,
    "title": "OpenCode Monitoring",
    "tags": ["opencode", "monitoring"],
    "timezone": "browser",
    "panels": [
      {
        "id": 1,
        "title": "Plugin Executions",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(opencode_plugin_executions_total[5m])",
            "legendFormat": "{{ plugin_name }} ({{ status }})"
          }
        ],
        "gridPos": {"h": 8, "w": 12, "x": 0, "y": 0},
        "yAxes": [
          {
            "label": "Executions/sec",
            "min": 0
          }
        ]
      },
      {
        "id": 2,
        "title": "Plugin Response Times",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(opencode_plugin_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          },
          {
            "expr": "histogram_quantile(0.50, rate(opencode_plugin_duration_seconds_bucket[5m]))",
            "legendFormat": "50th percentile"
          }
        ],
        "gridPos": {"h": 8, "w": 12, "x": 12, "y": 0},
        "yAxes": [
          {
            "label": "Seconds",
            "min": 0
          }
        ]
      },
      {
        "id": 3,
        "title": "System Resources",
        "type": "graph",
        "targets": [
          {
            "expr": "opencode_memory_usage_bytes / 1024 / 1024",
            "legendFormat": "Memory (MB)"
          },
          {
            "expr": "rate(opencode_cpu_usage_seconds_total[5m]) * 100",
            "legendFormat": "CPU (%)"
          }
        ],
        "gridPos": {"h": 8, "w": 24, "x": 0, "y": 8}
      },
      {
        "id": 4,
        "title": "Error Rate",
        "type": "singlestat",
        "targets": [
          {
            "expr": "rate(opencode_plugin_executions_total{status=\"error\"}[5m]) * 100",
            "legendFormat": "Error Rate %"
          }
        ],
        "gridPos": {"h": 4, "w": 6, "x": 0, "y": 16},
        "thresholds": "5,10",
        "colorBackground": true
      },
      {
        "id": 5,
        "title": "Active Plugins",
        "type": "singlestat",
        "targets": [
          {
            "expr": "opencode_active_plugins",
            "legendFormat": "Active Plugins"
          }
        ],
        "gridPos": {"h": 4, "w": 6, "x": 6, "y": 16}
      }
    ],
    "time": {
      "from": "now-1h",
      "to": "now"
    },
    "refresh": "5s"
  }
}
```

### Custom Dashboard Widgets

```typescript
// Custom dashboard component
import React, { useEffect, useState } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend } from 'recharts';

interface MetricData {
  timestamp: string;
  value: number;
  plugin?: string;
}

export const PluginPerformanceDashboard: React.FC = () => {
  const [responseTimeData, setResponseTimeData] = useState<MetricData[]>([]);
  const [errorRateData, setErrorRateData] = useState<MetricData[]>([]);

  useEffect(() => {
    const fetchMetrics = async () => {
      try {
        // Fetch response time metrics
        const responseTimeResponse = await fetch('/api/metrics/response-time');
        const responseTimeMetrics = await responseTimeResponse.json();
        setResponseTimeData(responseTimeMetrics);

        // Fetch error rate metrics
        const errorRateResponse = await fetch('/api/metrics/error-rate');
        const errorRateMetrics = await errorRateResponse.json();
        setErrorRateData(errorRateMetrics);
      } catch (error) {
        console.error('Failed to fetch metrics:', error);
      }
    };

    fetchMetrics();
    const interval = setInterval(fetchMetrics, 5000); // Update every 5 seconds

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="dashboard">
      <div className="dashboard-section">
        <h3>Plugin Response Times</h3>
        <LineChart width={800} height={300} data={responseTimeData}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="timestamp" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line type="monotone" dataKey="value" stroke="#8884d8" />
        </LineChart>
      </div>

      <div className="dashboard-section">
        <h3>Error Rate</h3>
        <div className="metric-cards">
          {errorRateData.map((metric, index) => (
            <div key={index} className={`metric-card ${metric.value > 5 ? 'warning' : 'normal'}`}>
              <div className="metric-value">{metric.value.toFixed(2)}%</div>
              <div className="metric-label">{metric.plugin || 'Overall'}</div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
```

## 🐛 Debugging Techniques

### Debug Mode Activation

```go
// Debug configuration
type DebugConfig struct {
    Enabled        bool   `yaml:"enabled"`
    LogLevel       string `yaml:"log_level"`
    ProfileEnabled bool   `yaml:"profile_enabled"`
    TraceEnabled   bool   `yaml:"trace_enabled"`
}

func EnableDebugMode(config *DebugConfig) {
    if config.Enabled {
        // Set log level to debug
        logger.SetLevel(slog.LevelDebug)
        
        // Enable profiling
        if config.ProfileEnabled {
            go func() {
                log.Println(http.ListenAndServe("localhost:6060", nil))
            }()
        }
        
        // Enable detailed tracing
        if config.TraceEnabled {
            otel.SetTracerProvider(
                trace.NewTracerProvider(
                    trace.WithSampler(trace.AlwaysSample()),
                ),
            )
        }
        
        logger.Info("Debug mode enabled",
            "log_level", config.LogLevel,
            "profile_enabled", config.ProfileEnabled,
            "trace_enabled", config.TraceEnabled,
        )
    }
}
```

### Performance Profiling

```go
package profiling

import (
    "context"
    "runtime"
    "runtime/pprof"
    "time"
)

// CPU profiling
func ProfileCPU(duration time.Duration, filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    if err := pprof.StartCPUProfile(f); err != nil {
        return err
    }
    defer pprof.StopCPUProfile()

    time.Sleep(duration)
    return nil
}

// Memory profiling
func ProfileMemory(filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    runtime.GC() // Get up-to-date statistics
    return pprof.WriteHeapProfile(f)
}

// Goroutine profiling
func ProfileGoroutines(filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    return pprof.Lookup("goroutine").WriteTo(f, 0)
}

// Automated profiling middleware
func ProfilingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("X-Profile-Request") == "true" {
            // Start CPU profiling for this request
            filename := fmt.Sprintf("profile-%d.prof", time.Now().Unix())
            go ProfileCPU(10*time.Second, filename)
            
            logger.Info("Profiling enabled for request",
                "method", r.Method,
                "path", r.URL.Path,
                "profile_file", filename,
            )
        }
        
        next.ServeHTTP(w, r)
    })
}
```

### Request Tracing

```typescript
import { AsyncLocalStorage } from 'async_hooks';
import { v4 as uuidv4 } from 'uuid';

interface RequestContext {
  requestId: string;
  timestamp: Date;
  path: string;
  method: string;
  userId?: string;
}

const requestContext = new AsyncLocalStorage<RequestContext>();

export class RequestTracer {
  static middleware(req: any, res: any, next: any) {
    const context: RequestContext = {
      requestId: uuidv4(),
      timestamp: new Date(),
      path: req.path,
      method: req.method,
      userId: req.user?.id
    };

    // Add request ID to response headers
    res.setHeader('X-Request-ID', context.requestId);

    // Run in async context
    requestContext.run(context, () => {
      logger.info('Request started', {
        request_id: context.requestId,
        method: context.method,
        path: context.path,
        user_id: context.userId
      });

      // Measure request duration
      const startTime = process.hrtime.bigint();
      
      res.on('finish', () => {
        const duration = Number(process.hrtime.bigint() - startTime) / 1000000; // Convert to ms
        
        logger.info('Request completed', {
          request_id: context.requestId,
          status_code: res.statusCode,
          duration_ms: duration
        });
      });

      next();
    });
  }

  static getContext(): RequestContext | undefined {
    return requestContext.getStore();
  }

  static trace<T>(operation: string, fn: () => Promise<T>): Promise<T> {
    const context = this.getContext();
    if (!context) {
      return fn();
    }

    const startTime = process.hrtime.bigint();
    
    logger.debug('Operation started', {
      request_id: context.requestId,
      operation
    });

    return fn()
      .then(result => {
        const duration = Number(process.hrtime.bigint() - startTime) / 1000000;
        logger.debug('Operation completed', {
          request_id: context.requestId,
          operation,
          duration_ms: duration
        });
        return result;
      })
      .catch(error => {
        const duration = Number(process.hrtime.bigint() - startTime) / 1000000;
        logger.error('Operation failed', {
          request_id: context.requestId,
          operation,
          duration_ms: duration,
          error: error.message
        });
        throw error;
      });
  }
}

// Usage in Express app
app.use(RequestTracer.middleware);

app.get('/api/plugins/:name/execute', async (req, res) => {
  const { name } = req.params;
  
  try {
    const result = await RequestTracer.trace('plugin_execution', async () => {
      return await pluginManager.execute(name, req.body.args);
    });
    
    res.json(result);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});
```

## 🔧 Plan-Specific Monitoring

### Plan 1: Plugin Ecosystem Architecture

**Key Metrics**:
- Plugin discovery time
- Plugin isolation health
- Community plugin adoption rates
- HashiCorp go-plugin performance

**Monitoring Setup**:
```go
// Plugin ecosystem specific metrics
var (
    pluginDiscoveryDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name: "opencode_plugin_discovery_duration_seconds",
            Help: "Time spent discovering plugins",
        },
    )
    
    pluginIsolationFailures = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "opencode_plugin_isolation_failures_total",
            Help: "Number of plugin isolation failures",
        },
    )
    
    communityPluginUsage = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "opencode_community_plugin_usage_total",
            Help: "Usage of community plugins",
        },
        []string{"plugin_name", "author"},
    )
)
```

### Plan 2: SuperClaude Integration

**Key Metrics**:
- Template rendering performance (<50ms target)
- Live-reload success rate
- Persona activation accuracy
- Context collection time

**Monitoring Setup**:
```typescript
const templateRenderingTime = new client.Histogram({
  name: 'opencode_template_rendering_seconds',
  help: 'Time spent rendering templates',
  buckets: [0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0]
});

const liveReloadEvents = new client.Counter({
  name: 'opencode_live_reload_events_total',
  help: 'Number of live reload events',
  labelNames: ['status', 'template_type']
});

const personaActivations = new client.Counter({
  name: 'opencode_persona_activations_total',
  help: 'Number of persona activations',
  labelNames: ['persona_type', 'accuracy']
});
```

### Plan 3: MCP Bridge + Plugin + EDA

**Key Metrics**:
- MCP standard compliance
- Event system throughput
- Multi-tier architecture latency
- Risk management effectiveness

**Monitoring Setup**:
```go
// Enterprise-grade monitoring
var (
    mcpComplianceScore = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "opencode_mcp_compliance_score",
            Help: "MCP standard compliance score (0-100)",
        },
    )
    
    eventSystemThroughput = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "opencode_events_processed_total",
            Help: "Number of events processed",
        },
        []string{"event_type", "tier"},
    )
    
    tierLatency = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "opencode_tier_latency_seconds",
            Help: "Latency between architectural tiers",
        },
        []string{"from_tier", "to_tier"},
    )
)
```

### Plan 4: MCP Plugin EDA Architecture

**Key Metrics**:
- NATS JetStream performance
- Event-driven scalability
- Edge computing readiness
- Mutual TLS success rate

### Plan 5: SuperClaude Implementation

**Key Metrics**:
- Comprehensive monitoring coverage
- Error handling effectiveness
- Configuration management health
- Enterprise feature adoption

## 📚 Additional Resources

### Tools and Platforms

- **Prometheus**: Metrics collection and alerting
- **Grafana**: Visualization and dashboards
- **Jaeger**: Distributed tracing
- **OpenTelemetry**: Observability framework
- **ELK Stack**: Log aggregation and analysis
- **PagerDuty**: Incident response

### Best Practices

1. **Start Simple**: Begin with basic metrics and expand gradually
2. **Monitor User Experience**: Focus on metrics that impact users
3. **Automate Responses**: Use alerts to trigger automated responses where possible
4. **Regular Reviews**: Regularly review and update monitoring strategies
5. **Documentation**: Document all monitoring procedures and runbooks

### Troubleshooting Guides

- **High CPU Usage**: Check for infinite loops, optimize algorithms
- **Memory Leaks**: Use profiling tools, check for unreleased resources
- **Network Issues**: Monitor connection pools, check for timeouts
- **Plugin Failures**: Review plugin logs, check isolation boundaries

---

This monitoring guide provides a comprehensive foundation for observability across all OpenCode-SuperClaude integration plans. Adapt the specific implementations based on your chosen plan and requirements. 📊