package scheduler

import (
	"testing"

	"github.com/donnigundala/dg-core/foundation"
	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	app := foundation.New(".")
	cfg := DefaultConfig()

	// Create scheduler instance
	scheduler := NewWithConfig(nil, cfg)
	app.Instance("scheduler", scheduler)

	// Test Resolve
	s, err := Resolve(app)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, scheduler, s)
}

func TestResolve_Error(t *testing.T) {
	app := foundation.New(".")

	// Test Resolve without registration
	s, err := Resolve(app)
	assert.Error(t, err)
	assert.Nil(t, s)
}

func TestMustResolve(t *testing.T) {
	app := foundation.New(".")
	scheduler := NewWithConfig(nil, DefaultConfig())
	app.Instance("scheduler", scheduler)

	// Test MustResolve
	assert.NotPanics(t, func() {
		s := MustResolve(app)
		assert.NotNil(t, s)
	})
}

func TestMustResolve_Panic(t *testing.T) {
	app := foundation.New(".")

	// Test MustResolve panics without registration
	assert.Panics(t, func() {
		MustResolve(app)
	})
}

func TestInjectable(t *testing.T) {
	app := foundation.New(".")
	scheduler := NewWithConfig(nil, DefaultConfig())
	app.Instance("scheduler", scheduler)

	inject := NewInjectable(app)

	assert.NotPanics(t, func() {
		s := inject.Scheduler()
		assert.NotNil(t, s)
	})
}
