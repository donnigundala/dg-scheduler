package scheduler

import (
	"github.com/donnigundala/dg-core/contracts/foundation"
)

// SchedulerServiceProvider implements the PluginProvider interface.
// This provides a simple, plug-and-play integration for applications
// that don't need custom configuration.
//
// For advanced use cases requiring custom adapters or configuration,
// use the library functions (New, NewWithConfig) directly.
type SchedulerServiceProvider struct {
	// Config holds optional scheduler configuration
	// If not provided, defaults will be used
	Config *Config
}

// Name returns the name of the plugin.
func (p *SchedulerServiceProvider) Name() string {
	return "scheduler"
}

// Version returns the version of the plugin.
func (p *SchedulerServiceProvider) Version() string {
	return "1.1.0"
}

// Dependencies returns the list of dependencies.
// Scheduler has no required dependencies, but can optionally integrate with queue.
func (p *SchedulerServiceProvider) Dependencies() []string {
	return []string{}
}

// Register registers the scheduler service provider.
func (p *SchedulerServiceProvider) Register(app foundation.Application) error {
	// Use provided config or default
	cfg := p.Config
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Register the scheduler as a singleton
	app.Singleton("scheduler", func() (interface{}, error) {
		// Try to resolve queue dispatcher (optional)
		var dispatcher Dispatcher
		if queueInstance, err := app.Make("queue"); err == nil {
			// Queue is available, use it as dispatcher
			// We use interface{} to avoid importing dg-queue
			if queueDispatcher, ok := queueInstance.(Dispatcher); ok {
				dispatcher = queueDispatcher
			}
		}

		// Try to resolve logger (optional)
		if cfg.Logger == nil {
			if loggerInstance, err := app.Make("logger"); err == nil {
				// Adapt dg-core logger to scheduler.Logger interface
				if logger, ok := loggerInstance.(interface {
					Info(msg string, keysAndValues ...interface{})
					Error(msg string, keysAndValues ...interface{})
					Warn(msg string, keysAndValues ...interface{})
				}); ok {
					cfg.Logger = &loggerAdapter{logger: logger}
				}
			}
		}

		// Create scheduler with optional dispatcher and config
		return NewWithConfig(dispatcher, cfg), nil
	})

	return nil
}

// Boot boots the scheduler service provider.
func (p *SchedulerServiceProvider) Boot(app foundation.Application) error {
	// NOTE: Scheduler implements Runnable interface
	// It will be automatically started by app.StartServices()
	// No need to verify resolution here - it will be resolved when needed

	return nil
}

// Shutdown gracefully stops the scheduler.
func (p *SchedulerServiceProvider) Shutdown(app foundation.Application) error {
	schedulerInstance, err := app.Make("scheduler")
	if err != nil {
		return nil // Scheduler not initialized
	}

	s := schedulerInstance.(*Scheduler)
	s.Stop()

	return nil
}

// loggerAdapter adapts a generic logger to scheduler.Logger interface.
type loggerAdapter struct {
	logger interface {
		Info(msg string, keysAndValues ...interface{})
		Error(msg string, keysAndValues ...interface{})
		Warn(msg string, keysAndValues ...interface{})
	}
}

func (l *loggerAdapter) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, keysAndValues...)
}

func (l *loggerAdapter) Error(msg string, err error, keysAndValues ...interface{}) {
	args := append([]interface{}{"error", err}, keysAndValues...)
	l.logger.Error(msg, args...)
}

func (l *loggerAdapter) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warn(msg, keysAndValues...)
}
