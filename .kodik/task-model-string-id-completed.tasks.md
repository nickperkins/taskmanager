# Task List: Task Model with String ID and Completed Boolean

---

### Feature: Data Model Update

- [x] **Task 1.1: Refactor `Task` Struct to Use String ID and Completed Boolean**
  - **Status:** Complete
  - **Context:**
    - Update `/internal/model/task.go` to change `ID` from `uuid.UUID` to `string` and `Status` to `Completed` (bool).
    - Remove all references to `TaskStatus` and `status` field.
    - Add/keep `Description`, `CreatedAt`, `UpdatedAt` fields as specified.
  - **Dependencies:** None
  - **Validation:**
    - [ ] `Task` struct uses `ID string`, `Completed bool`.
    - [ ] No references to `TaskStatus` or `status` remain.
    - [ ] All struct tags and field names match the design.

- [x] **Task 1.2: Update Task Validation Logic**
  - **Status:** Complete
  - **Context:**
    - Validation rules: `ID` non-empty, unique, alphanumeric/dash, max 24 chars; `Title` required, 1-200 chars; `Description` optional, max 1000 chars; `Completed` required (default false).
    - Update `/internal/model/task.go` validation methods.
  - **Dependencies:** Task 1.1
  - **Validation:**
    - [ ] Validation enforces all new rules.
    - [ ] Invalid data returns clear error messages.
    - [ ] Tests for all edge cases exist in `/internal/model/task_test.go`.

- [ ] **Task 1.3: Update All References to Task Model**
  - **Status:** In Progress
  - **Context:**
    - Update all usages of `Task` struct in `/internal/handler/`, `/internal/service/`, `/internal/repository/` to use new fields and types.
    - Remove all code referencing `uuid.UUID` and `TaskStatus`.
  - **Dependencies:** Task 1.1
  - **Validation:**
    - [ ] All code compiles with new model.
    - [ ] No references to old types remain.

---

### Feature: API Endpoints

- [ ] **Task 2.1: Update POST /tasks Handler for String ID and Completed**
  - **Status:** Pending
  - **Context:**
    - Accepts `{ id?, title, description?, completed }`.
    - Generates unique string ID if not provided.
    - Returns 201 with created task.
    - Update `/internal/handler/task_handler.go` and related service/repo logic.
  - **Dependencies:** Task 1.1, Task 1.2, Task 1.3
  - **Validation:**
    - [ ] POST /tasks accepts and returns correct fields.
    - [ ] ID is string, completed is boolean.
    - [ ] Returns 201 Created and correct response body.

- [ ] **Task 2.2: Update GET /tasks and GET /tasks/{id} Handlers**
  - **Status:** Pending
  - **Context:**
    - Returns tasks with new model fields: `id`, `title`, `description`, `completed`, `created_at`, `updated_at`.
    - Update `/internal/handler/task_handler.go` and service/repo logic.
  - **Dependencies:** Task 1.1, Task 1.3
  - **Validation:**
    - [ ] GET endpoints return correct fields and types.
    - [ ] 404 returned for missing task.

- [ ] **Task 2.3: Update PUT /tasks/{id} Handler for Partial Update**
  - **Status:** Pending
  - **Context:**
    - Accepts `{ title?, description?, completed? }`.
    - Only updates provided fields; `id`, `createdAt`, `updatedAt` are immutable.
    - Returns 200 with updated task or 404 if not found.
    - Update `/internal/handler/task_handler.go` and service/repo logic.
  - **Dependencies:** Task 1.1, Task 1.2, Task 1.3
  - **Validation:**
    - [ ] PUT allows partial update of allowed fields only.
    - [ ] Immutable fields are not changed.
    - [ ] Returns 200 OK and correct response body.
    - [ ] Returns 404 for missing task.

- [ ] **Task 2.4: Update DELETE /tasks/{id} Handler**
  - **Status:** Pending
  - **Context:**
    - Deletes task by string ID.
    - Returns 204 No Content or 404 if not found.
    - Update `/internal/handler/task_handler.go` and service/repo logic.
  - **Dependencies:** Task 1.1, Task 1.3
  - **Validation:**
    - [ ] DELETE works with string IDs.
    - [ ] Returns 204 for success, 404 for missing task.

---

