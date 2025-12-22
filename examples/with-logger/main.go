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

func (l *CustomLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

func (l *CustomLogger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l *CustomLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

func (l *CustomLogger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

func (l *CustomLogger) With(args ...interface{}) scheduler.Logger {
	return &CustomLogger{logger: l.logger.With(args...)}
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
