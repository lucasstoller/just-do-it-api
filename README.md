# Just Do It API

![Just Do It](./logo.jpeg)

**Just Do It** is a simple todo app created by [Lucas Stoller](www.linkedin.com/in/lucasstoller) using Go. His idea was to create a simple todo app were you can focus on what import the most, do your tasks.

## Getting Started

### Prerequisites

- Go 1.x
- PostgreSQL
- Insomnia (for API testing)

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up PostgreSQL database
4. Run the server:
   ```bash
   go run main.go
   ```

## Testing

The project includes comprehensive test coverage for all handlers. To run the tests:

```bash
go test ./... -v
```

The test suite includes:

- Unit tests for all CRUD operations
- Input validation tests
- Error handling scenarios
- Mock database implementation for reliable testing

## API Documentation

### Authentication

All endpoints require Bearer token authentication:

```
Authorization: Bearer <your-token>
```

### Endpoints

#### Tasks

##### Get All Tasks

- **GET** `/v1/tasks`
- Returns all tasks for the authenticated user

##### Create Task

- **POST** `/v1/tasks`
- Request Body:
  ```json
  {
    "title": "Example Task",
    "description": "Task description",
    "deadline": "2025-01-27T10:00:00Z"
  }
  ```

##### Update Task

- **PUT** `/v1/tasks/:id`
- Request Body:
  ```json
  {
    "title": "Updated Task",
    "description": "Updated description",
    "deadline": "2025-01-27T11:00:00Z"
  }
  ```

##### Delete Task

- **DELETE** `/v1/tasks/:id`

##### Toggle Task Completion

- **PATCH** `/v1/tasks/:id/toggle`

#### Task Filters

##### Get Today's Tasks

- **GET** `/v1/tasks/today`
- Returns tasks due today

##### Get Tasks by Date

- **GET** `/v1/tasks?deadline=YYYY-MM-DD`
- Returns tasks for a specific date

##### Get Backlog Tasks

- **GET** `/v1/tasks/backlog`
- Returns overdue and incomplete tasks

### Insomnia Collection

An Insomnia collection is included in the repository (`insomnia.json`). To use it:

1. Open Insomnia
2. Import the collection from `insomnia.json`
3. Set up the environment variables:
   - `baseUrl`: `http://localhost:8080/v1`
   - `token`: Your authentication token

## License

This app is under AGPL License. Please check [this page](./LICENSE) before using the codebase.
