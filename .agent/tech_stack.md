# Tech Stack & Tooling: GoTicket

Ten dokument definiuje zatwierdzony stack technologiczny i narzędzia używane w projekcie. Zmiana tych narzędzi wymaga konsultacji architektonicznej.

## Core Backend

- **Language**: Go 1.23+
- **API**:
  - REST (Public Gateway): `net/http` (Standard Library, Go 1.22+ routing features).
  - gRPC (Internal): `google.golang.org/grpc`, `protoc` (komunikacja między mikroserwisami).
- **Domain Logic**: Czysty Go, Hexagonal Architecture, brak zewnętrznych zależności w domenie.

## Data Persistence & Messaging

- **Primary DB**: PostgreSQL 16+ (transakcyjność, ACID).
- **DB Drivers**: `pgx/v5` (najszybszy driver Go).
- **Migrations**: `go-migrate` lub `tern`.
- **Query Builder/ORM**: `sqlc` (Type-safe SQL) - preferowane nad GORM.
- **Caching/Locks**: Redis 7+ (Alpine).
- **Messaging**: RabbitMQ (lub NATS JetStream) - do asynchronicznych zdarzeń.

## DevOps & Infrastructure (The "Must Have" for Senior)

- **Containerization**: Docker, Docker Compose (Multi-stage builds).
- **IaC (Infrastructure as Code)**: OpenTofu / Terraform.
- **Cloud**: AWS Free Tier (EC2 t2.micro, S3, RDS db.t3.micro).
- **CI/CD**: GitHub Actions (Lint, Test, Build, Deploy).
- **Observability**:
  - Logs: `slog` (Structured Logging - Go 1.21+).
  - Metrics: Prometheus.
  - Tracing: OpenTelemetry (opcjonalnie w późniejszej fazie).

## Local Development Tools

- `air`: Live reload for Go.
- `golangci-lint`: Strict linter configuration.
- `k6`: Load testing framework.
- `Postman`/`Bruno`: API Testing.

## Middleware & Security Patterns

- **Recovery Middleware**: Panic handling to prevent application crashes.
- **Logging Middleware**: Request/response tracking with ResponseWriter wrapper pattern.
- **Rate Limiting**: Protection against path scanning attacks (future: Redis-backed).
- **Audit Logging**: Compliance logging for GDPR, PCI DSS, HIPAA (structured logs with slog).

## Observability Stack

- **Logs**: `slog` (Structured Logging - Go 1.21+) with JSON output.
- **Metrics**: Prometheus (future).
- **Tracing**: OpenTelemetry (future, optional in later phase).
- **Monitoring**: Grafana dashboards (future).
