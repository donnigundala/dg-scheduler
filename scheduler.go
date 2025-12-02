package scheduler

import (
	"context"
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

// Dispatcher is the minimal interface needed for queue integration.
// This allows the scheduler to work with any queue implementation,
// or to work standalone without a queue.
//
// The interface matches dg-queue's Queue.Dispatch signature but uses
// interface{} instead of *Job to avoid importing dg-queue.
type Dispatcher interface {
	Dispatch(name string, payload interface{}) (interface{}, error)
}

// Scheduler manages scheduled jobs using cron syntax.
// It can work standalone or dispatch jobs to a queue.
type Scheduler struct {
	cron       *cron.Cron
	dispatcher Dispatcher // Optional - can be nil
	config     *Config
	entries    map[string]cron.EntryID
	mu         sync.RWMutex
}

// New creates a new scheduler with default configuration.
// The dispatcher parameter is optional - pass nil to use scheduler standalone.
// If dispatcher is provided, you can use ScheduleJob() to dispatch to the queue.
// If dispatcher is nil, use Schedule() with custom handlers.
func New(dispatcher Dispatcher) *Scheduler {
	return NewWithConfig(dispatcher, DefaultConfig())
}

// NewWithConfig creates a new scheduler with custom configuration.
func NewWithConfig(dispatcher Dispatcher, config *Config) *Scheduler {
	if config == nil {
		config = DefaultConfig()
	}

	return &Scheduler{
		cron:       cron.New(),
		dispatcher: dispatcher,
		config:     config,
		entries:    make(map[string]cron.EntryID),
	}
}

// Schedule schedules a custom handler using cron syntax.
// This method works with or without a dispatcher.
//
// cronExpr: Cron expression (e.g., "*/5 * * * *" for every 5 minutes)
// name: Unique name for this scheduled job
// handler: Function to execute on schedule
func (s *Scheduler) Schedule(cronExpr, name string, handler func() error) error {
	// Validate inputs
	if err := ValidateScheduleName(name); err != nil {
		return err
	}

	if err := ValidateCronExpression(cronExpr); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if already scheduled
	if _, exists := s.entries[name]; exists {
		return fmt.Errorf("schedule '%s' already exists", name)
	}

	// Add to cron
	entryID, err := s.cron.AddFunc(cronExpr, func() {
		if err := handler(); err != nil {
			s.handleError(name, err)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to add schedule: %w", err)
	}

	s.entries[name] = entryID

	// Log schedule creation
	s.logInfo("Schedule created", "name", name, "cron", cronExpr)

	return nil
}

// ScheduleJob schedules a job to be dispatched to the queue on a cron schedule.
// This is a convenience method that requires a dispatcher to be provided during New().
//
// Returns an error if the scheduler was created without a dispatcher.
func (s *Scheduler) ScheduleJob(cronExpr, jobName string, payload interface{}) error {
	if s.dispatcher == nil {
		return fmt.Errorf("scheduler was created without a dispatcher - use Schedule() with a custom handler instead")
	}

	return s.Schedule(cronExpr, "schedule_"+jobName, func() error {
		_, err := s.dispatcher.Dispatch(jobName, payload)
		return err
	})
}

// Remove removes a scheduled job.
func (s *Scheduler) Remove(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entryID, exists := s.entries[name]
	if !exists {
		return fmt.Errorf("schedule '%s' not found", name)
	}

	s.cron.Remove(entryID)
	delete(s.entries, name)

	s.logInfo("Schedule removed", "name", name)

	return nil
}

// Start starts the scheduler.
func (s *Scheduler) Start() {
	s.cron.Start()
	s.logInfo("Scheduler started", "schedules", len(s.entries))
}

// Stop stops the scheduler gracefully.
// Returns a context that will be done when the scheduler has stopped.
func (s *Scheduler) Stop() context.Context {
	s.logInfo("Scheduler stopping", "schedules", len(s.entries))
	return s.cron.Stop()
}

// Count returns the number of scheduled jobs.
func (s *Scheduler) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}

// handleError handles errors from scheduled jobs.
func (s *Scheduler) handleError(name string, err error) {
	// Use custom error handler if provided
	if s.config.ErrorHandler != nil {
		s.config.ErrorHandler(name, err)
		return
	}

	// Otherwise log the error
	s.logError("Scheduled job failed", err, "name", name)
}

// logInfo logs an informational message.
func (s *Scheduler) logInfo(msg string, keysAndValues ...interface{}) {
	if s.config.Logger != nil {
		s.config.Logger.Info(msg, keysAndValues...)
	} else {
		// Fallback to fmt.Printf
		fmt.Printf("[Scheduler] %s", msg)
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				fmt.Printf(" %v=%v", keysAndValues[i], keysAndValues[i+1])
			}
		}
		fmt.Println()
	}
}

// logError logs an error message.
func (s *Scheduler) logError(msg string, err error, keysAndValues ...interface{}) {
	if s.config.Logger != nil {
		s.config.Logger.Error(msg, err, keysAndValues...)
	} else {
		// Fallback to fmt.Printf
		fmt.Printf("[Scheduler] ERROR: %s: %v", msg, err)
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				fmt.Printf(" %v=%v", keysAndValues[i], keysAndValues[i+1])
			}
		}
		fmt.Println()
	}
}
