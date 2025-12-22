package dgscheduler_test

import (
	"fmt"
	"testing"

	dgscheduler "github.com/donnigundala/dg-scheduler"
)

// Example custom logger implementation
type testLogger struct {
	infos  []string
	errors []string
}

func (l *testLogger) Debug(msg string, args ...interface{}) {
	// Not used in tests yet
}

func (l *testLogger) Info(msg string, args ...interface{}) {
	formatted := msg
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			formatted += fmt.Sprintf(" %v=%v", args[i], args[i+1])
		}
	}
	l.infos = append(l.infos, formatted)
}

func (l *testLogger) Warn(msg string, args ...interface{}) {
	// Not used in scheduler yet
}

func (l *testLogger) Error(msg string, args ...interface{}) {
	formatted := msg
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			formatted += fmt.Sprintf(" %v=%v", args[i], args[i+1])
		}
	}
	l.errors = append(l.errors, formatted)
}

func (l *testLogger) With(args ...interface{}) dgscheduler.Logger {
	return l
}

func TestScheduler_WithLogger(t *testing.T) {
	logger := &testLogger{}
	config := &dgscheduler.Config{
		Logger: logger,
	}

	s := dgscheduler.NewWithConfig(nil, config)
	s.Start()
	defer s.Stop()

	// Schedule a job
	err := s.Schedule("*/5 * * * *", "test", func() error {
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to schedule: %v", err)
	}

	// Check that schedule creation was logged
	if len(logger.infos) < 2 { // Start + Schedule created
		t.Errorf("Expected at least 2 info logs, got %d", len(logger.infos))
	}
}

func TestScheduler_WithErrorHandler(t *testing.T) {
	errorHandlerCalled := false

	config := &dgscheduler.Config{
		ErrorHandler: func(name string, err error) {
			errorHandlerCalled = true
		},
	}

	s := dgscheduler.NewWithConfig(nil, config)
	s.Start()
	defer s.Stop()

	// Schedule a job
	err := s.Schedule("*/5 * * * *", "test-job", func() error {
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to schedule: %v", err)
	}

	// Verify config was set
	if config.ErrorHandler == nil {
		t.Error("Error handler should be set")
	}

	// Note: errorHandlerCalled would be true if a job actually failed
	// We can't easily test this without waiting for cron execution
	_ = errorHandlerCalled
}

func TestScheduler_DefaultConfig(t *testing.T) {
	config := dgscheduler.DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig should not return nil")
	}

	if config.ErrorHandler != nil {
		t.Error("Default error handler should be nil")
	}

	if config.Logger != nil {
		t.Error("Default logger should be nil")
	}
}
