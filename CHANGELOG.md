# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of prefect-go
- Complete Prefect Server v3 API client implementation
- Support for all major resources (Flows, FlowRuns, Deployments, etc.)
- Automatic retry logic with exponential backoff
- Pagination support with Iterator pattern
- Structured error types for API errors
- Context-aware operations with cancellation support
- Comprehensive test coverage
- Example programs for common use cases
- Full GoDoc documentation

### Services
- Flows service with CRUD operations
- FlowRuns service with state management and monitoring
- Deployments service with pause/resume support
- TaskRuns service with state transitions
- WorkPools and WorkQueues services
- Variables service for key-value storage
- Logs service for log management
- Admin service for health checks and version info

### Features
- Options pattern for flexible client configuration
- Retry configuration with customizable backoff
- Pagination helpers for large result sets
- Flow run monitoring with Wait() method
- Type-safe API with generated models
- Thread-safe client operations

### Developer Experience
- Integration test suite
- Example programs for learning
- Contributing guidelines
- Comprehensive README
- API documentation

## [1.0.0] - TBD

Initial stable release.

[Unreleased]: https://git.schultes.dev/schultesdev/prefect-go/compare/v1.0.0...HEAD
[1.0.0]: https://git.schultes.dev/schultesdev/prefect-go/releases/tag/v1.0.0
