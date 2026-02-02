# OpenAPI Specification for Prefect Server API

This directory contains the OpenAPI specification for the Prefect v3 Server REST API.

## How to obtain the OpenAPI spec

The OpenAPI specification is available from a running Prefect server instance.

### Option 1: From a local Prefect server

1. Start a Prefect server:
```bash
prefect server start
```

2. Download the OpenAPI specification:
```bash
curl http://localhost:4200/openapi.json -o api/openapi.json
```

### Option 2: From Prefect Cloud

```bash
curl https://api.prefect.cloud/api/openapi.json -o api/openapi.json
```

## Generating Go code

After obtaining the OpenAPI specification, run:

```bash
make generate
```

This will generate Go types and client code from the OpenAPI spec.

## Note

The `openapi.json` file is not included in the repository by default. You need to obtain it from a running Prefect server instance as described above.

For development purposes, a minimal placeholder spec is provided to allow the project to compile without a full Prefect server setup.
