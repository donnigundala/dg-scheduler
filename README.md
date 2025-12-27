# dg-scheduler

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A lightweight, cron-based job scheduler for Go. Works standalone or integrates with any queue system.

## Features

- ⏰ **Cron Scheduling** - Schedule jobs using standard cron expressions
- 🔌 **Flexible Integration** - Works standalone OR with queue systems
- 🎯 **Simple API** - Easy to use, minimal configuration
- ✅ **Battle-tested** - Uses [robfig/cron/v3](https://github.com/robfig/cron) under the hood
- 🧹 **Graceful Shutdown** - Clean shutdown with context support
- 🚀 **Zero Dependencies** - No required dependencies (queue integration is optional)

## Installation

```bash
go get github.com/donnigundala/dg-scheduler
```

## Quick Start

### Standalone Usage (No Queue Required)

```go
package main

import (
    "fmt"
    
    "github.com/donnigundala/dg-scheduler"
)

func main() {
    // Create scheduler without queue
    s := scheduler.New(nil)
    s.Start()
    defer s.Stop()
    
    // Schedule with custom handler
    s.Schedule("*/5 * * * *", "cleanup", func() error {
        fmt.Println("Running cleanup...")
        return nil
    })
    
    // Schedule hourly report
    s.Schedule("0 * * * *", "report", func() error {
        fmt.Println("Generating report...")
        return nil
    })
}
```

```go
package main

import (
    "github.com/donnigundala/dg-core/foundation"
    "github.com/donnigundala/dg-scheduler"
    "github.com/donnigundala/dg-queue"
)

func main() {
    app := foundation.New(".")
    
    // 1. Register Queue (Optional, required for ScheduleJob)
    app.Register(dgqueue.NewQueueServiceProvider(nil))
    
    // 2. Register Scheduler
    app.Register(dgscheduler.NewSchedulerServiceProvider(nil))
    
    app.Start()
}
```

### Integration via InfrastructureSuite
In your `bootstrap/app.go`, the scheduler is typically registered directly:

```go
func InfrastructureSuite(workerMode bool) []foundation.ServiceProvider {
	return []foundation.ServiceProvider{
		dgqueue.NewQueueServiceProvider(nil),
		dgscheduler.NewSchedulerServiceProvider(nil),
		// ...
	}
}
```
## Cron Expression Format

```
┌───────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌───────────── day of month (1 - 31)
│ │ │ ┌───────────── month (1 - 12)
│ │ │ │ ┌───────────── day of week (0 - 6) (Sunday to Saturday)
│ │ │ │ │
* * * * *
```

### Examples

- `*/5 * * * *` - Every 5 minutes
- `0 * * * *` - Every hour
- `0 0 * * *` - Every day at midnight
- `0 9 * * 1` - Every Monday at 9 AM
- `*/15 9-17 * * 1-5` - Every 15 minutes during business hours (9-5, Mon-Fri)

## API Reference

### Creating a Scheduler

**Standalone (no queue):**
```go
scheduler := scheduler.New(nil)
```

**With queue integration:**
```go
scheduler := scheduler.New(queueInstance)
```

**With custom dispatcher:**
```go
// Implement the Dispatcher interface
type MyDispatcher struct{}

func (d *MyDispatcher) Dispatch(name string, payload interface{}) (interface{}, error) {
    // Your custom dispatch logic
    return nil, nil
}

scheduler := scheduler.New(&MyDispatcher{})
```

### Scheduling Jobs

**Schedule with custom handler (works with or without queue):**
```go
err := scheduler.Schedule(cronExpr, name string, handler func() error)
```

**Schedule job to be dispatched to queue (requires queue):**
```go
err := scheduler.ScheduleJob(cronExpr, jobName string, payload interface{})
// Returns error if scheduler was created without a dispatcher
```

### Managing Schedules

**Remove a scheduled job:**
```go
err := scheduler.Remove(name string)
```

**Get count of scheduled jobs:**
```go
count := scheduler.Count()
```

### Lifecycle

**Start the scheduler:**
```go
scheduler.Start()
```

**Stop the scheduler:**
```go
ctx := scheduler.Stop() // Returns context that's done when stopped
<-ctx.Done() // Wait for shutdown
```

## Configuration

### Custom Logger

Implement the `Logger` interface to integrate with your logging system:

```go
type MyLogger struct {
    logger *slog.Logger
}

func (l *MyLogger) Info(msg string, keysAndValues ...interface{}) {
    l.logger.Info(msg, keysAndValues...)
}

func (l *MyLogger) Error(msg string, err error, keysAndValues ...interface{}) {
    args := append([]interface{}{"error", err}, keysAndValues...)
    l.logger.Error(msg, args...)
}

func (l *MyLogger) Warn(msg string, keysAndValues ...interface{}) {
    l.logger.Warn(msg, keysAndValues...)
}

// Use custom logger
config := &scheduler.Config{
    Logger: &MyLogger{logger: slog.Default()},
}
s := scheduler.NewWithConfig(queueInstance, config)
```

### Custom Error Handler

Handle schedule errors your way:

```go
config := &scheduler.Config{
    ErrorHandler: func(name string, err error) {
        // Send to error tracking service
        sentry.CaptureException(err)
        
        // Log with context
        log.Error("Schedule failed",
            "schedule", name,
            "error", err,
        )
        
        // Send alert
        if isCritical(name) {
            sendPagerDutyAlert(name, err)
        }
    },
}

s := scheduler.NewWithConfig(nil, config)
```

### Complete Configuration Example

```go
config := &scheduler.Config{
    Logger: myLogger,
    ErrorHandler: func(name string, err error) {
        metrics.IncrementScheduleError(name)
        logger.Error("Schedule error", "name", name, "error", err)
    },
}

s := scheduler.NewWithConfig(queueInstance, config)
s.Start()
defer s.Stop()
```

## Validation

The scheduler validates all inputs before adding schedules:

### Cron Expression Validation

```go
// Validate before using
err := scheduler.ValidateCronExpression("*/5 * * * *")
if err != nil {
    log.Fatal(err) // Invalid cron expression
}

// Parse and get next run time
schedule, err := scheduler.ParseCronExpression("0 9 * * 1")
if err != nil {
    log.Fatal(err)
}
nextRun := schedule.Next(time.Now())
fmt.Println("Next run:", nextRun)
```

### Schedule Name Validation

```go
// Validates:
// - Not empty
// - Max 100 characters
// - No newlines, tabs, or control characters
err := scheduler.ValidateScheduleName("my-schedule")
if err != nil {
    log.Fatal(err)
}
```

### Automatic Validation

All schedules are automatically validated:

```go
// This will return an error if invalid
err := s.Schedule("invalid cron", "bad-name\n", handler)
// Error: schedule name contains invalid characters (newlines, tabs)
```

## Migration from dg-queue

If you were using `Manager.Schedule()` from dg-queue v1.x:

### Before
```go
q := queue.New(config)
q.Schedule("*/5 * * * *", "cleanup", handler)
q.Start(ctx)
```

### After
```go
// Queue for job processing
q := queue.New(config)
q.Start(ctx)

// Scheduler for cron jobs (separate package)
s := scheduler.New(q)
s.Schedule("*/5 * * * *", "cleanup", handler)
s.Start()
```

## Integration with dg-core

```go
// In your service provider
type SchedulerServiceProvider struct {
    scheduler *scheduler.Scheduler
}

func (p *SchedulerServiceProvider) Boot(app foundation.Application) error {
    queueInstance, _ := app.Make("queue")
    q := queueInstance.(queue.Queue)
    
    p.scheduler = scheduler.New(q)
    
    // Register your scheduled jobs
    p.scheduler.ScheduleJob("0 * * * *", "hourly-cleanup", nil)
    
    p.scheduler.Start()
    return nil
}

func (p *SchedulerServiceProvider) Shutdown(app foundation.Application) error {
    ctx := p.scheduler.Stop()
    <-ctx.Done()
    return nil
}
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Related Packages

- [dg-queue](https://github.com/donnigundala/dg-queue) - Queue system for job processing
- [dg-core](https://github.com/donnigundala/dg-core) - Core framework
- [dg-cache](https://github.com/donnigundala/dg-cache) - Cache abstraction
