package dgscheduler

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	instrumentationName = "github.com/donnigundala/dg-scheduler"
)

// Observability wraps a scheduler handler with metrics.
func Observability(name string, handler func() error) func() error {
	meter := otel.GetMeterProvider().Meter(instrumentationName)

	executionCounter, err := meter.Int64Counter(
		"scheduler.job.execution.count",
		metric.WithDescription("Total number of scheduled job executions"),
		metric.WithUnit("{execution}"),
	)
	if err != nil {
		// Log error
	}

	durationHistogram, err := meter.Float64Histogram(
		"scheduler.job.execution.duration",
		metric.WithDescription("Duration of scheduled job executions"),
		metric.WithUnit("ms"),
	)
	if err != nil {
	}

	return func() error {
		start := time.Now()
		err := handler()
		duration := float64(time.Since(start).Milliseconds())

		ctx := context.Background()
		status := "success"
		if err != nil {
			status = "error"
		}

		attrs := metric.WithAttributes(
			attribute.String("job.name", name),
			attribute.String("job.status", status),
		)

		executionCounter.Add(ctx, 1, attrs)
		durationHistogram.Record(ctx, duration, attrs)

		return err
	}
}
