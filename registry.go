package dgscheduler

import (
	"fmt"
	"log/slog"
)

// Registry holds all registered scheduled jobs.
type Registry struct {
	jobs   []Job
	logger *slog.Logger
}

// NewRegistry creates a new job registry.
func NewRegistry(logger *slog.Logger) *Registry {
	return &Registry{
		jobs:   make([]Job, 0),
		logger: logger,
	}
}

// Register adds a job to the registry.
func (r *Registry) Register(job Job) {
	r.jobs = append(r.jobs, job)
	if r.logger != nil {
		r.logger.Debug("Job registered", "name", job.Name(), "schedule", job.Schedule(), "enabled", job.IsEnabled())
	}
}

// GetEnabledJobs returns all enabled jobs.
func (r *Registry) GetEnabledJobs() []Job {
	enabled := make([]Job, 0)
	for _, job := range r.jobs {
		if job.IsEnabled() {
			enabled = append(enabled, job)
		}
	}
	return enabled
}

// SchedulerEngine defines the interface for the underlying scheduler (e.g. cron).
type SchedulerEngine interface {
	Schedule(cronExpr, name string, handler func() error) error
}

// ScheduleAll schedules all enabled jobs with the provided engine.
func (r *Registry) ScheduleAll(engine SchedulerEngine) error {
	enabledJobs := r.GetEnabledJobs()

	for _, job := range enabledJobs {
		// job.Handle matches func() error signature
		if err := engine.Schedule(job.Schedule(), job.Name(), job.Handle); err != nil {
			return fmt.Errorf("failed to schedule job '%s': %w", job.Name(), err)
		}
		if r.logger != nil {
			r.logger.Info("Job scheduled", "name", job.Name(), "schedule", job.Schedule())
		}
	}

	if r.logger != nil {
		r.logger.Info("All jobs scheduled", "total", len(r.jobs), "enabled", len(enabledJobs))
	}
	return nil
}
