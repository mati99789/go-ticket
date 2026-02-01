# GoTicket 🎫

> **"Top of the TOP" Mentorship Project**
> A high-performance, scalable Event Ticketing System built in Go.

## 🎯 About The Project

GoTicket is not just another ToDo app. It is a rigorous engineering training ground designed to master **Advanced Go Backend Development** and **DevOps** practices.

The goal is to build a system capable of handling high-concurrency traffic (e.g., ticket sales for massive events like a Taylor Swift concert), evolving from a **Modular Monolith** to a full **Microservices** architecture.

### 🧠 Core Philosophy (The "No Code" Rule)

- **Zero Magic**: We avoid frameworks like Gin or Fiber. We use the standard library (`net/http`) to understand how HTTP works under the hood.
- **"Why?" over "How?"**: Every line of code must be justified. We choose `pgxpool` over `sql.DB` for specific performance reasons, not just because it's popular.
- **Domain-Driven Design (DDD)**: Business logic is the heart of the system, isolated in the `domain` layer, agnostic of database or transport.
- **Production Standards**: We don't write "tutorial code". We write code ready for production from day one (Graceful Shutdown, Structured Logging, Linter, Docker).

## 🏗️ Architecture (C4 Model)

The project follows **Clean Architecture** principles and implements a **Modular Monolith** structure (Phase 1).

```mermaid
graph TD
    User((User))

    subgraph "GoTicket App (Modular Monolith)"
        HTTP[API Layer / Router]

        subgraph "Modules (Internal Packages)"
            Mod_Event[Module: Event]
            Mod_Auth[Module: Auth (Todo)]
            Mod_Booking[Module: Booking]
        end

        DB[(PostgreSQL)]
    end

    User -->|HTTP| HTTP
    HTTP --> Mod_Event
    HTTP --> Mod_Booking

    Mod_Event --> DB
    Mod_Booking --> DB
```

### Components

1.  **Domain Layer (`internal/domain`)**: Pure Go structs and logic. No external dependencies.
2.  **Repository Layer (`internal/postgres`)**: Database implementation using `sqlc` for type-safe SQL.
3.  **Service Layer (`internal/services`)**: Transaction orchestration and business workflows.
4.  **API Layer (`internal/api`)**: HTTP Handlers, Routing (`net/http`), DTOs, and Error Mapping.
5.  **App Wiring (`cmd/app`)**: Dependency Injection and Graceful Shutdown logic.

## 🛠️ Tech Stack & Tooling

### Core Backend

- **Language**: Go 1.25+
- **Router**: `net/http` (Standard Library) + Go 1.22 routing features
- **Database**: PostgreSQL 16
- **Driver**: `pgx/v5` + `pgxpool` (High performance connection pooling)
- **SQL Generator**: `sqlc` (Compiles SQL to Go)
- **Logging**: `log/slog` (Structured JSON logging)

### Testing

- **Integration Tests**: `testcontainers-go` (Real Postgres in Docker)
- **Assertions**: `testify/assert`
- **Race Detection**: `go test -race`

### DevOps & Infrastructure

- **Docker & Compose**: Local development environment
- **Migrations**: `golang-migrate`
- **Linter**: `golangci-lint` (Strict configuration)
- **(Planned) IaC**: OpenTofu / Terraform
- **(Planned) Observability**: Prometheus + Grafana
- **(Planned) Cloud**: AWS Free Tier

## 🚀 Getting Started

### Prerequisites

- Go 1.22+ to build
- Docker & Docker Compose to run infrastructure

### Quick Start

1.  **Clone the Repository**:

    ```bash
    git clone https://github.com/mati/go-ticket.git
    cd go-ticket
    ```

2.  **Start Infrastructure (Postgres)**:

    ```bash
    docker compose up -d
    ```

3.  **Run Database Migrations**:

    ```bash
    migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/go_ticket?sslmode=disable" up
    ```

4.  **Run the Application**:

    ```bash
    go run cmd/app/main.go
    ```

    _The server will start on port :8080_

5.  **Run Tests**:

    ```bash
    # Unit tests
    go test ./internal/domain/... -v

    # Integration tests (requires Docker)
    go test ./internal/postgres/... -v

    # Race detection
    go test -race ./...
    ```

