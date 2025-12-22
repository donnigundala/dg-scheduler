package dgscheduler

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

// ValidateCronExpression validates a cron expression format.
// Returns an error if the expression is invalid.
func ValidateCronExpression(expr string) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(expr)
	if err != nil {
		return fmt.Errorf("invalid cron expression '%s': %w", expr, err)
	}
	return nil
}

// ValidateScheduleName validates a schedule name.
// Names must be non-empty and contain only alphanumeric characters, hyphens, and underscores.
func ValidateScheduleName(name string) error {
	if name == "" {
		return fmt.Errorf("schedule name cannot be empty")
	}

	if len(name) > 100 {
		return fmt.Errorf("schedule name too long (max 100 characters): %d", len(name))
	}

	// Check for valid characters (alphanumeric, hyphen, underscore)
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return fmt.Errorf("schedule name '%s' contains invalid character '%c'", name, char)
		}
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
