# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-12-27

### Added
- Initial stable release of the `dg-scheduler` plugin.
- **Cron-Based Scheduling**: Robust job scheduling using `robfig/cron/v3`.
- **Queue Integration**: Seamless integration with `dg-queue` for job dispatch.
- **Schedule Management**: Add, remove, and count scheduled jobs.
- **Container Integration**: Auto-registration with Injectable pattern.
- **Observability**: OpenTelemetry metrics for scheduled job execution.

### Features
- Flexible cron expression support
- Automatic job dispatching to queue
- Thread-safe schedule management
- Graceful shutdown support
- Production-ready with comprehensive test coverage

### Documentation
- Complete README with examples
- Cron expression guide
- Queue integration documentation

---

## Development History

The following versions represent the development journey leading to v1.0.0:

### 2025-12-03
- Extracted from dg-queue as standalone package
- Enhanced schedule management API
- Added observability metrics

### 2025-11-24
- Initial implementation with cron support
- Queue integration for job dispatch
