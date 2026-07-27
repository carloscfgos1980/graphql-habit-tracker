# GraphQL Habit Tracker API

A Go GraphQL API for tracking daily habits, check-ins, and streak-related data.

## Tech Stack

- Go
- gqlgen
- Gin
- SQLite (modernc.org/sqlite)
- JWT authentication

## Features

- User registration and login
- JWT-protected GraphQL endpoint
- CRUD operations for habits
- Habit check-ins (habit logs)
- Query habit logs by habit

## Project Structure

- `cmd/api/main.go`: app entrypoint
- `internal/graph/schema.graphqls`: GraphQL schema
- `internal/graph/resolvers/`: GraphQL resolvers
- `internal/repository/`: database access layer
- `internal/middleware/`: auth middleware
- `internal/database/sqlite.go`: SQLite initialization
- `migrations/`: SQL migrations

## Prerequisites

- Go 1.26+
- `goose` for migrations

Install goose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Environment Variables

Create/update `.env`:

```env
DATABASE_PATH=./data/habit.db
PORT=3000
JWT_SECRET=your_jwt_secret_key
```

## Install Dependencies

```bash
go mod tidy
```

## Database Migrations

Run migrations against the same DB file used by `DATABASE_PATH`:

```bash
goose -dir ./migrations sqlite3 ./data/habit.db up
```

Check migration status:

```bash
goose -dir ./migrations sqlite3 ./data/habit.db status
```

Rollback one migration:

```bash
goose -dir ./migrations sqlite3 ./data/habit.db down
```

## Run the API

```bash
go run ./cmd/api/main.go
```

- Playground: `http://localhost:3000/playground`
- GraphQL endpoint: `POST http://localhost:3000/graphql`

## Authentication

- `register` and `login` are public.
- Most other operations require a JWT.
- Send JWT in header:

```http
Authorization: Bearer <token>
```

## cURL Examples (Authenticated)

Login and copy the `token` from the response:

```bash
curl -X POST http://localhost:3000/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"mutation { login(email: \"carlos@example.com\", password: \"StrongPass123!\") { token } }"}'
```

Set your JWT in a shell variable:

```bash
TOKEN="<paste-jwt-token-here>"
```

Create a habit (authenticated):

```bash
curl -X POST http://localhost:3000/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"query":"mutation { createHabit(name: \"Workout\", description: \"30 minutes\") { id name } }"}'
```

List habits (authenticated):

```bash
curl -X POST http://localhost:3000/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"query":"query { habits { id name description } }"}'
```

Check in a habit with GraphQL variables (authenticated):

```bash
curl -X POST http://localhost:3000/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "query":"mutation CheckIn($habitId: ID!, $date: String) { checkInHabit(habitId: $habitId, date: $date) { id completedDate createdAt } }",
    "variables": {"habitId":"<habit-id>", "date":"2026-07-27"}
  }'
```

Get habit logs by habit with GraphQL variables (authenticated):

```bash
curl -X POST http://localhost:3000/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "query":"query HabitLogs($habitId: ID!) { habitLogs(habitId: $habitId) { id completedDate createdAt habit { id name } } }",
    "variables": {"habitId":"<habit-id>"}
  }'
```

## GraphQL Operations

### Register

```graphql
mutation {
  register(name: "Carlos", email: "carlos@example.com", password: "StrongPass123!") {
    token
    user {
      id
      name
      email
    }
  }
}
```

### Login

```graphql
mutation {
  login(email: "carlos@example.com", password: "StrongPass123!") {
    token
    user {
      id
      name
      email
    }
  }
}
```

### Create Habit

```graphql
mutation {
  createHabit(name: "Workout", description: "30 minutes of training") {
    id
    name
    description
    createdAt
  }
}
```

### Check In Habit

`date` is optional and supports RFC3339 or `YYYY-MM-DD`.

```graphql
mutation {
  checkInHabit(habitId: "<habit-id>", date: "2026-07-27") {
    id
    completedDate
    createdAt
  }
}
```

### List Habits

```graphql
query {
  habits {
    id
    name
    description
    currentStreak
    longestStreak
    totalCompletions
  }
}
```

### Habit Logs by Habit

```graphql
query {
  habitLogs(habitId: "<habit-id>") {
    id
    completedDate
    createdAt
    habit {
      id
      name
    }
  }
}
```

## Generate GraphQL Code

If you modify the schema:

```bash
go run github.com/99designs/gqlgen generate
```

## Run Checks

```bash
go test ./...
```

## Notes

- Keep migration target DB path consistent with `DATABASE_PATH`.
- If you run into table-not-found errors, verify migration status and DB file path first.
