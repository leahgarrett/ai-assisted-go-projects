# Todo API

A Go API to read and write task data stored in `data.json`.

## Setup

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run the application:
   ```bash
   go run cmd/todo-api/main.go
   ```

3. Access the API at `http://localhost:8080`.

## Features

- List all tasks
- Get a task by ID
- Add a new task
- Update an existing task
- Delete a task
