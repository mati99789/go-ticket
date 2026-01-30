# Project: "GoTicket" - High-Scale Event Ticketing System

## Objective

Create an Enterprise-grade ticket reservation system capable of handling "The Taylor Swift Problem" traffic (high load, massive race conditions).

## Architecture (High Level)

The system follows **Hexagonal Architecture (Ports & Adapters)** with **Domain-Driven Design (DDD)**.
It starts as a Modular Monolith and will evolve into Microservices.

### Data Flow

1.  **Client** sends request with JWT token to **API Gateway**.
2.  **Auth Middleware** validates JWT, extracts user role.
3.  **API Gateway** (GraphQL/REST) routes to specific modules based on permissions.
4.  **Booking Core** handles reservation logic (Atomic Transactions).
5.  **Redis** acts as a Distributed Lock to prevent overselling.
6.  **Postgres** persists the "hard" state (ACID).
7.  **Kafka/RabbitMQ** publishes `BookingConfirmed` events.
8.  **Notification Service** consumes events to send emails.

## Tech Stack (Senior Level)

### 1. Backend & Logic

- **Language**: Go 1.24+
- **API Styles**:
  - **REST**: `net/http` (Standard Library) for CRUD.
  - **GraphQL**: `99designs/gqlgen` for flexible data fetching.
- **Communication**: gRPC (internal service-to-service).
- **Documentation**: Swagger/OpenAPI (`swaggo`).
- **Architecture**: Clean Architecture / Hexagonal.

### 2. Authentication & Authorization (Security)

- **Authentication**:
  - **JWT (JSON Web Tokens)**: Stateless auth for API requests.
  - **OAuth2 / OIDC**: Social login (Google, GitHub) for users.
  - **Refresh Tokens**: Stored in Redis for session management.
- **Authorization**:
  - **RBAC (Role-Based Access Control)**: Roles: `admin`, `organizer`, `user`.
  - **Middleware**: JWT validation on protected routes.
- **Password Security**: bcrypt hashing (cost factor 12+).
- **Libraries**: `golang-jwt/jwt`, `coreos/go-oidc`.

### 3. Data & Persistence

- **Database**: PostgreSQL 16+.
- **Driver**: `pgx/v5` (performance, binary protocol).
- **Migrations**: `golang-migrate` or native SQL.
- **Caching & Locking**: **Redis** (Critical for Race Conditions).

### 3. Asynchronous Messaging (Event Driven)

- **Task Queue (RabbitMQ)**: Best for "Fire & Forget" jobs like Sending Emails, PDF Generation.
- **Event Stream (Kafka)**: Best for Domain Events (`BookingCreated`, `PaymentProcessed`) where persistence and strict ordering matter (Analytics, Audit Log).
- **Pattern**: Outbox Pattern (guaranteed delivery).

### 4. DevOps & Infrastructure

- **Containerization**: **Docker** (Multi-stage builds, Distroless images).
- **Orchestration**: **Kubernetes (K8s)**.
  - Local: Kind / Minikube.
  - Prod: EKS/GKE (simulated via Terraform).
- **IaC (Infrastructure as Code)**: **Terraform** / OpenTofu.
- **CI/CD**: **GitHub Actions** (Lint, Test, Build, Push).

### 5. Quality Assurance

- **Unit Tests**: Standard `testing` package + Table Driven Tests.
- **Integration Tests**: **Testcontainers** (Real DB in Docker).
- **Linting**: `golangci-lint` (strict rules).

## Evolution Roadmap

1.  **Phase 1**: Solid Monolith (REST, Postgres, Docker, Testcontainers).
2.  **Phase 2**: Advanced Features (GraphQL, Redis, Async).
3.  **Phase 3**: DevOps (K8s, Terraform, CI/CD).
4.  **Phase 4**: Microservices Split.
