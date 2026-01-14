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
            Mod_Booking[Module: Booking (Todo)]
        end

        DB[(PostgreSQL)]
    end

    User -->|HTTP| HTTP
    HTTP --> Mod_Event

    Mod_Event --> DB
```

### Components

1.  **Domain Layer (`internal/domain`)**: Pure Go structs and logic. No external dependencies.
2.  **Repository Layer (`internal/postgres`)**: Database implementation using `sqlc` for type-safe SQL.
3.  **API Layer (`internal/api`)**: HTTP Handlers, Routing (`net/http`), and DTOs.
4.  **App Wiring (`cmd/app`)**: Dependency Injection and Graceful Shutdown logic.

## 🛠️ Tech Stack & Tooling

### Core Backend

- **Language**: Go 1.25+
- **Router**: `net/http` (Standard Library) + Go 1.22 routing features
- **Database**: PostgreSQL 16
- **Driver**: `pgx/v5` + `pgxpool` (High performance connection pooling)
- **SQL Generator**: `sqlc` (Compiles SQL to Go)
- **Logging**: `log/slog` (Structured JSON logging)

### DevOps & Infrastructure

- **Docker & Compose**: Local development environment
- **Migrations**: `golang-migrate`
- **Linter**: `golangci-lint` (Strict configuration)
- **(Planned) IaC**: OpenTofu / Terraform
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

## 🧪 Testing the API

You can use `curl` or Postman to interact with the API.

| Method   | Endpoint       | Description                   |
| :------- | :------------- | :---------------------------- |
| `POST`   | `/events`      | Create a new event            |
| `GET`    | `/events/{id}` | Get event details             |
| `PUT`    | `/events/{id}` | Update event name or schedule |
| `DELETE` | `/events/{id}` | Delete an event               |
| `GET`    | `/events`      | List all events               |

**Example Request:**

```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Metallica Live",
    "price": 25000,
    "startAt": "2024-12-01T20:00:00Z",
    "endAt": "2024-12-01T23:00:00Z"
  }'
```

## ✅ Definition of Done

To pass Code Review, every PR must meet the strict standards defined in `.agent/definition_of_done.md`:

- [ ] Linter clean (`golangci-lint run`)
- [ ] Zero magic numbers
- [ ] Domain logic unit tested
- [ ] Clean Git history (Conventional Commits)

## 🗺️ Roadmap & Progress

- [x] **Phase 0: Setup** (Go, Linter, Project Structure)
- [x] **Phase 1: Domain Core** (Entities, Validation, Unit Tests)
- [x] **Phase 2: Persistence** (Postgres, pgxpool, sqlc, Migrations)
- [x] **Phase 3: HTTP API** (REST Handlers, JSON, Graceful Shutdown)
- [ ] **Phase 4: Optimization** (Integration Tests, Caching)
- [ ] **Phase 5: Booking Logic** (Transactions, Distributed Locking)
- [ ] **Phase 6: Cloud** (AWS Deployment)
