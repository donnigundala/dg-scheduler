package main

import (
	"fmt"

	scheduler "github.com/donnigundala/dg-scheduler"
)

func main() {
	// Create scheduler without queue (standalone)
	s := scheduler.New(nil)
	s.Start()
	defer s.Stop()

	// Schedule a simple task
	s.Schedule("*/1 * * * *", "hello", func() error {
		fmt.Println("Hello from scheduler!")
		return nil
	})

	// Schedule a cleanup task
	s.Schedule("*/2 * * * *", "cleanup", func() error {
		fmt.Println("Running cleanup...")
		// Your cleanup logic here
		return nil
	})

	// Keep running
	fmt.Println("Scheduler started. Press Ctrl+C to stop.")
	select {}
}
