package main

import (
	"errors"
	"fmt"
	"math/rand"

	scheduler "github.com/donnigundala/dg-scheduler"
)

func main() {
	// Configure with custom error handler
	config := &scheduler.Config{
		ErrorHandler: func(name string, err error) {
			// Custom error handling logic
			fmt.Printf("⚠️  Schedule '%s' failed: %v\n", name, err)

			// You could:
			// - Send to error tracking (Sentry, Rollbar)
			// - Send alerts (PagerDuty, Slack)
			// - Increment metrics
			// - Log to external service
		},
	}

	s := scheduler.NewWithConfig(nil, config)
	s.Start()
	defer s.Stop()

	// Schedule a task that sometimes fails
	s.Schedule("*/1 * * * *", "flaky-task", func() error {
		if rand.Intn(2) == 0 {
			return errors.New("random failure")
		}
		fmt.Println("✓ Task succeeded")
		return nil
	})

	// Schedule a task that always fails
	s.Schedule("*/2 * * * *", "failing-task", func() error {
		return errors.New("this task always fails")
	})

	fmt.Println("Scheduler with error handler started.")
	fmt.Println("Watch for custom error messages...")
	select {}
}
