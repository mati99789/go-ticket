# Tech Stack & Tooling: GoTicket

This document defines the **approved technology stack** and tooling used in the project. Changing any of these tools requires an architectural consultation and a new ADR.

---

## Core Backend

| Technology | Choice | Reason |
|---|---|---|
| Language | Go 1.23+ | Performance, concurrency model, single binary deployment |
| Public API | `net/http` (Standard Library) | Understand HTTP fundamentals, not framework magic |
| Internal API | gRPC (`google.golang.org/grpc`) | Typed contracts between microservices |
| Domain Logic | Pure Go, no external deps | Hexagonal Architecture — domain must be infrastructure-agnostic |

---

## Data Persistence & Messaging

| Technology | Choice | Reason |
|---|---|---|
| Primary DB | PostgreSQL 16+ | ACID guarantees, mature, battle-tested |
| DB Driver | `pgx/v5` | Fastest Go driver, binary protocol |
| Migrations | `golang-migrate` | Version-controlled schema changes |
| Query Layer | `sqlc` | Type-safe SQL (preferred over GORM — know your queries) |
| Caching/Locks | Redis 7+ (Alpine) | Distributed locking, rate limiting, session cache |
| Event Stream | Kafka (KRaft mode) | Ordered domain events, at-least-once delivery, Outbox Pattern |
| Task Queue | RabbitMQ | Fire-and-forget jobs (email, PDF, notifications) |

---

## DevOps & Infrastructure

| Technology | Choice | Reason |
|---|---|---|
| Containerization | Docker (multi-stage, distroless) | Minimal attack surface, smallest image size |
| Local Orchestration | Docker Compose | Full stack local development |
| IaC | OpenTofu / Terraform | Reproducible infrastructure |
| Cloud Target | AWS Free Tier (EC2, RDS, S3) | Industry standard, vast ecosystem |
| CI/CD | GitHub Actions | Native GitHub integration, free for public repos |

---

## Observability Stack

| Technology | Purpose | Status |
|---|---|---|
| `log/slog` | Structured JSON logging | Active |
| Prometheus | Metrics collection | Planned |
| Grafana | Metrics dashboards | Planned |
| OpenTelemetry | Distributed tracing | Planned (later phase) |
| Jaeger | Trace visualization | Planned (later phase) |

---

## Testing

| Technology | Purpose |
|---|---|
| `testing` (stdlib) | Unit tests, Table-Driven Tests |
| `testify/assert` | Readable assertions |
| `testcontainers-go` | Real PostgreSQL in Docker for integration tests |
| `go test -race` | Race condition detection |
| `k6` | Load testing and concurrency verification |

---

## Local Development Tools

| Tool | Purpose |
|---|---|
| `air` | Live reload for Go |
| `golangci-lint` | Strict linter (must be zero errors) |
| `Bruno` / `Postman` | API testing |
| `k6` | Load testing framework |

---

## Middleware & Security Patterns

| Pattern | Implementation |
|---|---|
| Recovery Middleware | Panic → structured error response, no crash |
| Logging Middleware | Request/response tracking with ResponseWriter wrapper |
| Rate Limiting | Redis INCR+EXPIRE: IP-based (auth, 5/15min) + User-based (API, 100/1min) |
| JWT Auth | `golang-jwt/jwt` — stateless authentication |
| RBAC | `RequireRole()` middleware — roles: `user`, `organizer`, `admin` |
| Audit Logging | Structured slog entries for compliance (GDPR, PCI DSS, HIPAA) |

---

> **Why these choices?** Every tool here is chosen for a specific, production-justified reason. If you want to swap a tool, write an ADR first — explain what problem you're solving that the current tool cannot.
