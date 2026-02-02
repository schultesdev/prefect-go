# Prefect Go API Client

[![Go Reference](https://pkg.go.dev/badge/github.com/gregor/prefect-go.svg)](https://pkg.go.dev/github.com/gregor/prefect-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/gregor/prefect-go)](https://goreportcard.com/report/github.com/gregor/prefect-go)

Ein vollständiger Go-Client für die Prefect v3 Server REST API mit Unterstützung für alle Endpoints, automatischen Retries und Pagination.

## Features

- 🚀 Vollständige API-Abdeckung aller Prefect Server v3 Endpoints
- 🔄 Automatische Retry-Logik mit exponential backoff
- 📄 Pagination-Support mit Iterator-Pattern
- 🔐 Authentifizierung für Prefect Cloud
- ✅ Type-safe API mit generierten Go-Structs
- 🧪 Vollständig getestet
- 📝 Umfangreiche Dokumentation und Beispiele

## Installation

```bash
go get github.com/gregor/prefect-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/gregor/prefect-go/pkg/client"
    "github.com/gregor/prefect-go/pkg/models"
)

func main() {
    // Client erstellen
    c := client.NewClient(
        client.WithBaseURL("http://localhost:4200/api"),
    )

    ctx := context.Background()

    // Flow erstellen
    flow, err := c.Flows.Create(ctx, &models.FlowCreate{
        Name: "my-flow",
        Tags: []string{"example", "go-client"},
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Flow erstellt: %s (ID: %s)\n", flow.Name, flow.ID)

    // Flow-Run starten
    run, err := c.FlowRuns.Create(ctx, &models.FlowRunCreate{
        FlowID: flow.ID,
        Name:   "my-first-run",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Flow-Run gestartet: %s (ID: %s)\n", run.Name, run.ID)
}
```

## Authentifizierung

### Lokaler Prefect Server

```go
client := client.NewClient(
    client.WithBaseURL("http://localhost:4200/api"),
)
```

### Prefect Cloud

```go
client := client.NewClient(
    client.WithBaseURL("https://api.prefect.cloud/api/accounts/{account_id}/workspaces/{workspace_id}"),
    client.WithAPIKey("your-api-key"),
)
```

## Verwendungsbeispiele

### Flows verwalten

```go
// Flow erstellen
flow, err := client.Flows.Create(ctx, &models.FlowCreate{
    Name: "my-flow",
    Tags: []string{"production"},
})

// Flow abrufen
flow, err := client.Flows.Get(ctx, flowID)

// Flows auflisten mit Pagination
iter := client.Flows.ListAll(ctx, &models.FlowFilter{
    Tags: []string{"production"},
})
for iter.Next(ctx) {
    flow := iter.Value()
    fmt.Printf("Flow: %s\n", flow.Name)
}
if err := iter.Err(); err != nil {
    log.Fatal(err)
}

// Flow löschen
err := client.Flows.Delete(ctx, flowID)
```

### Flow-Runs verwalten

```go
// Flow-Run erstellen und starten
run, err := client.FlowRuns.Create(ctx, &models.FlowRunCreate{
    FlowID:     flowID,
    Name:       "scheduled-run",
    Parameters: map[string]interface{}{
        "param1": "value1",
        "param2": 42,
    },
})

// Auf Completion warten
finalRun, err := client.FlowRuns.Wait(ctx, run.ID, 5*time.Second)
fmt.Printf("Flow-Run Status: %s\n", finalRun.StateType)
```

### Deployments erstellen

```go
deployment, err := client.Deployments.Create(ctx, &models.DeploymentCreate{
    Name:         "my-deployment",
    FlowID:       flowID,
    WorkPoolName: "default-pool",
    Schedules: []models.DeploymentScheduleCreate{
        {
            Schedule: models.IntervalSchedule{
                Interval: 3600, // Jede Stunde
            },
            Active: true,
        },
    },
})
```

## Konfiguration

### Client-Optionen

```go
client := client.NewClient(
    client.WithBaseURL("http://localhost:4200/api"),
    client.WithAPIKey("your-api-key"),
    client.WithTimeout(30 * time.Second),
    client.WithRetry(retry.Config{
        MaxAttempts:   5,
        InitialDelay:  time.Second,
        MaxDelay:      30 * time.Second,
        BackoffFactor: 2.0,
    }),
)
```

### Retry-Konfiguration

Der Client retried automatisch bei transienten Fehlern (5xx, 429):

```go
import "github.com/gregor/prefect-go/pkg/retry"

retryConfig := retry.Config{
    MaxAttempts:     5,
    InitialDelay:    time.Second,
    MaxDelay:        30 * time.Second,
    BackoffFactor:   2.0,
    RetryableErrors: []int{429, 500, 502, 503, 504},
}
```

## Verfügbare Services

Der Client bietet folgende Services:

- **Flows**: Flow-Management
- **FlowRuns**: Flow-Run-Operationen
- **Deployments**: Deployment-Verwaltung
- **TaskRuns**: Task-Run-Management
- **WorkPools**: Work-Pool-Operationen
- **WorkQueues**: Work-Queue-Verwaltung
- **Artifacts**: Artifact-Handling
- **Blocks**: Block-Dokumente und -Typen
- **Events**: Event-Management
- **Automations**: Automation-Konfiguration
- **Variables**: Variable-Verwaltung
- **ConcurrencyLimits**: Concurrency-Limit-Management
- **Logs**: Log-Abfragen
- **Admin**: Administrative Endpoints

## Entwicklung

### Code generieren

```bash
make generate
```

### Tests ausführen

```bash
make test
```

### Integration-Tests

Integration-Tests erfordern einen laufenden Prefect-Server:

```bash
# Prefect Server starten
prefect server start

# Tests ausführen
go test -tags=integration ./...
```

## Beispiele

Weitere Beispiele finden Sie im [examples/](examples/) Verzeichnis:

- [basic](examples/basic/main.go): Grundlegende Verwendung
- [deployment](examples/deployment/main.go): Deployment mit Schedule
- [monitoring](examples/monitoring/main.go): Flow-Run-Monitoring
- [pagination](examples/pagination/main.go): Pagination-Beispiel

## Dokumentation

Vollständige API-Dokumentation: https://pkg.go.dev/github.com/gregor/prefect-go

## Lizenz

MIT License - siehe [LICENSE](LICENSE) für Details.

## Beiträge

Beiträge sind willkommen! Bitte öffnen Sie ein Issue oder erstellen Sie einen Pull Request.