### Feature: Service and Repository Layer

- [ ] **Task 3.1: Refactor Service and Repository Interfaces to Use String IDs**
  - **Status:** Pending
  - **Context:**
    - All service and repository methods use `string` for task IDs.
    - Update `/internal/service/` and `/internal/repository/`.
    - Remove all `uuid.UUID` usage.
  - **Dependencies:** Task 1.1, Task 1.3
  - **Validation:**
    - [ ] All interfaces and implementations use `string` IDs.
    - [ ] No `uuid.UUID` remains in service/repo code.

- [ ] **Task 3.2: Enforce Unique String IDs in Repository**
  - **Status:** Pending
  - **Context:**
    - Repository must reject duplicate IDs.
    - Update `/internal/repository/inmemory_task_repository.go` and related code.
  - **Dependencies:** Task 3.1
  - **Validation:**
    - [ ] Creating a task with duplicate ID fails with clear error.
    - [ ] All repository tests cover uniqueness.

---

### Feature: Tests and Validation

- [ ] **Task 4.1: Update and Add Unit Tests for Model and Validation**
  - **Status:** Pending
  - **Context:**
    - Update `/internal/model/task_test.go` for new struct and validation rules.
    - Add tests for all edge cases (ID, title, description, completed).
  - **Dependencies:** Task 1.1, Task 1.2
  - **Validation:**
    - [ ] All validation rules are tested.
    - [ ] Tests pass.

- [ ] **Task 4.2: Update and Add Unit Tests for Service and Repository**
  - **Status:** Pending
  - **Context:**
    - Update `/internal/service/task_service_test.go` and `/internal/repository/inmemory_task_repository_test.go` for new model and ID type.
    - Add tests for unique ID enforcement, completed field, and error cases.
  - **Dependencies:** Task 3.1, Task 3.2
  - **Validation:**
    - [ ] All CRUD and validation logic is tested.
    - [ ] Tests pass.

- [ ] **Task 4.3: Update and Add Handler and Integration Tests**
  - **Status:** Pending
  - **Context:**
    - Update `/internal/handler/task_handler_test.go` and `/internal/handler/task_handler_integration_test.go` for new model and endpoints.
    - Add tests for all endpoint scenarios (create, get, update, delete, validation, error).
  - **Dependencies:** Task 2.1, Task 2.2, Task 2.3, Task 2.4
  - **Validation:**
    - [ ] All endpoint and error scenarios are tested.
    - [ ] Tests pass.

---

### Feature: Documentation and Migration

- [ ] **Task 5.1: Update API Documentation and OpenAPI/Swagger**
  - **Status:** Pending
  - **Context:**
    - Update API docs to reflect new model and endpoints.
    - Update or add OpenAPI/Swagger if present.
  - **Dependencies:** Task 1.1, Task 2.1-2.4
  - **Validation:**
    - [ ] API docs and OpenAPI/Swagger match implementation.
    - [ ] All fields and types are correct.

- [ ] **Task 5.2: Document Breaking Change for Clients**
  - **Status:** Pending
  - **Context:**
    - Clearly document that this is a breaking change: string IDs, completed boolean, no status/UUID.
    - Add migration notes for clients.
  - **Dependencies:** Task 5.1
  - **Validation:**
    - [ ] Migration/breaking change notes are present and clear.
    - [ ] All client-facing docs are updated.

---

### Feature: Final Review and Cleanup

- [ ] **Task 6.1: Codebase Audit and Cleanup**
  - **Status:** Pending
  - **Context:**
    - Search for and remove any remaining references to `uuid.UUID`, `TaskStatus`, or `status`.
    - Ensure all code, tests, and docs are consistent with new model.
  - **Dependencies:** All previous tasks
  - **Validation:**
    - [ ] No legacy references remain.
    - [ ] All code, tests, and docs are consistent.

- [ ] **Task 6.2: Final End-to-End Validation**
  - **Status:** Pending
  - **Context:**
    - Run all unit, integration, and concurrency tests.
    - Confirm all endpoints and business logic work as intended.
  - **Dependencies:** All previous tasks
  - **Validation:**
    - [ ] All tests pass with no failures.
    - [ ] API endpoints behave as described in the spec and design.
    - [ ] Concurrency and thread-safety are confirmed.
