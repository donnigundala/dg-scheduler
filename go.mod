module github.com/donnigundala/dg-scheduler

go 1.25.0

require (
	github.com/donnigundala/dg-core v0.0.0-00010101000000-000000000000
	github.com/donnigundala/dg-queue v1.3.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/donnigundala/dg-core => ../dg-core
