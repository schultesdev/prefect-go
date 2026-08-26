# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- **BREAKING**: `Variable.Value`, `VariableCreate.Value` and `VariableUpdate.Value`
  are now `json.RawMessage` instead of `string` (`*string` for the update case).
  Prefect 3.x stores a variable's value as any JSON type - string, integer,
  number, boolean, object, array or null - so a `string` field could not encode
  an object or array on write, and failed outright on read with
  `cannot unmarshal object into Go struct field .Alias.value of type string`.

  Migration:

  ```go
  // before
  v := &models.VariableCreate{Name: "cfg", Value: "/data/imports"}
  s := variable.Value

  // after
  value, err := models.VariableValue("/data/imports") // or any Go value
  v := &models.VariableCreate{Name: "cfg", Value: value}
  s, ok := variable.StringValue() // ok == false unless the value is a JSON string
  ```

  A nil `VariableUpdate.Value` is omitted from the request as before, leaving the
  stored value untouched; pass `json.RawMessage("null")` to set it to JSON null.

### Added
- `models.VariableValue(any) (json.RawMessage, error)` to encode an arbitrary Go
  value for a variable.
- `Variable.StringValue() (string, bool)` to read a variable whose value is a
  JSON string, reporting `false` for every other JSON type.
- Tests for the Variables service (`pkg/client/variables_test.go`), covering all
  JSON value types on read, object and string values on create, and the
  nil/explicit-null distinction on update.
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
