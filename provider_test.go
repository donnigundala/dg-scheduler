package scheduler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchedulerServiceProvider_Name(t *testing.T) {
	provider := &SchedulerServiceProvider{}
	assert.Equal(t, "scheduler", provider.Name())
}

func TestSchedulerServiceProvider_Version(t *testing.T) {
	provider := &SchedulerServiceProvider{}
	assert.Equal(t, "1.2.0", provider.Version())
}

func TestSchedulerServiceProvider_Dependencies(t *testing.T) {
	provider := &SchedulerServiceProvider{}
	deps := provider.Dependencies()

	assert.NotNil(t, deps)
	assert.Empty(t, deps, "dg-scheduler should have no required dependencies")
}

func TestSchedulerServiceProvider_ConfigDefaults(t *testing.T) {
	provider := &SchedulerServiceProvider{}
	assert.Nil(t, provider.Config, "Config should be nil by default")
}

func TestSchedulerServiceProvider_CustomConfig(t *testing.T) {
	customConfig := &Config{
		ErrorHandler: func(name string, err error) {},
	}

	provider := &SchedulerServiceProvider{
		Config: customConfig,
	}

	assert.NotNil(t, provider.Config)
	assert.NotNil(t, provider.Config.ErrorHandler)
}
