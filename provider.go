package dgscheduler

import (
	"reflect"

	"github.com/donnigundala/dg-core/contracts/foundation"
)

// SchedulerServiceProvider implements the PluginProvider interface.
// This provides a simple, plug-and-play integration for applications
// that don't need custom configuration.
//
// For advanced use cases requiring custom adapters or configuration,
// use the library functions (New, NewWithConfig) directly.
type SchedulerServiceProvider struct {
	// If not provided, defaults will be used
	Config *Config
}

// NewSchedulerServiceProvider creates a new scheduler service provider.
func NewSchedulerServiceProvider(config *Config) *SchedulerServiceProvider {
	return &SchedulerServiceProvider{
		Config: config,
	}
}

// Name returns the name of the plugin.
func (p *SchedulerServiceProvider) Name() string {
	return Binding
}

// Version returns the version of the plugin.
func (p *SchedulerServiceProvider) Version() string {
	return Version
}

// Dependencies returns the list of dependencies.
// Scheduler has no required dependencies, but can optionally integrate with queue.
func (p *SchedulerServiceProvider) Dependencies() []string {
	return []string{}
}

// Register registers the scheduler service provider.
func (p *SchedulerServiceProvider) Register(app foundation.Application) error {
	app.Singleton(Binding, func() (interface{}, error) {
		// Use provided config or default
		cfg := p.Config
		if cfg == nil {
			cfg = DefaultConfig()
		}

		// Try to resolve logger from container if not already configured
		if cfg.Logger == nil {
			if loggerInstance, err := app.Make("logger"); err == nil {
				// Adapt the logger to scheduler.Logger interface
				if adapted, ok := loggerInstance.(interface {
					Debug(msg string, args ...interface{})
					Info(msg string, args ...interface{})
					Warn(msg string, args ...interface{})
					Error(msg string, args ...interface{})
				}); ok {
					cfg.Logger = &loggerAdapter{logger: adapted}
				}
			}
		}

		// Create the scheduler instance
		schedulerInstance := NewWithConfig(nil, cfg)
		return schedulerInstance, nil
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
		Debug(msg string, args ...interface{})
		Info(msg string, args ...interface{})
		Warn(msg string, args ...interface{})
		Error(msg string, args ...interface{})
	}
}

func (l *loggerAdapter) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

func (l *loggerAdapter) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l *loggerAdapter) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

func (l *loggerAdapter) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

func (l *loggerAdapter) With(args ...interface{}) Logger {
	// Try to call With(args...) via reflection to support different return types
	v := reflect.ValueOf(l.logger)
	m := v.MethodByName("With")
	if m.IsValid() {
		valArgs := make([]reflect.Value, len(args))
		for i, arg := range args {
			valArgs[i] = reflect.ValueOf(arg)
		}
		results := m.Call(valArgs)
		if len(results) == 1 {
			if nextLogger, ok := results[0].Interface().(interface {
				Debug(msg string, args ...interface{})
				Info(msg string, args ...interface{})
				Warn(msg string, args ...interface{})
				Error(msg string, args ...interface{})
			}); ok {
				return &loggerAdapter{logger: nextLogger}
			}
		}
	}
	return l
}
