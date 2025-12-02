# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-12-03

### Added
- Initial release of dg-scheduler
- Cron-based job scheduling using robfig/cron/v3
- Standalone operation (no dependencies required)
- Optional queue integration via Dispatcher interface
- Configurable error handler for custom error handling
- Structured logging interface for integration with logging systems
- Schedule validation (cron expressions and schedule names)
- Graceful shutdown with context support
- Thread-safe schedule management
- Duplicate schedule prevention
- Helper functions: `ValidateCronExpression`, `ValidateScheduleName`, `ParseCronExpression`
- Comprehensive test coverage
- Examples directory with working code samples
- Full API documentation

### Features
- **Standalone Usage**: Works without any queue system
- **Flexible Integration**: Dispatcher interface allows integration with any queue
- **Production Ready**: Error handling, logging, validation
- **Zero Required Dependencies**: Only robfig/cron/v3 (optional queue integration)
- **Developer Friendly**: Clear API, good documentation, examples

### Breaking Changes
- None (initial release)

### Migration
- Extracted from dg-queue v1.x
- Users of `dg-queue.Manager.Schedule()` should migrate to dg-scheduler
- See README.md for migration guide

[1.0.0]: https://github.com/donnigundala/dg-scheduler/releases/tag/v1.0.0
