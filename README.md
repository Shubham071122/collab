# Collab Backend

Backend API for Collab, a collaborative workspace application. The service is built with Go, Gin, PostgreSQL, JWT authentication, and a layered project structure that keeps HTTP handling, business logic, and database access separate.

## Architecture

The request flow follows a simple layered pattern:

```txt
Route -> Handler -> Service -> Repository -> Database
```

- `Handler` reads HTTP input and writes responses.
- `Service` contains validation and business rules.
- `Repository` is responsible for database queries.

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/shubham071122/collab.git
cd collab
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Configure Environment

Create a `.env` file in the project root:

```env
PORT=8080
JWT_SECRET=your-secret
DATABASE_URL=postgres://user:password@localhost:5432/collab?sslmode=disable
```

`DATABASE_URL` is used by the API and by the migration command.

### 4. Run Migrations

Install `golang-migrate` if you do not have it:

```bash
brew install golang-migrate
```

Apply all migrations:

```bash
migrate -path internal/database/migrations -database "$DATABASE_URL" up
```

Rollback the latest migration:

```bash
migrate -path internal/database/migrations -database "$DATABASE_URL" down 1
```

Create a new migration:

```bash
migrate create -ext sql -dir internal/database/migrations create_example_table
```

This creates one `.up.sql` file for applying changes and one `.down.sql` file for rolling them back.

### 5. Run the API

For normal execution:

```bash
go run ./cmd/api
```

For live reload during development:

```bash
go install github.com/air-verse/air@latest
air
```

If `air` is not available after installation:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Development Notes

- Keep database changes in `internal/database/migrations`.
- Keep route definitions in `internal/routes/routes.go`.
- Add domain-specific logic inside its own package under `internal`.
- Use services for validation and repositories for SQL queries.
