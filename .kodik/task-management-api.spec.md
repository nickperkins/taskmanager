# Task: Task Management API in Go

## Project Overview

A simple, in-memory Task Management API built in Go that provides RESTful HTTP endpoints for basic CRUD operations on tasks. The system emphasizes thread-safety for concurrent access and follows SOLID principles for maintainable, extensible code architecture.

## System Requirements

- HTTP/2 server with RESTful API endpoints
- Thread-safe in-memory data store for concurrent operations
- Basic CRUD functionality for task management
- SOLID principles implementation
- JSON request/response format
- Proper error handling and HTTP status codes

---

### Feature: Task Creation

**Description**: Users can create new tasks through the API with essential task information including title, description, and initial status.

**User Story**: As a client application, I want to create new tasks via HTTP POST request so that users can add tasks to their task list.

**Scenario: Create a new task successfully**

```gherkin
Given the task management API is running
When I send a POST request to "/tasks" with valid task data:

  {
    "title": "Complete project documentation",
    "description": "Write comprehensive documentation for the task management API",
    "status": "pending"
  }

Then the response status should be 201 Created
And the response body should contain the created task with a unique ID
And the response should include:
```

- Generated unique task ID
- Title: "Complete project documentation"
- Description: "Write comprehensive documentation for the task management API"
- Status: "pending"
- Created timestamp
- Updated timestamp

**Scenario: Handle invalid task creation data**

```gherkin
Given the task management API is running
When I send a POST request to "/tasks" with invalid data:
  {
    "title": "",
    "description": "Task with empty title"
  }

Then the response status should be 400 Bad Request
And the response body should contain validation error details

```

---

### Feature: Task Retrieval

**Description**: The API provides endpoints to retrieve all tasks or a specific task by ID, supporting both list and detail views.

**User Story**: As a client application, I want to retrieve tasks from the API so that users can view their task list and individual task details.

**Scenario: Retrieve all tasks**

```gherkin
Given the task management API is running
And there are 3 tasks in the system:
  | id | title | status |
  | 1  | Task 1 | pending |
  | 2  | Task 2 | completed |
  | 3  | Task 3 | in-progress |
When I send a GET request to "/tasks"
Then the response status should be 200 OK
And the response body should contain an array of 3 tasks
And each task should include id, title, description, status, created_at, and updated_at fields
```

**Scenario: Retrieve empty task list**

```gherkin
Given the task management API is running
And there are no tasks in the system
When I send a GET request to "/tasks"
Then the response status should be 200 OK
And the response body should contain an empty array
```

**Scenario: Retrieve a specific task by ID**

```gherkin
Given the task management API is running
And there is a task with ID "123" in the system
When I send a GET request to "/tasks/123"
Then the response status should be 200 OK
And the response body should contain the task details:
  - ID: "123"
  - Title, description, status
  - Created and updated timestamps
```

**Scenario: Handle non-existent task retrieval**

```gherkin
Given the task management API is running
When I send a GET request to "/tasks/999" for a non-existent task
Then the response status should be 404 Not Found
And the response body should contain an appropriate error message
```

---

### Feature: Task Updates

**Description**: Clients can update existing tasks by sending PUT requests with modified task data, including status changes and content updates.

**User Story**: As a client application, I want to update existing tasks so that users can modify task details and mark tasks as complete.

**Scenario: Update an existing task successfully**

```gherkin
Given the task management API is running
And there is a task with ID "456" with status "pending"
When I send a PUT request to "/tasks/456" with updated data:
  ```json
  {
    "title": "Updated task title",
    "description": "Updated description",
    "status": "completed"
  }
  ```

Then the response status should be 200 OK
And the response body should contain the updated task
And the task's updated_at timestamp should be more recent than created_at
And the task status should be "completed"

```

**Scenario: Handle partial task updates**
```gherkin
Given the task management API is running
And there is a task with ID "789"
When I send a PUT request to "/tasks/789" with partial data:
  ```json
  {
    "status": "in-progress"
  }
  ```

Then the response status should be 200 OK
And only the status field should be updated
And other fields should remain unchanged
And the updated_at timestamp should be updated

```

**Scenario: Handle update of non-existent task**
```gherkin
Given the task management API is running
When I send a PUT request to "/tasks/999" for a non-existent task
Then the response status should be 404 Not Found
And the response body should contain an appropriate error message
```

---

### Feature: Task Deletion

**Description**: The API allows clients to permanently remove tasks from the system using DELETE requests.

**User Story**: As a client application, I want to delete tasks so that users can remove unwanted or completed tasks from their list.

**Scenario: Delete an existing task successfully**

```gherkin
Given the task management API is running
And there is a task with ID "321" in the system
When I send a DELETE request to "/tasks/321"
Then the response status should be 204 No Content
And the response body should be empty
And the task should no longer exist in the system
```

**Scenario: Handle deletion of non-existent task**

```gherkin
Given the task management API is running
When I send a DELETE request to "/tasks/999" for a non-existent task
Then the response status should be 404 Not Found
And the response body should contain an appropriate error message
```

---

### Feature: Thread-Safe Concurrent Operations

**Description**: The in-memory store must handle concurrent read and write operations safely without data corruption or race conditions.

**User Story**: As a system administrator, I want the API to handle multiple concurrent requests safely so that data integrity is maintained under load.

**Scenario: Handle concurrent read operations**

