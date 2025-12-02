package scheduler

import (
	"fmt"
	"strings"

	"github.com/robfig/cron/v3"
)

// ValidateCronExpression validates a cron expression.
// Returns nil if valid, error with helpful message if invalid.
func ValidateCronExpression(expr string) error {
	if expr == "" {
		return fmt.Errorf("cron expression cannot be empty")
	}

	// Parse the expression
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(expr)
	if err != nil {
		return fmt.Errorf("invalid cron expression '%s': %w", expr, err)
	}

	return nil
}

// ValidateScheduleName validates a schedule name.
// Returns nil if valid, error if invalid.
func ValidateScheduleName(name string) error {
	if name == "" {
		return fmt.Errorf("schedule name cannot be empty")
	}

	if len(name) > 100 {
		return fmt.Errorf("schedule name too long (max 100 characters): %d", len(name))
	}

	// Check for invalid characters
	if strings.ContainsAny(name, "\n\r\t") {
		return fmt.Errorf("schedule name contains invalid characters (newlines, tabs)")
	}

	return nil
}

// ParseCronExpression parses a cron expression and returns the schedule.
// This is useful for getting the next run time before adding the schedule.
func ParseCronExpression(expr string) (cron.Schedule, error) {
	if err := ValidateCronExpression(expr); err != nil {
		return nil, err
	}

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	return parser.Parse(expr)
}
