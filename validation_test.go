package scheduler

import (
	"strings"
	"testing"
	"time"
)

func TestValidateCronExpression(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr bool
	}{
		{"valid every minute", "* * * * *", false},
		{"valid every 5 minutes", "*/5 * * * *", false},
		{"valid specific time", "0 9 * * 1", false},
		{"empty expression", "", true},
		{"invalid format", "invalid", true},
		{"too many fields", "* * * * * *", true},
		{"invalid minute", "60 * * * *", true},
		{"invalid hour", "* 25 * * *", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCronExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCronExpression() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateScheduleName(t *testing.T) {
	tests := []struct {
		name     string
		schedule string
		wantErr  bool
	}{
		{"valid name", "my-schedule", false},
		{"valid with underscore", "my_schedule", false},
		{"valid with numbers", "schedule123", false},
		{"empty name", "", true},
		{"too long", strings.Repeat("a", 101), true},
		{"with newline", "schedule\nname", true},
		{"with tab", "schedule\tname", true},
		{"with carriage return", "schedule\rname", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateScheduleName(tt.schedule)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateScheduleName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseCronExpression(t *testing.T) {
	schedule, err := ParseCronExpression("*/5 * * * *")
	if err != nil {
		t.Fatalf("ParseCronExpression() error = %v", err)
	}

	if schedule == nil {
		t.Error("Expected non-nil schedule")
	}

	// Test that we can get next run time
	now := time.Now()
	next := schedule.Next(now)

	if next.Before(now) {
		t.Error("Next run time should be after now")
	}
}

func TestParseCronExpression_Invalid(t *testing.T) {
	_, err := ParseCronExpression("invalid")
	if err == nil {
		t.Error("Expected error for invalid cron expression")
	}
}

func TestScheduler_ValidationInSchedule(t *testing.T) {
	s := New(nil)
	defer s.Stop()

	// Test empty name
	err := s.Schedule("* * * * *", "", func() error { return nil })
	if err == nil {
		t.Error("Expected error for empty schedule name")
	}

	// Test invalid cron
	err = s.Schedule("invalid", "test", func() error { return nil })
	if err == nil {
		t.Error("Expected error for invalid cron expression")
	}

	// Test valid schedule
	err = s.Schedule("*/5 * * * *", "valid", func() error { return nil })
	if err != nil {
		t.Errorf("Unexpected error for valid schedule: %v", err)
	}
}