```gherkin
Given the task management API is running
And there are multiple tasks in the system
When 10 concurrent GET requests are sent to "/tasks"
Then all requests should return 200 OK
And all responses should contain consistent data
And no race conditions should occur
```

**Scenario: Handle concurrent write operations**

```gherkin
Given the task management API is running
When multiple concurrent POST requests are sent to create tasks
Then each request should receive a unique task ID
And all tasks should be stored successfully
And no data corruption should occur
```

**Scenario: Handle mixed concurrent operations**

```gherkin
Given the task management API is running
And there are existing tasks in the system
When concurrent read, write, update, and delete operations are performed
Then all operations should complete successfully
And data consistency should be maintained
And no deadlocks should occur
```

---

### Feature: Data Model Structure

**Description**: Tasks are represented as structured data with specific fields and validation requirements.

**User Story**: As a developer, I want a well-defined task data model so that the API maintains consistent data structure and validation.

**Scenario: Task data model validation**

```gherkin
Given a task is being created or updated
When the task data is processed
Then the task should have the following structure:
  - id: string (UUID format, auto-generated)
  - title: string (required, 1-200 characters)
  - description: string (optional, max 1000 characters)
  - status: string (enum: "pending", "in-progress", "completed")
  - created_at: timestamp (ISO 8601 format)
  - updated_at: timestamp (ISO 8601 format)
```

---

### Feature: Error Handling and HTTP Status Codes

**Description**: The API provides appropriate HTTP status codes and error messages for different scenarios.

**User Story**: As a client application developer, I want clear error responses so that I can handle different error conditions appropriately.

**Scenario: Return appropriate HTTP status codes**

```gherkin
Given the task management API is running
When various operations are performed
Then the API should return appropriate status codes:
  - 200 OK for successful GET and PUT operations
  - 201 Created for successful POST operations
  - 204 No Content for successful DELETE operations
  - 400 Bad Request for invalid request data
  - 404 Not Found for non-existent resources
  - 500 Internal Server Error for server-side errors
```

**Scenario: Provide detailed error responses**

```gherkin
Given the task management API encounters an error
When an error response is returned
Then the response should include:
  - Appropriate HTTP status code
  - Error message describing the issue
  - Error code for programmatic handling
  - Timestamp of the error
```

---

## Architecture Considerations

### SOLID Principles Implementation

**Single Responsibility Principle (SRP)**

- Separate handlers for each HTTP endpoint
- Dedicated service layer for business logic
- Isolated repository layer for data operations
- Distinct validation and serialization components

**Open/Closed Principle (OCP)**

- Interface-based design for extensibility
- Plugin architecture for adding new storage backends
- Configurable validation rules
- Extensible middleware chain

**Liskov Substitution Principle (LSP)**

- Storage interface that can be implemented by different backends
- Consistent behavior across all implementations
- Proper error handling contracts

**Interface Segregation Principle (ISP)**

- Small, focused interfaces for different concerns
- Separate read and write operations where appropriate
- Client-specific interface definitions

**Dependency Inversion Principle (DIP)**

- Depend on abstractions, not concretions
- Dependency injection for testability
- Configuration-driven component wiring

### Concurrency Safety Requirements

**Thread-Safe Data Store**

- Use Go's sync.RWMutex for concurrent access control
- Protect all read and write operations to the data store
- Ensure atomic operations for complex state changes

**Resource Locking Strategy**

- Reader-writer locks for optimal read performance
- Minimize lock duration to prevent bottlenecks
- Avoid nested locks to prevent deadlocks

**Data Consistency**

- Ensure all CRUD operations maintain data integrity
- Implement proper error handling for concurrent modifications
- Use appropriate Go concurrency patterns

### Implementation Guidelines

**Project Structure (Idiomatic Go, per community standards)**

```
/cmd/server/          - Application entry point (main.go)
/internal/handler/    - HTTP handlers (private to app)
/internal/service/    - Business logic (private)
/internal/repository/ - Data access layer (private)
/internal/model/      - Data models (private)
/internal/config/     - Configuration (private)
/pkg/                 - Public packages (only if truly reusable)
/docs/                - Project documentation (optional)
```

- Do not use `/src/` or a separate `/tests` folder. Place all tests alongside the code they test.
- Use `/test/` only for external test data or integration helpers if needed.
- Keep the structure simple and focused on clarity and idiomatic Go.

**Key Components**

- HTTP/2 router and middleware setup (use standard library or community-accepted router)
- Request/response serialization (encoding/json)
- Input validation (hand-rolled for this challenge, as is idiomatic)
- Structured logging (use go.uber.org/zap)
- Graceful shutdown handling

**Testing Strategy**

- Use the standard `testing` package for all tests
- Use `github.com/stretchr/testify` for assertions and mocks (community standard)
- Place all test files (`*_test.go`) in the same directory as the code under test
- Unit tests for all components
- Integration tests for API endpoints (if time allows)
- Error scenario coverage
- Concurrency/race condition tests for the in-memory store

**Technology Stack**

- Go standard library HTTP/2 server
- encoding/json for JSON
- github.com/google/uuid for UUIDs
- go.uber.org/zap for logging
- github.com/stretchr/testify for test assertions and mocks

**Other Guidelines**

- Use `gofmt` and `staticcheck` for formatting and linting
- Use Go doc comments for all exported types and functions
- Keep the codebase simple, clean, and idiomatic—showcase your ability to write production-quality Go
- Avoid over-engineering; focus on clarity, maintainability, and idiomatic Go patterns
