
# Task Manager API

A simple, modular, and testable RESTful API for managing tasks, written in Go. This project demonstrates clean architecture, in-memory storage, and includes a demo script for API usage.

## Features

- RESTful API for CRUD operations on tasks
- In-memory repository (no external DB required)
- Modular architecture (handler, service, repository, model)
- Graceful shutdown and logging (zap)
- Demo script for API usage
- Comprehensive unit and integration tests

## Architecture

```mermaid
graph TD
    A[Client] -- HTTP/JSON --> B[Handler]
    B -- Calls --> C[Service]
    C -- Uses --> D[Repository]
    D -- Stores --> E[In-Memory Map]
    C -- Uses --> F[ID Generator]
    B -- Logs --> G[Logger-zap]
```

## Data Flow

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Service
    participant Repository
    Client->>Handler: HTTP Request (JSON)
    Handler->>Service: Validate & Process
    Service->>Repository: CRUD Operation
    Repository-->>Service: Result
    Service-->>Handler: Response Data
    Handler-->>Client: HTTP Response (JSON)
```

## Getting Started

### Prerequisites

- Go 1.20+
- `jq` (for demo script)

### Build & Run

```sh
make build
./bin/taskmanager
```

Or run directly:

```sh
go run ./cmd/server
```

### API Endpoints

- `GET    /tasks`         - List all tasks
- `POST   /tasks`         - Create a new task
- `GET    /tasks/{id}`    - Get a task by ID
- `PUT    /tasks/{id}`    - Update a task by ID
- `DELETE /tasks/{id}`    - Delete a task by ID

#### Task JSON Example

```json
{
  "id": "string (auto-generated if omitted)",
  "title": "Task title",
  "description": "Optional description",
  "completed": false
}
```

### Demo Script

Run the provided demo script to see the API in action:

```sh
bash demo_tasks.sh
```

## Testing

Run all tests:

```sh
make test
```

Or:

```sh
go test ./...
```

## Development

- Main entry: `cmd/server/main.go`
- Handlers: `internal/handler/`
- Services: `internal/service/`
- Repository: `internal/repository/`
- Models: `internal/model/`
