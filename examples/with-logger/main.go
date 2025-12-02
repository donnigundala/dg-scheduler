package main

import (
	"fmt"
	"log/slog"
	"os"

	scheduler "github.com/donnigundala/dg-scheduler"
)

// CustomLogger wraps slog to implement scheduler.Logger interface
type CustomLogger struct {
	logger *slog.Logger
}

func (l *CustomLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, keysAndValues...)
}

func (l *CustomLogger) Error(msg string, err error, keysAndValues ...interface{}) {
	args := append([]interface{}{"error", err}, keysAndValues...)
	l.logger.Error(msg, args...)
}

func (l *CustomLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warn(msg, keysAndValues...)
}

func main() {
	// Create structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Configure scheduler with custom logger
	config := &scheduler.Config{
		Logger: &CustomLogger{logger: logger},
	}

	s := scheduler.NewWithConfig(nil, config)
	s.Start()
	defer s.Stop()

	// Schedule tasks
	s.Schedule("*/1 * * * *", "task1", func() error {
		fmt.Println("Task 1 executed")
		return nil
	})

	s.Schedule("*/2 * * * *", "task2", func() error {
		fmt.Println("Task 2 executed")
		return nil
	})

	fmt.Println("Scheduler with custom logger started. Check JSON logs above.")
	select {}
}
