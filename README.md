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
4. Run migrations and start the server:

   ```bash
   # Run migrations and start server
   go run main.go

   # Or, to reset database and rerun all migrations
   go run main.go -reset
   ```

### Database Migrations

The project uses `golang-migrate` for database migrations. Migration files are located in the `migrations` directory:

- `000001_create_users_table.up.sql`: Creates users table
- `000001_create_users_table.down.sql`: Drops users table
- `000002_add_user_id_to_tasks.up.sql`: Adds user_id to tasks table
- `000002_add_user_id_to_tasks.down.sql`: Removes user_id from tasks table

Migrations are automatically run when starting the server. Use the `-reset` flag to drop all tables and rerun migrations:

```bash
go run main.go -reset
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

#### Register

- **POST** `/api/auth/register`
- Request Body:
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- Response:
  ```json
  {
    "token": "jwt-token",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "created_at": "2025-01-27T05:00:00Z",
      "updated_at": "2025-01-27T05:00:00Z"
    }
  }
  ```

#### Login

- **POST** `/api/auth/login`
- Request Body:
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- Response:
  ```json
  {
    "token": "jwt-token",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "created_at": "2025-01-27T05:00:00Z",
      "updated_at": "2025-01-27T05:00:00Z"
    }
  }
  ```

#### Protected Endpoints

All task endpoints require Bearer token authentication:

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
   - `baseUrl`: `http://localhost:8080`
   - `token`: JWT token from login response (automatically set after login)

## License

This app is under AGPL License. Please check [this page](./LICENSE) before using the codebase.
