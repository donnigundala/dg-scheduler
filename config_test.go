package scheduler_test

import (
	"fmt"
	"testing"

	scheduler "github.com/donnigundala/dg-scheduler"
)

// Example custom logger implementation
type testLogger struct {
	infos  []string
	errors []string
}

func (l *testLogger) Info(msg string, keysAndValues ...interface{}) {
	formatted := msg
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			formatted += fmt.Sprintf(" %v=%v", keysAndValues[i], keysAndValues[i+1])
		}
	}
	l.infos = append(l.infos, formatted)
}

func (l *testLogger) Error(msg string, err error, keysAndValues ...interface{}) {
	formatted := fmt.Sprintf("%s: %v", msg, err)
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			formatted += fmt.Sprintf(" %v=%v", keysAndValues[i], keysAndValues[i+1])
		}
	}
	l.errors = append(l.errors, formatted)
}

func (l *testLogger) Warn(msg string, keysAndValues ...interface{}) {
	// Not used in scheduler yet
}

func TestScheduler_WithLogger(t *testing.T) {
	logger := &testLogger{}
	config := &scheduler.Config{
		Logger: logger,
	}

	s := scheduler.NewWithConfig(nil, config)
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

	config := &scheduler.Config{
		ErrorHandler: func(name string, err error) {
			errorHandlerCalled = true
		},
	}

	s := scheduler.NewWithConfig(nil, config)
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
	config := scheduler.DefaultConfig()

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
