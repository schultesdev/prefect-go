# Integration Tests

Integration tests require a running Prefect server instance.

## Running Integration Tests

### 1. Start a Prefect Server

```bash
prefect server start
```

This will start a Prefect server on `http://localhost:4200`.

### 2. Run the Integration Tests

```bash
go test -v -tags=integration ./pkg/client/
```

Or using make:

```bash
make test-integration
```

### 3. Custom Server URL

If your Prefect server is running on a different URL, set the `PREFECT_API_URL` environment variable:

```bash
export PREFECT_API_URL="http://your-server:4200/api"
go test -v -tags=integration ./pkg/client/
```

## Test Coverage

The integration tests cover:

- **Flow Lifecycle**: Create, get, update, list, and delete flows
- **Flow Run Lifecycle**: Create, get, set state, and delete flow runs
- **Deployment Lifecycle**: Create, get, pause, resume, and delete deployments
- **Pagination**: Manual pagination and iterator-based pagination
- **Variables**: Create, get, update, and delete variables
- **Admin Endpoints**: Health check and version
- **Flow Run Monitoring**: Wait for flow run completion

## Notes

- Integration tests use the `integration` build tag to separate them from unit tests
- Each test creates and cleans up its own resources
- Tests use unique UUIDs in names to avoid conflicts
- Some tests may take several seconds to complete (e.g., flow run wait tests)

## Troubleshooting

### Connection Refused

If you see "connection refused" errors, make sure:
1. Prefect server is running (`prefect server start`)
2. The server is accessible at the configured URL
3. No firewall is blocking the connection

### Test Timeout

If tests timeout:
1. Increase the test timeout: `go test -timeout 5m -tags=integration ./pkg/client/`
2. Check server logs for errors
3. Ensure the server is not under heavy load
