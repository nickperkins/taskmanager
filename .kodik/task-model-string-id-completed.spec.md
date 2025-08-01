# Task: Update Task Model to Support String ID and Completed Boolean

**Detailed Feature Description:**

The current Task Management API uses a UUID (type `uuid.UUID`) for the `id` field and a string-based `status` field (with values like "pending", "in-progress", "completed"). However, the original requirement specifies that a task should have at least an `id` (string), a `title` (string), and a `completed` (bool) status. This update will bring the data model and API into compliance with the original request.

This change will:

- Change the `id` field in the Task model from `uuid.UUID` to `string`.
- Allow the client to provide an `id` when creating a task; if not provided, the server will generate a unique string (e.g., a UUID as a string).
- Replace the `status` field (string/enum) with a `completed` boolean field.
- Update all validation, serialization, and business logic to use the new model.
- Update all endpoints, tests, and documentation to reflect the new model.
- Ensure backward-incompatible changes are clearly documented.

---

## Feature: Task Model Update

**User Story:**
As a developer or client application, I want the Task model to use a string `id` and a boolean `completed` field so that the API matches the original requirements and is easier to integrate with systems that use string IDs and boolean completion status.

### Cucumber Scenarios

**Scenario: Create a new task with client-provided ID**
Given the task management API is running
When I send a POST request to "/tasks" with:

```json
{
  "id": "my-custom-id",
  "title": "My Task",
  "completed": false
}
```

Then the response status should be 201 Created
And the response body should contain the task with id "my-custom-id" and completed false

**Scenario: Create a new task without providing an ID**
Given the task management API is running
When I send a POST request to "/tasks" with:

```json
{
  "title": "Auto-ID Task",
  "completed": false
}
```

Then the response status should be 201 Created
And the response body should contain a task with a generated string id and completed false

**Scenario: Retrieve a task and see completed as boolean**
Given there is a task with id "abc" and completed true
When I send a GET request to "/tasks/abc"
Then the response body should include:

- id: "abc"
- completed: true

**Scenario: Update a task's completed status**
Given there is a task with id "xyz" and completed false
When I send a PUT request to "/tasks/xyz" with:

```json
{
  "completed": true
}
```

Then the response body should show completed: true

**Scenario: List tasks and see completed as boolean**
Given there are multiple tasks in the system
When I send a GET request to "/tasks"
Then each task in the response should have:

- id: string
- title: string
- completed: boolean

---

## Migration and Compatibility

- All code, tests, and documentation referencing `status` or `uuid.UUID` must be updated.
- The API will no longer accept or return a `status` field; instead, it will use `completed` (bool).
- The API will accept string IDs, either provided by the client or generated by the server.

---

## Acceptance Criteria

- The Task model uses `id` (string), `title` (string), and `completed` (bool).
- All endpoints and tests are updated to use the new model.
- The API is fully backward-incompatible with the previous model, and this is documented.
- All validation, error handling, and concurrency guarantees remain in place.