## 🧪 Testing the API

You can use `curl` or Postman to interact with the API.

### Event Endpoints

| Method   | Endpoint       | Description                   |
| :------- | :------------- | :---------------------------- |
| `POST`   | `/events`      | Create a new event            |
| `GET`    | `/events/{id}` | Get event details             |
| `PUT`    | `/events/{id}` | Update event name or schedule |
| `DELETE` | `/events/{id}` | Delete an event               |
| `GET`    | `/events`      | List all events               |

### Booking Endpoints

| Method | Endpoint                | Description      |
| :----- | :---------------------- | :--------------- |
| `POST` | `/events/{id}/bookings` | Create a booking |

**Example Request:**

```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Metallica Live",
    "price": 25000,
    "startAt": "2024-12-01T20:00:00Z",
    "endAt": "2024-12-01T23:00:00Z",
    "capacity": 1000
  }'
```

## 🎓 What I Learned

### Advanced Go Patterns

- **Domain-Driven Design**: Separating business logic from infrastructure
- **Repository Pattern**: Clean abstraction over database operations
- **DTO Pattern**: API contracts independent of domain models
- **Options Pattern**: Flexible test fixtures
- **Error Mapping**: Layered error handling (DB → Domain → HTTP)

### Concurrency & Race Conditions

- **Atomic Operations**: `sync/atomic` for thread-safe counters
- **WaitGroups**: Coordinating goroutines in tests
- **Database Transactions**: Preventing double-booking with `UPDATE ... WHERE`
- **Race Detection**: Verified with `go test -race`

**Proof:** Concurrent test with 200 goroutines trying to book 100 spots - only 100 succeeded! ✅

### Testing Best Practices

- **Testcontainers**: Real Postgres in tests (not mocks!)
- **Integration Tests**: End-to-end verification with real database
- **Test Fixtures**: Reusable test data with options pattern
- **Isolation**: Each test gets clean database instance

### Production-Ready Code

- **Structured Logging**: `slog` with JSON output
- **Graceful Shutdown**: Proper cleanup on SIGTERM
- **Error Handling**: Domain errors mapped to HTTP status codes
- **Type Safety**: `sqlc` generates type-safe database code

## ✅ Definition of Done

To pass Code Review, every PR must meet the strict standards defined in `.agent/definition_of_done.md`:

- [x] Linter clean (`golangci-lint run`)
- [x] Zero magic numbers
- [x] Domain logic unit tested
- [x] Integration tests passing
- [x] No race conditions (`go test -race`)
- [x] Clean Git history (Conventional Commits)

## 🗺️ Roadmap & Progress

- [x] **Phase 0: Setup** (Go, Linter, Project Structure)
- [x] **Phase 1: Domain Core** (Entities, Validation, Unit Tests)
- [x] **Phase 2: Persistence** (Postgres, pgxpool, sqlc, Migrations)
- [x] **Phase 3: HTTP API** (REST Handlers, JSON, Graceful Shutdown)
- [x] **Phase 3.5: Booking System** (Transactions, Atomic Reservations)
- [x] **Phase 4: Polish** (DTO Layer, Error Mapping, Integration Tests)
- [ ] **Phase 5: DevOps** (Docker, CI/CD, Observability)
- [ ] **Phase 6: Cloud** (AWS Deployment, Terraform)

### Current Focus

**Integration Testing** - Verifying system behavior with real database:

- ✅ Testcontainers setup
- ✅ Repository tests (CRUD, error handling)
- ✅ Race condition tests (200 concurrent bookings)
- 🔄 Service tests (transaction verification)
- 🔄 Handler tests (end-to-end API)

## 📈 Test Coverage

```bash
$ go test ./... -cover

internal/domain         coverage: 95.2%
internal/postgres       coverage: 87.3%
internal/services       coverage: 82.1%
internal/api            coverage: 78.9%
```

## 🤝 Contributing

This is a mentorship project following the **"Zero Gotowego Kodu"** (No Ready Code) rule. The mentor provides guidance, pseudocode, and architecture decisions, but all code is written by the mentee.

## 📄 License

MIT License - See LICENSE file for details

---

**Built with ❤️ and strict mentorship principles**
