package dgscheduler

// Logger is the interface for structured logging.
// Implement this interface to integrate with your logging system.
type Logger interface {
	// Debug logs a debug message
	Debug(msg string, args ...interface{})

	// Info logs an informational message
	Info(msg string, args ...interface{})

	// Warn logs a warning message
	Warn(msg string, args ...interface{})

	// Error logs an error message
	Error(msg string, args ...interface{})

	// With returns a new logger with the given fields added
	With(args ...interface{}) Logger
}

// Config holds configuration for the scheduler.
type Config struct {
	// ErrorHandler is called when a scheduled job returns an error.
	// If nil, errors are logged using the Logger (or printed if no logger).
	ErrorHandler func(name string, err error)

	// Logger is used for structured logging.
	// If nil, falls back to fmt.Printf.
	Logger Logger
}

// DefaultConfig returns a config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		ErrorHandler: nil,
		Logger:       nil,
	}
}
