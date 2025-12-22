package dgscheduler

import (
	"fmt"

	"github.com/donnigundala/dg-core/contracts/foundation"
)

// Resolve resolves the Scheduler from the container.
func Resolve(app foundation.Application) (*Scheduler, error) {
	instance, err := app.Make("scheduler")
	if err != nil {
		return nil, err
	}
	return instance.(*Scheduler), nil
}

// MustResolve resolves the Scheduler or panics.
func MustResolve(app foundation.Application) *Scheduler {
	s, err := Resolve(app)
	if err != nil {
		panic(fmt.Sprintf("failed to resolve scheduler: %v", err))
	}
	return s
}

// Injectable can be embedded in structs to provide easy access to Scheduler.
type Injectable struct {
	app foundation.Application
}

// NewInjectable creates a new Injectable instance.
func NewInjectable(app foundation.Application) *Injectable {
	return &Injectable{
		app: app,
	}
}

// Scheduler returns the resolved Scheduler.
// It panics if the scheduler cannot be resolved.
func (i *Injectable) Scheduler() *Scheduler {
	return MustResolve(i.app)
}
