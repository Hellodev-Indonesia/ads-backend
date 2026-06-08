# Ads Backend

This project is a backend API written in Go, utilizing Gin, GORM, MySQL, Redis, and Centrifugo. It follows a Domain-Driven Module-Based architecture, ready to scale into an ERP system.

## Prerequisites

Before you begin, ensure you have the following installed on your machine:

- **[Go](https://golang.org/doc/install)** (1.20+ recommended)
- **[Docker & Docker Compose](https://docs.docker.com/get-docker/)** (for running MySQL, Redis, and Centrifugo locally)
- **[golang-migrate](https://github.com/golang-migrate/migrate)** (for running database migrations)
- **[Make](https://www.gnu.org/software/make/)** (for running Makefile commands)

## Installation Guide

Follow these steps to set up the project locally:

1. **Clone the Repository**
   ```bash
   git clone <repository-url>
   cd ads_backend
   ```

2. **Set Up Environment Variables**
   Copy the example environment file and update it with your local credentials if necessary.
   ```bash
   cp .env.example .env
   ```

3. **Install Go Dependencies**
   Download all required Go modules.
   ```bash
   go mod download
   ```

4. **Start the Infrastructure**
   Spin up MySQL, Redis, and Centrifugo using Docker Compose.
   ```bash
   docker-compose up -d
   ```

5. **Run Database Migrations**
   Apply the database schema using `golang-migrate`. The Makefile command will read your `.env` variables.
   ```bash
   make migrate-up
   ```

6. **Run Database Seeders (Optional)**
   Populate your database with initial required data.
   ```bash
   make seed
   ```

7. **Generate API Documentation (Optional)**
   Generate Swagger docs.
   ```bash
   make docs
   ```

8. **Start the API Server**
   Run the development server.
   ```bash
   make dev
   ```

---

## Makefile Usage Guide

The `Makefile` contains many helpful shortcuts for development, testing, and database management. 

To see a list of all available commands directly in your terminal, run:
```bash
make help
```

### Development Commands

| Command | Description |
|---|---|
| `make dev` | Starts the API server using `go run cmd/api/main.go`. |
| `make build` | Compiles the Go code and builds binaries for the API and seeders into the `bin/` directory. |
| `make docs` | Generates the Swagger OpenAPI documentation using `swaggo`. |
| `make mock` | Generates or updates test mocks using `mockery`. |

### Database & Migration Commands

| Command | Description |
|---|---|
| `make migrate-up` | Applies all pending migrations to update the database schema. |
| `make migrate-down` | Rolls back the most recently applied migration by 1 step. |
| `make migrate-fresh` | Drops all tables completely and re-runs all migrations from scratch. **(Use with caution)** |
| `make seed` | Runs the seeder scripts to populate the database with initial/dummy data. |

### Testing Commands

| Command | Description |
|---|---|
| `make test` | Runs all tests across the entire project. |
| `make test-unit` | Runs only the unit tests (skips tests requiring a database/integration). |
| `make test-integration` | Runs only the integration tests. |
| `make test-coverage` | Runs all tests, generates a coverage profile, and opens the HTML coverage report in your browser. |
